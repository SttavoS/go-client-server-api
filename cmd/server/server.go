package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	http.HandleFunc("/cotacao", BuscaCotacao)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func BuscaCotacao(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Cotacao{
		Bid: "5.50",
	}
	json.NewEncoder(w).Encode(response)
}
