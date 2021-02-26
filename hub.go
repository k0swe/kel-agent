// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"log"
)

type WebsocketMessage struct {
	// kel-agent version info
	Version string        `json:"version"`
	Wsjtx   WsjtxMessage  `json:"wsjtx,omitempty"`
	Hamlib  HamlibMessage `json:"hamlib,omitempty"`
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
	wsjtx chan WsjtxMessage

	// HamLib message channel
	hamlib chan HamlibMessage
}

var versionInfo string

func newHub(version string) *Hub {
	versionInfo = version
	wsjtChan := make(chan WsjtxMessage, 5)
	go handleWsjtx(wsjtChan)
	hamlibChan := make(chan HamlibMessage, 5)
	go handleHamlib(hamlibChan)

	return &Hub{
		command:    make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		wsjtx:      wsjtChan,
		hamlib:     hamlibChan,
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
			h.broadcast(WebsocketMessage{
				Version: versionInfo,
				Wsjtx:   wsjtxMessage,
			})
		case hamlibMessage := <-h.hamlib:
			h.broadcast(WebsocketMessage{
				Version: versionInfo,
				Hamlib:  hamlibMessage,
			})
		}
	}
}

func (h *Hub) broadcast(message WebsocketMessage) {
	jsn, _ := json.Marshal(message)
	for client := range h.clients {
		select {
		case client.send <- jsn:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}
