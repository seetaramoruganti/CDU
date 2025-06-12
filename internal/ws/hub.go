package ws

import (
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the sensor poller (JSON bytes).
	Broadcast chan []byte

	// Register requests from the WebSocket connections.
	Register chan *Client

	// Unregister requests from WebSocket connections.
	Unregister chan *Client
}

// NewHub creates and returns a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop, processing register, unregister, and broadcast events.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
			log.Printf("client registered, total: %d", len(h.clients))

		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				log.Printf("client unregistered, total: %d", len(h.clients))
			}

		case message := <-h.Broadcast:
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					// If client can't receive, unregister it
					delete(h.clients, client)
					close(client.Send)
				}
			}
		}
	}
}
