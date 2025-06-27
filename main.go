package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"tech-tarot-backend/cardstack"
	"tech-tarot-backend/controllers"
)

func main() {
	cfg := LoadConfig()

	err := cardstack.Init(cfg.CardStackPath)
	if err != nil {
		panic("no cardstack found")
	}

	ctx := context.Background()
	mux := http.NewServeMux()
	cardStackController := controllers.NewCardStackController()
	sseController := controllers.NewSSEController(ctx)

	cardStackController.RegisterRoutes(mux)
	sseController.RegisterRoutes(mux)

	mux.HandleFunc("/", welcome)
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(cfg.HTTPServerPort), mux))
}

func welcome(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Welcome to the tech tarot backend!")
}
