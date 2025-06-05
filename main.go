package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"tech-tarot-backend/cardstack"
	"tech-tarot-backend/server"
)

func main() {
	ctx := context.Background()
	sse := server.NewSSEServer(ctx)
	sse.Start()

	http.HandleFunc("/", welcome)
	http.HandleFunc("/connect", sse.Connect)
	http.HandleFunc("/start", sse.StartGame)
	http.HandleFunc("/cards", withCORS(draw))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		fn(w, r)
	}
}

func draw(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(cardstack.NewStack().Draw3())
}

func welcome(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Welcome to the tech tarot backend!")
}
