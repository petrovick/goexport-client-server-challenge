package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

type QuotationsModel struct {
	ID        string `gorm:"primaryKey"`
	Value     float64
	UpdatedAt string
	gorm.Model
}

func main() {
	log.Println("Initializing server")
	http.HandleFunc("/cotacao", quotationHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
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
	db, err := gorm.Open(sqlite.Open("goexpert-database.db"), &gorm.Config{})
	if err != nil {
		log.Println("Erro ao conectar no banco de dados")
		panic(err)
	}

	db.AutoMigrate(&QuotationsModel{})

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	quotationsModel := QuotationsModel{
		ID:        (uuid.New().String()),
		Value:     quotation.USDBRL.Bid,
		UpdatedAt: quotation.USDBRL.Create_date,
	}

	db.WithContext(ctx).Create(&quotationsModel)
	ContextExecution(ctx, "DB")
}

func ContextExecution(ctx context.Context, name string) {
	select {
	case <-ctx.Done():
		panic("Time exceeded in " + name)
	case <-time.After(time.Millisecond * 10):
		println("success")
	}
}
