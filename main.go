package main

import (
	"fmt"

	"github.com/lucascprazeres/cepg/cep"
)

func main() {
	result, _ := cep.Cep("05010000")

	fmt.Println("final", result)
}
