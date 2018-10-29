package main

import (
	"fmt"
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

// newHub create instance of connections hub structure object
func newHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// getClients returns slice of clients remote addresses
func (h *Hub) getClients() []string {
	clientsMap := []string{}
	for cl := range h.clients {
		record := fmt.Sprintf("ip=%s, dir=%s", cl.conn.RemoteAddr().String(), cl.dir)
		clientsMap = append(clientsMap, record)
	}

	return clientsMap
}

// run function listens for channels signals to add or remove client connection
// to hub clients collection
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
	}
}
