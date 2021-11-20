package server

import (
	"fmt"
	"net/http"

	"github.com/k0swe/kel-agent/internal/config"
	"github.com/rs/zerolog/log"
)

type Server struct {
	conf config.Config
	hub  *Hub
	Stop chan bool
}

func Start(c config.Config) (*Server, error) {
	if c.Websocket.Key != "" && c.Websocket.Cert == "" ||
		c.Websocket.Key == "" && c.Websocket.Cert != "" {
		return &Server{}, fmt.Errorf("-key and -cert must be used together")
	}

	hub := newHub(&c)
	go hub.run()
	server := Server{c, hub, make(chan bool, 1)}

	secure := false
	protocol := "ws://"
	if c.Websocket.Key != "" && c.Websocket.Cert != "" {
		secure = true
		protocol = "wss://"
	}
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		server.serveWs(hub, w, r)
	})
	http.HandleFunc("/", server.indexHandler)
	addrAndPort := fmt.Sprintf("%s:%d", c.Websocket.Address, c.Websocket.Port)
	log.Info().Msgf("ready to serve at %s%s", protocol, addrAndPort)
	if secure {
		go func() {
			log.Fatal().Err(
				http.ListenAndServeTLS(addrAndPort, c.Websocket.Cert, c.Websocket.Key, nil)).Msg("websocket dying")
			server.Stop <- true
		}()
	} else {
		go func() {
			log.Fatal().Err(http.ListenAndServe(addrAndPort, nil)).Msg("websocket dying")
			server.Stop <- true
		}()
	}
	return &server, nil
}

func (s Server) indexHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Congratulations, you've reached kel-agent! " +
		"If you can see this, you should be able to connect to the websocket."))
}
