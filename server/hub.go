package main

import (
	"fmt"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		}

		for _ = range time.Tick(3000 * time.Millisecond) {
			fmt.Println("CONNCETED WORKERS")
			for cl := range h.clients {
				fmt.Println(" ==> " + cl.conn.RemoteAddr().String())
			}
		}
	}
}
