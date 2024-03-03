package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/lucascprazeres/cepg/models"
)

type CepData struct {
	Cep        string `json:"cep"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
	Logradouro string `json:"logradouroDNEC"`
}

type CorreiosResponse struct {
	Dados []CepData `json:"dados"`
}

func Correios(cep *string) models.Service {
	return func() (models.Result, error) {
		endpoint := "https://buscacepinter.correios.com.br/app/endereco/carrega-cep-endereco.php"

		formData := url.Values{}
		formData.Set("endereco", *cep)
		formData.Set("tipoCEP", "ALL")

		headers := map[string]string{
			"content-type": "application/x-www-form-urlencoded; charset=UTF-8",
			"referer":      "https://buscacepinter.correios.com.br/app/endereco/index.php",
		}

		request, _ := http.NewRequest("POST", endpoint, strings.NewReader(formData.Encode()))

		for key, value := range headers {
			request.Header.Set(key, value)
		}

		client := &http.Client{}

		response, err := client.Do(request)
		if err != nil {
			return models.Result{}, err
		}
		defer response.Body.Close()

		return parseCorreiosResponse(response)
	}
}

func parseCorreiosResponse(response *http.Response) (models.Result, error) {
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return models.Result{}, errors.New("failed reading correios response")
	}

	var result CorreiosResponse
	err = json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println(err)
		return models.Result{}, errors.New("failed unmarshalling correios response")
	}

	parsedStreetName := strings.Split(result.Dados[0].Logradouro, " - ")[0]

	return models.Result{
		Service:      "correios",
		Cep:          result.Dados[0].Cep,
		State:        result.Dados[0].Uf,
		City:         result.Dados[0].Localidade,
		Neighborhood: result.Dados[0].Bairro,
		Street:       parsedStreetName,
	}, nil
}
