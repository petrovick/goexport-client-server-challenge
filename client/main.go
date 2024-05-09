package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type QuotationUSDBRL struct {
	Bid float64 `json:"bid"`
}

func main() {
	body := getRealAndDollarPrice()
	log.Println(body)

	var quotation QuotationUSDBRL
	err := json.Unmarshal(body, &quotation)
	if err != nil {
		log.Println("Erro ao parsear dados do JSON")
		panic(err)
	}

	log.Println(quotation)
	writeToFile(&quotation)
}

func writeToFile(quotation *QuotationUSDBRL) {
	f, err := os.Create("arquivo.txt")
	if err != nil {
		panic(err)
	}

	tamanho, err := f.Write([]byte(fmt.Sprintf("Dólar: %f", quotation.Bid)))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Arquivo criado com sucesso! Tamanho: %d bytes\n", tamanho)
	f.Close()
}

func getRealAndDollarPrice() []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Println("Erro ao montar requisição dos dados da cotação")
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro ao buscar dados da cotação")
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Erro ao ler os dados da API")
		panic(err)
	}

	return body
}
