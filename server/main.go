package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Cotacao struct {
	USDBRL CotacaoUSDBRL `json:"USDBRL"`
}

type CotacaoUSDBRL struct {
	Code        string  `json:"code"`
	Codein      string  `json:"codein"`
	Name        string  `json:"name"`
	High        string  `json:"high"`
	Low         string  `json:"low"`
	Bid         float64 `json:"bid,string"`
	Ask         string  `json:"ask"`
	Create_date string  `json:"create_date"`
}

func main() {
	log.Println("Initializing server")
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request iniciada")

	body := getRealAndDollarPrice()

	var cotacao Cotacao
	err := json.Unmarshal(body, &cotacao)
	if err != nil {
		log.Println("Erro ao parsear dados do JSON")
		panic(err)
	}
	saveToDatabase(&cotacao)

	defer log.Println("Request finalizada")
}

func getRealAndDollarPrice() []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
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

func saveToDatabase(c *Cotacao) {
	db, err := sql.Open("mysql", "root:12345678@tcp(mysql-5.7:3306)/goexpert")
	if err != nil {
		log.Println("Erro ao conectar no banco de dados")
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Millisecond)
	defer cancel()

	stmt, err := db.Prepare("insert into prices(id, value, updated_at) values(?, ?, ?)")
	if err != nil {
		log.Println("Erro ao preparar stetement no banco de dados")
		panic(err)
	}

	_, err = stmt.ExecContext(ctx, uuid.New().String(), c.USDBRL.Bid, c.USDBRL.Create_date)
	if err != nil {
		log.Println("Erro ao executar insert no banco de dados")
		panic(err)
	}
}
