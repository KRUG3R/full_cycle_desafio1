package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Valor string `json:"cotacao"`
}

func main() {
	fmt.Println("iniciando Client")
	var cotacao Cotacao
	cotacao.GetBid()
	fmt.Println("Valor da cotação:")
	fmt.Println(cotacao.Valor)
	WriteFile(cotacao.Valor)
	fmt.Println("fim do client")
}

func (cotacao *Cotacao) GetBid() {
	ctx1 := context.Background()
	ctx1, cancel := context.WithTimeout(ctx1, time.Millisecond*300)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx1, "GET", "http://localhost:8080/cotacao", nil)
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
	err = json.Unmarshal(body, cotacao)
	if err != nil {
		panic(err)
	}
}

func WriteFile(valor string) {
	file, err := os.OpenFile("arquivo.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString("Dólar: {" + valor + "}\n")

	if err != nil {
		panic(err)
	}
}
