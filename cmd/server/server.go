package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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

func main() {
	http.HandleFunc("/cotacao", BuscaCotacao)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func BuscaCotacao(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
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
