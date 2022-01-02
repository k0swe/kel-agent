package webview

import (
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"

	"github.com/k0swe/kel-agent/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/webview/webview"
)

var conf *config.Config

//go:embed kel-agent-gui/dist/kel-agent-gui/*
var angularFS embed.FS

func StartWebView(c *config.Config) {
	conf = c
	fsContent, err := fs.Sub(fs.FS(angularFS), "kel-agent-gui/dist/kel-agent-gui")
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	angularFileServer := http.FileServer(http.FS(fsContent))
	mux := http.NewServeMux()
	mux.Handle("/", angularFileServer)
	// port 0: let OS choose port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	port := listener.Addr().(*net.TCPAddr).Port
	address := fmt.Sprintf("http://localhost:%d", port)
	log.Info().Str("address", address).Msg("Serving webview")
	go func() {
		if err := http.Serve(listener, mux); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()
	webViewStart(address)
}

func webViewStart(address string) {
	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("kel-agent")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate(address)
	w.Run()
}
