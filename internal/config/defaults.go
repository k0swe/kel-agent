package config

var defaultConf = Config{
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
		Address: "224.0.0.1",
		Port:    2237,
	},
}
