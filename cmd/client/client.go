package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Printf("Erro ao criar a requisição com contexto: %v", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Timeout na requisição: %v", err)
			return
		}
		log.Printf("Erro ao chamar a API: %v", err)
		return

	}
	defer resp.Body.Close()

	var data CotacaoResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Printf("Erro ao decodificar a resposta: %v", err)
		return
	}

	fmt.Printf("BID: %v\n", data.Bid)
}
