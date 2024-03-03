package models

type Result struct {
	Service      string `json:"service"`
	Cep          string `json:"cep"`
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
}

type Service func() (Result, error)
