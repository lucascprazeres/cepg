package cep

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/lucascprazeres/cepg/models"
	"github.com/lucascprazeres/cepg/services"
)

func Cep(cep string) (models.Result, error) {
	removeSpecialCharacters(&cep)
	err := validateCepLength(&cep)

	if err != nil {
		return models.Result{}, err
	}

	result, err := fetchCepFromServices(&cep)

	if err != nil {
		return models.Result{}, err
	}

	return result, nil
}

func removeSpecialCharacters(cepRawValue *string) {
	r := regexp.MustCompile("[^0-9]")
	*cepRawValue = r.ReplaceAllString(*cepRawValue, "")
}

func validateCepLength(cep *string) error {
	const CEP_MAXIMUM_LENGTH = 8

	if len(*cep) > CEP_MAXIMUM_LENGTH {
		return fmt.Errorf("%v is not a valid cep since it has more than 8 digits", *cep)
	}

	return nil
}

func fetchCepFromServices(cep *string) (models.Result, error) {
	services := map[string]models.Service{
		"viacep":   services.Viacep(cep),
		"correios": services.Correios(cep),
		"widenet":  services.Widenet(cep),
	}

	var wg sync.WaitGroup
	resultCh := make(chan models.Result, 1)

	for _, service := range services {
		wg.Add(1)
		go fetchCep(service, &wg, &resultCh)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	result := <-resultCh

	return result, nil
}

func fetchCep(service models.Service, wg *sync.WaitGroup, ch *chan models.Result) {
	defer wg.Done()

	result, err := service()

	if err == nil {
		*ch <- result
	}
}
