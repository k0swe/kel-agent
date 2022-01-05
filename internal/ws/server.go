package ws

import (
	"fmt"
	"net/http"

	"github.com/k0swe/kel-agent/internal/config"
	"github.com/rs/zerolog/log"
)

type Server struct {
	conf    config.Config
	hub     *Hub
	Started chan bool
	Stop    chan bool
}

const wsPath = "/websocket"

func Start(c config.Config) (*Server, error) {
	if c.Websocket.Key != "" && c.Websocket.Cert == "" ||
		c.Websocket.Key == "" && c.Websocket.Cert != "" {
		return &Server{}, fmt.Errorf("-key and -cert must be used together")
	}

	hub := newHub(&c)
	go hub.run()
	server := Server{c, hub, make(chan bool, 1), make(chan bool, 1)}

	secure := false
	protocol := "ws://"
	if c.Websocket.Key != "" && c.Websocket.Cert != "" {
		secure = true
		protocol = "wss://"
	}
	mux := http.NewServeMux()
	mux.HandleFunc(wsPath, func(w http.ResponseWriter, r *http.Request) {
		server.serveWs(hub, w, r)
	})
	mux.HandleFunc("/", server.indexHandler)
	addrAndPort := fmt.Sprintf("%s:%d", c.Websocket.Address, c.Websocket.Port)
	log.Info().Str("address", fmt.Sprintf("%s%s%s", protocol, addrAndPort, wsPath)).Msg("Serving websocket")
	go func() {
		var err error
		server.Started <- true
		if secure {
			err = http.ListenAndServeTLS(addrAndPort, c.Websocket.Cert, c.Websocket.Key, mux)
		} else {
			err = http.ListenAndServe(addrAndPort, mux)
		}
		log.Fatal().Err(err).Msg("websocket dying")
		server.Stop <- true
	}()
	return &server, nil
}

func (s Server) indexHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Congratulations, you've reached kel-agent! " +
		"If you can see this, you should be able to connect to the websocket."))
}
