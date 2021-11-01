package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/k0swe/kel-agent/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var versionInfo string
var conf config.Config

func main() {
	versionInfo = fmt.Sprintf("kel-agent %v (%v)", Version, GitCommit)
	fmt.Printf("%v %v %v %v %v\n",
		versionInfo, runtime.Version(), runtime.GOOS, runtime.GOARCH, BuildTime)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	conf = config.ParseAllConfigs()
	log.Debug().Msg("Verbose output enabled")

	if conf.Websocket.Key != "" && conf.Websocket.Cert == "" ||
		conf.Websocket.Key == "" && conf.Websocket.Cert != "" {
		panic("-key and -cert must be used together")
	}
	secure := false
	protocol := "ws://"
	if conf.Websocket.Key != "" && conf.Websocket.Cert != "" {
		secure = true
		protocol = "wss://"
	}

	hub := newHub()
	go hub.run()
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.HandleFunc("/", indexHandler)
	addrAndPort := fmt.Sprintf("%s:%d", conf.Websocket.Address, conf.Websocket.Port)
	log.Info().Msgf("ready to serve at %s%s", protocol, addrAndPort)
	if secure {
		log.Fatal().Err(
			http.ListenAndServeTLS(addrAndPort, conf.Websocket.Cert, conf.Websocket.Key, nil)).Msg("")
	} else {
		log.Fatal().Err(http.ListenAndServe(addrAndPort, nil)).Msg("")
	}
}

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Congratulations, you've reached kel-agent! " +
		"If you can see this, you should be able to connect to the websocket."))
}
