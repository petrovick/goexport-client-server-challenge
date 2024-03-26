package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Quotation struct {
	USDBRL QuotationUSDBRL `json:"USDBRL"`
}

type QuotationUSDBRL struct {
	Code        string  `json:"code"`
	Codein      string  `json:"codein"`
	Name        string  `json:"name"`
	High        string  `json:"high"`
	Low         string  `json:"low"`
	Bid         float64 `json:"bid,string"`
	Ask         string  `json:"ask"`
	Create_date string  `json:"create_date"`
}

type QuotationDTO struct {
	Bid float64
}

func main() {
	log.Println("Initializing server")
	http.HandleFunc("/cotacao", quotationHandler)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}

func quotationHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Request iniciada")

	body := getRealAndDollarQuotation()

	var quotation Quotation
	err := json.Unmarshal(body, &quotation)
	if err != nil {
		log.Println("Erro ao parsear dados do JSON")
		panic(err)
	}
	saveToDatabase(&quotation)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var quotationDTO QuotationDTO = QuotationDTO{}
	quotationDTO.Bid = quotation.USDBRL.Bid

	result, err := json.Marshal(quotationDTO)

	w.Write(result)
	defer log.Println("Request finalizada")
}

func getRealAndDollarQuotation() []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	awesomeUrl := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	req, err := http.NewRequestWithContext(ctx, "GET", awesomeUrl, nil)
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

func saveToDatabase(quotation *Quotation) {
	db, err := sql.Open("sqlite3", "./data/goexpert-database.db")
	if err != nil {
		log.Println("Erro ao conectar no banco de dados")
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	stmt, err := db.Prepare("insert into quotations(id, value, updated_at) values(?, ?, ?)")
	if err != nil {
		log.Println("Erro ao preparar stetement no banco de dados")
		panic(err)
	}

	_, err = stmt.ExecContext(ctx, uuid.New().String(), quotation.USDBRL.Bid, quotation.USDBRL.Create_date)
	if err != nil {
		log.Println("Erro ao executar insert no banco de dados")
		panic(err)
	}
}
