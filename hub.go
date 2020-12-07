// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
)

type WebsocketMessage struct {
	Wsjtx WsjtxMessage `json:"wsjtx"`
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered websocket clients.
	clients map[*Client]bool

	// Inbound messages from the websocket clients.
	command chan []byte

	// Register requests from the websocket clients.
	register chan *Client

	// Unregister requests from websocket clients.
	unregister chan *Client

	// WSJT-X message channel
	wsjtx chan []byte
}

func newHub() *Hub {
	wsjtChan := make(chan []byte, 5)
	go handleWsjtx(wsjtChan)

	return &Hub{
		command:    make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		wsjtx:      wsjtChan,
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Established websocket session with %v", client.conn.RemoteAddr())
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Printf("Disconnected from %v", client.conn.RemoteAddr())
				delete(h.clients, client)
				close(client.send)
			}
		case command := <-h.command:
			// TODO: route this to a backend
			log.Printf("Command from client: %v", command)
		case wsjtxMessage := <-h.wsjtx:
			h.broadcast(wsjtxMessage)
		}
	}
}

func (h *Hub) broadcast(message []byte) {
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}
