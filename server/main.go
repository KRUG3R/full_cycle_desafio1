package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("Iniciando o Server")
	defer fmt.Println("Finalizando Server")
	http.HandleFunc("/cotacao", handler)
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
	Timestamp  string `json:"timestamp"`
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
		if ctx1.Err() == context.DeadlineExceeded {
			fmt.Println("Timeout da chamada atingido")
		}
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
	err := PersisteDB(cambio)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		cotacao := map[string]interface{}{"cotacao": cambio.USDBRL.Bid}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cotacao)
	}

}

func PersisteDB(cambio Cambio) error {
	query := `INSERT INTO cotacao (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) Values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	db, err := sql.Open("sqlite3", "../sqlite_setup/desafio1.db")
	if err != nil {
		msg := fmt.Errorf("Erro ao abrir o banco de dados: %v", err)
		fmt.Println(msg)
		return msg
	}
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		msg := fmt.Errorf("Erro ao preparar a query: %v", err)
		fmt.Println(msg)
		return msg
	}
	defer stmt.Close()

	ctx2 := context.Background()
	ctx2, cancel := context.WithTimeout(ctx2, time.Millisecond*10)
	defer cancel()

	_, err = stmt.ExecContext(ctx2, cambio.USDBRL.Code, cambio.USDBRL.Codein, cambio.USDBRL.Name, cambio.USDBRL.High, cambio.USDBRL.Low, cambio.USDBRL.VarBid, cambio.USDBRL.PctChange, cambio.USDBRL.Bid, cambio.USDBRL.Ask, cambio.USDBRL.Timestamp, cambio.USDBRL.CreateDate)
	if err != nil {
		if ctx2.Err() == context.DeadlineExceeded {
			fmt.Println("Timeout da chamada ao SQLITE atingido")
		}
		msg := fmt.Errorf("Erro ao executar a query: %v", err)
		fmt.Println(msg)
		return msg
	}
	return nil

}
