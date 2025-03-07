package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	// _ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("Iniciando o Server")
	defer fmt.Println("Finalizando Server")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

type Detalhes struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestap   string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type Cambio struct {
	USDBRL Detalhes `json:"USDBRL"`
}

func (cambio *Cambio) GetCambio() {
	ctx1 := context.Background()
	ctx1, cancel := context.WithTimeout(ctx1, time.Millisecond*200)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx1, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, cambio)
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	var cambio Cambio
	cambio.GetCambio()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cambio)
}

// func PersisteDB(cambio Cambio) {
// 	db, err := sql.Open("sqlite3", "./cambio.db")
// 	insertSQL := `INSERT INTO cambio (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) Values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

// }
