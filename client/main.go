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

type CotacaoUSDBRL struct {
	Bid float64 `json:"bid"`
}

func main() {
	body := getRealAndDollarPrice()
	log.Println(body)

	var cotacao CotacaoUSDBRL
	err := json.Unmarshal(body, &cotacao)
	if err != nil {
		log.Println("Erro ao parsear dados do JSON")
		panic(err)
	}

	log.Println(cotacao)
	writeToFile(&cotacao)
}

func writeToFile(cotacao *CotacaoUSDBRL) {
	f, err := os.Create("arquivo.txt")
	if err != nil {
		panic(err)
	}

	tamanho, err := f.Write([]byte(fmt.Sprintf("Dólar: %f", cotacao.Bid)))
	// tamanho, err := f.WriteString("Hello, World!")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Arquivo criado com sucesso! Tamanho: %d bytes\n", tamanho)
	f.Close()
}

func getRealAndDollarPrice() []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://go-server:8081/", nil)
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
