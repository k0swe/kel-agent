package config

import "github.com/rs/zerolog"

type Config struct {
	Websocket WebsocketConfig `json:"websocket,omitempty"`
	Wsjtx     WsjtxConfig     `json:"wsjtx,omitempty"`
	LogLevel  zerolog.Level
}

type WebsocketConfig struct {
	// Address is the IP or hostname from which to serve the websocket HTTP
	Address string `json:"address,omitempty"`
	// Port is the TCP port from which to serve the websocket HTTP
	Port uint `json:"port,omitempty"`
	// key is the path to the TLS private key file (needed only if serving securely)
	Key string `json:"key,omitempty"`
	// cert is the path to the TLS public certificate file (needed only if serving securely)
	Cert string `json:"cert,omitempty"`
	// allowedOrigins are the web origins which are allowed by CORS
	AllowedOrigins []string `json:"allowedOrigins,omitempty"`
}

type WsjtxConfig struct {
	// Enabled is whether to listen to WSJT-X
	Enabled bool `json:"enabled"`
	// Address is the IP or hostname on which to listen to WSJT-X
	Address string `json:"address,omitempty"`
	// Port is the UDP port on which to listen to WSJT-X
	Port uint `json:"port,omitempty"`
}
