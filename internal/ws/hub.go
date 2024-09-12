// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"encoding/json"
	"github.com/k0swe/kel-agent/internal/config"
	rigctld "github.com/k0swe/kel-agent/internal/rigctld_wrapper"
	wsjtx "github.com/k0swe/kel-agent/internal/wsjtx_wrapper"
	"github.com/rs/zerolog/log"
)

type WebsocketMessage struct {
	// Version is kel-agent version info
	Version string        `json:"version,omitempty"`
	Wsjtx   wsjtx.Message `json:"wsjtx,omitempty"`
	Rigctld interface{}   `json:"rigctld,omitempty"`
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	conf *config.Config

	// Registered websocket clients.
	clients map[*Client]bool

	// Inbound messages from the websocket clients.
	command chan []byte

	// Register requests from the websocket clients.
	register chan *Client

	// Unregister requests from websocket clients.
	unregister chan *Client

	// Wrapper for the WSJT-X connection
	wsjtxHandler *wsjtx.Handler

	// WSJT-X message channel
	wsjtx chan wsjtx.Message

	// Wrapper for the rigctld connection
	rigctldHandler *rigctld.Handler
}

func newHub(c *config.Config) *Hub {
	var wh *wsjtx.Handler
	wsjtChan := make(chan wsjtx.Message, 5)
	if c.Wsjtx.Enabled {
		var err error
		wh, err = wsjtx.NewHandler(c)
		if err != nil {
			log.Warn().Err(err).Msgf("couldn't connect to WSJTX")
		} else {
			go wh.ListenToWsjtx(wsjtChan)
		}
	}

	var rh *rigctld.Handler
	if c.Rigctld.Enabled {
		var err error
		rh, err = rigctld.NewHandler(c)
		if err != nil {
			log.Warn().Err(err).Msgf("couldn't connect to rigctld")
		}
	}

	return &Hub{
		conf:           c,
		command:        make(chan []byte),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		clients:        make(map[*Client]bool),
		rigctldHandler: rh,
		wsjtxHandler:   wh,
		wsjtx:          wsjtChan,
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
			log.Debug().Msgf("Command from client: %v", string(command))
			h.handleClientCommand(command)
		case wsjtxMessage := <-h.wsjtx:
			h.broadcast(WebsocketMessage{
				Version: h.conf.VersionInfo,
				Wsjtx:   wsjtxMessage,
			})
		}
	}
}

func (h *Hub) broadcast(message WebsocketMessage) {
	jsn, _ := json.Marshal(message)
	log.Trace().Msgf("broadcasting: %v", string(jsn))
	for client := range h.clients {
		select {
		case client.send <- jsn:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) handleClientCommand(command []byte) {
	var msg = &WebsocketMessage{}
	log.Trace().Msgf("received client command: %v", string(command))
	if err := json.Unmarshal(command, msg); err != nil {
		log.Warn().Err(err).Msg("failed to parse client command; dropping")
		return
	}
	if msg.Wsjtx.MsgType != "" {
		// Don't know all the payload types here, so re-marshal just that and handle in wrapper
		payload, _ := json.Marshal(msg.Wsjtx.Payload)
		if err := h.wsjtxHandler.HandleClientCommand(msg.Wsjtx.MsgType, payload); err != nil {
			log.Warn().Err(err).Msg("failed to handle wsjtx client command; dropping")
			return
		}
	}
	if msg.Rigctld != "" {
		// Don't know all the payload types here, so re-marshal just that and handle in wrapper
		payload, _ := json.Marshal(msg.Rigctld)
		response, err := h.rigctldHandler.HandleClientCommand(payload)
		if err != nil {
			log.Warn().Err(err).Msg("failed to handle rigctld client command")

		}
		h.broadcast(WebsocketMessage{
			Version: h.conf.VersionInfo,
			Rigctld: response,
		})
	}
}
