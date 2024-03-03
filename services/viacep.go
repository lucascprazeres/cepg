package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/lucascprazeres/cepg/models"
)

func Viacep(cep *string) models.Service {
	return func() (models.Result, error) {
		url := fmt.Sprintf("https://viacep.com.br/ws/%v/json/", *cep)
		response, err := http.Get(url)
		if err != nil {
			return models.Result{}, errors.New("failed fetching viacep")
		}
		defer response.Body.Close()

		return parseViacepResponse(response)
	}
}

func parseViacepResponse(response *http.Response) (models.Result, error) {
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return models.Result{}, errors.New("failed reading viacep response")
	}

	var result map[string]string

	err = json.Unmarshal(data, &result)
	if err != nil {
		return models.Result{}, errors.New("failed unmarshalling viacep response")
	}

	return models.Result{
		Service:      "viacep",
		Cep:          result["cep"],
		State:        result["uf"],
		City:         result["localidade"],
		Neighborhood: result["bairro"],
		Street:       result["logradouro"],
	}, nil
}
