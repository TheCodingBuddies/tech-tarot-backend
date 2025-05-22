package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"tech-tarot-backend/server"
)

func main() {
	http.HandleFunc("/", welcome)
	http.HandleFunc("/start", start)
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
	json.NewEncoder(writer).Encode(server.NewStack().Draw3())
}

func welcome(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Welcome to the tech tarot backend!")
}

func start(writer http.ResponseWriter, request *http.Request) {
	/* ToDo: refactor early return */
	if request.Method != http.MethodPost {
		http.Error(writer, "Only Post method is allowed", http.StatusMethodNotAllowed)
		return
	}
	/* ToDo:send start to frontend */
	fmt.Fprint(writer, "Tech Tarot started")
}
