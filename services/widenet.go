package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lucascprazeres/cepg/models"
)

type WidenetResponse struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

func Widenet(cep *string) models.Service {
	return func() (models.Result, error) {
		cepValue := *cep
		cepWithDash := strings.Join([]string{cepValue[:5], "-", cepValue[5:]}, "")

		url := fmt.Sprintf("https://cdn.apicep.com/file/apicep/%v.json", cepWithDash)

		response, err := http.Get(url)
		if err != nil {
			return models.Result{}, errors.New("failed fetching widenet")
		}
		defer response.Body.Close()

		return parseWidenetResponse(response)
	}
}

func parseWidenetResponse(response *http.Response) (models.Result, error) {
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return models.Result{}, errors.New("failed reading viacep response")
	}

	var result WidenetResponse

	err = json.Unmarshal(data, &result)
	if err != nil {
		return models.Result{}, errors.New("failed unmarshalling viacep response")
	}

	parsedStreetName := strings.Split(result.Address, " - ")[0]

	return models.Result{
		Service:      "widenet",
		Cep:          result.Code,
		State:        result.State,
		City:         result.City,
		Neighborhood: result.District,
		Street:       parsedStreetName,
	}, nil
}
