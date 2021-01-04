package main

import (
	"bytes"
	"fmt"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func broadcast(h *Hub, message []byte) {
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) run() {
	gameStatus := 0 // 0: not started; 1: started
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			log.Println(string(message[:]))
			switch {
			case bytes.Equal(message, []byte("newWinner")):
				output := []byte(`{"update":1}`)
				broadcast(h, output)
			case bytes.Equal(message, []byte("countdown")):
				output := []byte(`{"countdown":1}`)
				gameStatus = 1
				broadcast(h, output)
			case bytes.Equal(message, []byte("gameStart")):
				gameStatus = 1
				output := []byte(fmt.Sprintf(`{"gameStatus":%d}`, gameStatus))
				broadcast(h, output)
			case bytes.Equal(message, []byte("gameNotStart")):
				gameStatus = 0
				output := []byte(fmt.Sprintf(`{"gameStatus":%d}`, gameStatus))
				broadcast(h, output)
			case bytes.Equal(message, []byte("gameStatus")):
				output := []byte(fmt.Sprintf(`{"gameStatus":%d}`, gameStatus))
				broadcast(h, output)
			}
		}
	}
}
