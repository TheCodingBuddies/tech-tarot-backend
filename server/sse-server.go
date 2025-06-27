package server

import (
	"context"
	"fmt"
	"log"
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

func (sseServer *EventServer) Broadcast(message string) {
	broadcastMessage := fmt.Sprintf("%s\n\n", message)
	for c := range sseServer.clients {
		select {
		case c <- []byte(broadcastMessage):
		case <-time.After(1 * time.Second):
			log.Printf("Broadcast timeout for client: %v", c)
		}
	}

}
