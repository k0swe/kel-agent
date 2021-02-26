package config

import (
	"runtime"

	"github.com/xylo04/goHamlib"
)

var defaultConf Config

func init() {
	var defaultWsjtxAddr string
	switch runtime.GOOS {
	case "windows":
		defaultWsjtxAddr = "127.0.0.1"
	default:
		defaultWsjtxAddr = "224.0.0.1"
	}
	defaultConf = Config{
		Websocket: WebsocketConfig{
			Address: "localhost",
			Port:    8081,
			Key:     "",
			Cert:    "",
			AllowedOrigins: []string{
				"https://forester.radio",
			},
		},
		Wsjtx: WsjtxConfig{
			Enabled: true,
			Address: defaultWsjtxAddr,
			Port:    2237,
		},
		Hamlib: HamlibConfig{
			Enabled:      false,
			RetrySeconds: 10,
			RigModel:     3073, // Icom IC-7300
			RigPort:      goHamlib.RigPortName[goHamlib.RigPortSerial],
			PortName:     "/dev/ttyUSB0",
			BaudRate:     9600,
			DataBits:     8,
			StopBits:     1,
			Parity:       byte(goHamlib.ParityNone),
			Handshake:    byte(goHamlib.HandshakeNone),
		},
	}
}
