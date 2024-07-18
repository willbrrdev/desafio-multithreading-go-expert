package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	urlViaCep    = "https://viacep.com.br/ws/%s/json/"
	urlBrasilApi = "https://brasilapi.com.br/api/cep/v1/%s"
)

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
}

type BrasilApiResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

func getViaCep(cep string) ViaCepResponse {
	url := fmt.Sprintf(urlViaCep, cep)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var viaCepResponse ViaCepResponse
	err = json.Unmarshal(body, &viaCepResponse)
	if err != nil {
		log.Fatal(err)
	}

	return viaCepResponse
}

func getBrasilApi(cep string) BrasilApiResponse {
	url := fmt.Sprintf(urlBrasilApi, cep)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var brasilApiResponse BrasilApiResponse
	err = json.Unmarshal(body, &brasilApiResponse)
	if err != nil {
		log.Fatal(err)
	}

	return brasilApiResponse
}

func main() {
	cep := "54780455"

	viacep := make(chan ViaCepResponse)
	brasilapi := make(chan BrasilApiResponse)

	// ViaCep
	go func() {
		for {
			response := getViaCep(cep)
			viacep <- response
		}

	}()

	// BrasilApi
	go func() {
		for {
			response := getBrasilApi(cep)
			brasilapi <- response
		}
	}()

	for i := 0; i < 1; i++ {
		select {
		case resp := <-viacep:
			fmt.Println("ViaCep", resp)
		case resp := <-brasilapi:
			fmt.Println("BrasilApi", resp)
		case <-time.After(time.Second):
			println("Timeout")
		}
	}
}
