package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"tech-tarot-backend/server"
	"time"
)

type UserData struct {
	Name string `json:"user"`
}

type SSEController struct {
	sseServer *server.EventServer
}

func NewSSEController(context context.Context) *SSEController {
	sse := server.NewSSEServer(context)
	sse.Start()
	return &SSEController{sseServer: sse}
}

func (controller *SSEController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/connect", withCORS(controller.connect))
	mux.HandleFunc("/start", controller.startGame)
}

func (controller *SSEController) connect(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	writer.Header().Set("X-Accel-Buffering", "no")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	flusher, ok := writer.(http.Flusher)
	if !ok {
		http.Error(writer, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithCancel(request.Context())
	defer cancel()

	message := make(chan []byte)
	controller.sseServer.ConnectClient <- message

	keepAliveTicker := time.NewTicker(10 * time.Second)

	for {
		select {
		case msg, ok := <-message:
			log.Printf("message to broadcast: %v\n", string(msg))
			if !ok {
				return
			}
			_, err := fmt.Fprintf(writer, "data: %s\n\n", msg)
			if err != nil {
				log.Printf("error writing message: %v", err)
				return
			}
			flusher.Flush()
		case <-keepAliveTicker.C:
			_, err := fmt.Fprintf(writer, "data: keepalive\n\n")
			if err != nil {
				log.Printf("error writing keep alive message: %v\n", err)
				return
			}
			flusher.Flush()
		case <-ctx.Done():
			controller.sseServer.CloseClient <- message
			return
		}
	}
}

func (controller *SSEController) startGame(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Only Post method is allowed", http.StatusMethodNotAllowed)
		return
	}
	userData := UserData{}
	err := json.NewDecoder(request.Body).Decode(&userData)
	if err != nil {
		http.Error(writer, "Invalid user data", http.StatusBadRequest)
		return
	}
	controller.sseServer.Broadcast(userData.Name)
	fmt.Fprint(writer, "Tech Tarot started")
}
