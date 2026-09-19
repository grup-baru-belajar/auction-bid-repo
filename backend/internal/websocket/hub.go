package websocket

import (
	"sync"
)

type Client struct {
	auctionID int64
	send      chan []byte
	hub       *Hub
}

type Hub struct {
	mu      sync.RWMutex
	clients map[int64]map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[int64]map[*Client]struct{}),
	}
}

func (h *Hub) register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[c.auctionID]; !ok {
		h.clients[c.auctionID] = make(map[*Client]struct{})
	}
	h.clients[c.auctionID][c] = struct{}{}
}

func (h *Hub) unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if group, ok := h.clients[c.auctionID]; ok {
		delete(group, c)
		if len(group) == 0 {
			delete(h.clients, c.auctionID)
		}
	}
}

func (h *Hub) Broadcast(auctionID int64, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for c := range h.clients[auctionID] {
		select {
		case c.send <- msg:
		default:
		}
	}
}
