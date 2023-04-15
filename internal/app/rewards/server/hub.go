package server

import (
	"sync"
	"ttv-cli/internal/pkg/config"
)

type Hub struct {
	clients     map[*Client]bool
	register    chan *Client
	unregister  chan *Client
	broadcast   chan []byte
	rewardsById sync.Map
}

func newHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan []byte),
		rewardsById: sync.Map{},
	}
}

func (h *Hub) run(config config.Config, streamer string) {
	go h.pumpEvents(config, streamer)
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.rewardsById.Range(func(_, value any) bool {
				r := value.(*reward)
				client.send <- r.toBytes()
				return true
			})

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// If the send buffer is full, assume that the client is dead or stuck and unregister them
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
