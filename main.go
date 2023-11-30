package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	canalViaCep := make(chan ViaCep, 1)
	// Esta API no momento em que estou codificando está com falha, por este motivo incluí a OpenCep
	canalApiCep := make(chan ApiCep, 1)
	canalOpenCep := make(chan OpenCep, 1)
	canalBrasilApiCep := make(chan BrasilApiCep, 1)

	cepEntrada := "88354270"

	go func() {
		resultado, err := BuscaCepComViaCep(ctx, cepEntrada)
		if err != nil {
			println("Falha ao buscar CEP com ViaCep:", err.Error())
			return
		}
		if resultado != nil {
			canalViaCep <- *resultado
		}
	}()

	go func() {
		resultado, err := buscaCepComApiCep(ctx, cepEntrada)
		if err != nil {
			println("Falha ao buscar CEP com ApiCep:", err.Error())
			return
		}

		if resultado != nil {
			canalApiCep <- *resultado
		}
	}()

	go func() {
		resultado, err := BuscaCepComOpenCep(ctx, cepEntrada)
		if err != nil {
			println("Falha ao buscar CEP com OpenCep:", err.Error())
			return
		}
		if resultado != nil {
			canalOpenCep <- *resultado
		}
	}()

	go func() {
		resultado, err := BuscaCepComBrasilApiCep(ctx, cepEntrada)
		if err != nil {
			println("Falha ao buscar CEP com BrasilApiCep:", err.Error())
			return
		}
		if resultado != nil {
			canalBrasilApiCep <- *resultado
		}
	}()

	select {
	case <-ctx.Done():
		println("Timeout atingido")
	case res := <-canalViaCep:
		fmt.Printf("Resultado obtido com ViaCep: %+v\n", res)
	case res := <-canalApiCep:
		fmt.Printf("Resultado obtido com ApiCep: %+v\n", res)
	case res := <-canalOpenCep:
		fmt.Printf("Resultado obtido com OpenCep: %+v\n", res)
	case res := <-canalBrasilApiCep:
		fmt.Printf("Resultado obtido com BrasilApiCep: %+v\n", res)
	}
}
