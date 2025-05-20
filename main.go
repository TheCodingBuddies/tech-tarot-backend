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
	http.HandleFunc("/cards", draw)
	log.Fatal(http.ListenAndServe(":8080", nil))
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
