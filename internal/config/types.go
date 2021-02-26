package config

import (
	"encoding/json"
	"github.com/rs/zerolog"
)

type Config struct {
	Websocket   WebsocketConfig `json:"websocket,omitempty" yaml:"websocket,omitempty"`
	Wsjtx       WsjtxConfig     `json:"wsjtx,omitempty" yaml:"wsjtx,omitempty"`
	Hamlib      HamlibConfig    `json:"hamlib,omitempty" yaml:"hamlib,omitempty"`
	VersionInfo string          `json:"-" yaml:"-"`
}

func (c Config) MarshalZerologObject(e *zerolog.Event) {
	j, _ := json.Marshal(c)
	e.RawJSON("config", j)
}

type WebsocketConfig struct {
	// Address is the IP or hostname from which to serve the websocket HTTP
	Address string `json:"address,omitempty" yaml:"address,omitempty"`
	// Port is the TCP port from which to serve the websocket HTTP
	Port uint `json:"port,omitempty" yaml:"port,omitempty"`
	// key is the path to the TLS private key file (needed only if serving securely)
	Key string `json:"key,omitempty" yaml:"key,omitempty"`
	// cert is the path to the TLS public certificate file (needed only if serving securely)
	Cert string `json:"cert,omitempty" yaml:"cert,omitempty"`
	// allowedOrigins are the web origins which are allowed by CORS
	AllowedOrigins []string `json:"allowedOrigins,omitempty" yaml:"allowedOrigins,omitempty"`
}

type WsjtxConfig struct {
	// Enabled is whether to listen to WSJT-X
	Enabled bool `json:"enabled" yaml:"enabled"`
	// Address is the IP or hostname on which to listen to WSJT-X
	Address string `json:"address,omitempty" yaml:"address,omitempty"`
	// Port is the UDP port on which to listen to WSJT-X
	Port uint `json:"port,omitempty" yaml:"port,omitempty"`
}

type HamlibConfig struct {
	// Enabled is whether to listen to a rig via Hamlib
	Enabled bool `json:"enabled" yaml:"enabled"`
	// RetrySeconds is the time to wait between connection attempts
	RetrySeconds int `json:"retrySeconds,omitempty" yaml:"retrySeconds,omitempty"`
	// RigModel is the Hamlib rig model number (see `rigctl -l`)
	RigModel int `json:"rigModel,omitempty" yaml:"rigModel,omitempty"`
	// RigPort is the port type, e.g. "RIG_PORT_SERIAL"
	// (see https://github.com/xylo04/goHamlib/blob/3752aec70bb9298eedaed8a58e834ffd92261ce0/goHamlib.go#L758-L772)
	RigPort string `json:"rigPort,omitempty" yaml:"rigPort,omitempty" jsonschema:"enum=RIG_PORT_NONE,enum=RIG_PORT_SERIAL,enum=RIG_PORT_NETWORK,enum=RIG_PORT_DEVICE,enum=RIG_PORT_PACKET,enum=RIG_PORT_DTMF,enum=RIG_PORT_ULTRA,enum=RIG_PORT_RPC,enum=RIG_PORT_PARALLEL,enum=RIG_PORT_USB,enum=RIG_PORT_UDP_NETWORK,enum=RIG_PORT_CM108"`
	// PortName is the name of the port, e.g. "/dev/ttyUSB0"
	PortName string `json:"portName,omitempty" yaml:"portName,omitempty"`
	// BaudRate is the baud rate for communicating between the computer and radio
	BaudRate int `json:"baudRate,omitempty" yaml:"baudRate,omitempty"`
	// DataBits is the number of bits per character
	DataBits int `json:"dataBits,omitempty" yaml:"dataBits,omitempty"`
	// StopBits is the number of bits after each character
	StopBits int `json:"stopBits,omitempty" yaml:"stopBits,omitempty"`
	// Parity is whether parity error checking is used: none, even or odd
	// (see https://github.com/xylo04/goHamlib/blob/3752aec70bb9298eedaed8a58e834ffd92261ce0/goHamlib.go#L15-L20)
	Parity byte `json:"parity,omitempty" yaml:"parity,omitempty" jsonschema:"enum=0,enum=1,enum=2"`
	// Handshake is the handshake method used: none or RTSCTS
	// (see https://github.com/xylo04/goHamlib/blob/3752aec70bb9298eedaed8a58e834ffd92261ce0/goHamlib.go#L25-L29)
	Handshake byte `json:"handshake,omitempty" yaml:"handshake,omitempty" jsonschema:"enum=0,enum=1"`
}
