package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type UsdBrl struct {
	Cotacao Cotacao `json:"USDBRL"`
}

type Cotacao struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	Bid        string `json:"bid"`
	CreateDate string `json:"create_date"`
}

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite", "./data/cotacoes.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cotacoes (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	bid TEXT NOT NULL,
    	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/cotacao", BuscaCotacao)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func BuscaCotacao(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		log.Printf("Erro ao criar a requisição com contexto: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Erro ao realizar requisição")
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Timeout na requisição: %v", err)
			writeErrorResponse(w, http.StatusGatewayTimeout, "Requisição demorou demais")
			return
		}
		log.Printf("Erro ao chamar a API: %v", err)
		writeErrorResponse(w, http.StatusBadGateway, "Erro ao realizar requisição")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API respondeu com status %v", resp.StatusCode)
		writeErrorResponse(w, http.StatusBadGateway, "Erro ao buscar dados")
		return
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Erro ao ler a resposta: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Erro ao ler a resposta")
		return
	}

	var data UsdBrl
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		log.Printf("Erro ao converter a resposta: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Erro ao ler a resposta")
		return
	}

	ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelDB()
	_, err = db.ExecContext(ctxDB, "INSERT INTO cotacoes (bid) VALUES (?)", data.Cotacao.Bid)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Timeout ao salvar no banco de dados: %v", err)
			writeErrorResponse(w, http.StatusGatewayTimeout, "Requisição demorou demais")
			return
		}
		log.Printf("Erro ao salvar a resposta no banco de dados: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Erro ao salvar a resposta")
		return
	}

	cotacao := CotacaoResponse{
		Bid: data.Cotacao.Bid,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao)
}

func writeErrorResponse(w http.ResponseWriter, status int, message string) {
	errResponse := ErrorResponse{
		StatusCode: status,
		Message:    message,
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errResponse)
}
