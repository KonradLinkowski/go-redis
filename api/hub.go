package main

import "sync"

type Hub struct {
	clients    map[chan string]bool
	register   chan chan string
	unregister chan chan string
	broadcast  chan string
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[chan string]bool),
		register:   make(chan chan string),
		unregister: make(chan chan string),
		broadcast:  make(chan string),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			delete(h.clients, client)
			close(client)
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			println("Broadcasting message to clients")
			for client := range h.clients {
				select {
				case client <- message:
				default:
				}
			}
			h.mu.Unlock()
		}
	}
}
