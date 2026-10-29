package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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

	if resp.StatusCode != http.StatusOK {
		log.Printf("API respondeu com status %v", resp.StatusCode)
		return
	}

	var data CotacaoResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Printf("Erro ao decodificar a resposta: %v", err)
		return
	}

	file, err := os.Create("./data/cotacao.txt")
	if err != nil {
		log.Printf("Erro ao criar arquivo: %v", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("Dólar: %v", data.Bid))
	if err != nil {
		log.Printf("Erro ao escrever no arquivo: %v", err)
		return
	}
	fmt.Printf("Dólar: %v\n", data.Bid)
}
