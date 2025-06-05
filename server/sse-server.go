package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type EventServer struct {
	context       context.Context
	cancel        context.CancelFunc
	clients       map[chan []byte]struct{} // Map to keep track of connected clients
	sync          sync.Mutex
	ConnectClient chan chan []byte
	CloseClient   chan chan []byte
	srvWg         sync.WaitGroup
	reqWg         sync.WaitGroup
}

func NewSSEServer(parentCtx context.Context) *EventServer {
	ctx, cancel := context.WithCancel(parentCtx)
	return &EventServer{
		context:       ctx,
		cancel:        cancel,
		clients:       make(map[chan []byte]struct{}),
		ConnectClient: make(chan chan []byte),
		CloseClient:   make(chan chan []byte),
	}
}

func (sseServer *EventServer) Start() {
	sseServer.srvWg.Add(1)
	go func() {
		sseServer.run()
	}()
}

func (sseServer *EventServer) run() {
	defer sseServer.srvWg.Done()

	for {
		select {
		case <-sseServer.context.Done():
			sseServer.reqWg.Wait()
			return

		case clientConnection := <-sseServer.ConnectClient:
			sseServer.sync.Lock()
			sseServer.clients[clientConnection] = struct{}{}
			log.Printf("connecting client %v\n", clientConnection)
			sseServer.sync.Unlock()

		case clientDisconnect := <-sseServer.CloseClient:
			sseServer.sync.Lock()
			log.Printf("closing client %v", clientDisconnect)
			delete(sseServer.clients, clientDisconnect)
			sseServer.sync.Unlock()
		}
	}
}

func (sseServer *EventServer) Connect(writer http.ResponseWriter, request *http.Request) {
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

	message := make(chan []byte)
	sseServer.ConnectClient <- message

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
		}
	}
}

func (sseServer *EventServer) StartGame(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Only Post method is allowed", http.StatusMethodNotAllowed)
		return
	}
	sseServer.broadcast("start-game")
	fmt.Fprint(writer, "Tech Tarot started")
}

func (sseServer *EventServer) broadcast(message string) {
	for c := range sseServer.clients {
		broadcastMessage := fmt.Sprintf("%s\n\n", message)
		c <- []byte(broadcastMessage)
	}
}
