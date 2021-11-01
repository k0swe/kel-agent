// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"

	"github.com/k0swe/kel-agent/internal/wsjtx"
	"github.com/rs/zerolog/log"
)

type WebsocketMessage struct {
	// Version is kel-agent version info
	Version string        `json:"version,omitempty"`
	Wsjtx   wsjtx.Message `json:"wsjtx,omitempty"`
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
	wsjtx chan wsjtx.Message
}

func newHub() *Hub {
	wsjtChan := make(chan wsjtx.Message, 5)
	if conf.Wsjtx.Enabled {
		go wsjtx.HandleWsjtx(conf, wsjtChan)
	}

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
			log.Debug().Msgf("Established websocket session with %v", client.conn.RemoteAddr())
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Debug().Msgf("Disconnected from %v", client.conn.RemoteAddr())
				delete(h.clients, client)
				close(client.send)
			}
		case command := <-h.command:
			// TODO: route this to a backend
			log.Debug().Msgf("Command from client: %v", command)
		case wsjtxMessage := <-h.wsjtx:
			h.broadcast(WebsocketMessage{
				Version: versionInfo,
				Wsjtx:   wsjtxMessage,
			})
		}
	}
}

func (h *Hub) broadcast(message WebsocketMessage) {
	log.Trace().Msgf("broadcasting: %v", message)
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
