package config

import "runtime"

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
		Rigctld: RigctldConfig{
			Enabled: false,
			Address: "localhost",
			Port:    4532,
		},
	}
}
