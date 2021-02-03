package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

const defaultAddr = "localhost:8081"

var allowedOrigins sliceFlag = []string{
	"https://forester.radio",
	"http://localhost:8080",
	"http://localhost:4200",
}
var debug *bool

func main() {
	log.Printf("kel-agent %v (%v) %v %v %v %v",
		Version, GitCommit, runtime.Version(), runtime.GOOS, runtime.GOARCH, BuildTime)
	flag.Var(&allowedOrigins, "origins", "comma-separated list of allowed origins")
	debug = flag.Bool("v", false, "verbose debugging output")
	addr := flag.String("host", defaultAddr, "hosting address")
	key := flag.String("key", "", "TLS key")
	cert := flag.String("cert", "", "TLS certificate")
	flag.Parse()
	if *key != "" && *cert == "" || *key == "" && *cert != "" {
		panic("-key and -cert must be used together")
	}
	secure := false
	protocol := "ws://"
	if *key != "" && *cert != "" {
		secure = true
		protocol = "wss://"
	}
	log.Println("Allowed origins are", allowedOrigins)

	hub := newHub(fmt.Sprintf("kel-agent %v (%v)", Version, GitCommit))
	go hub.run()
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.HandleFunc("/", indexHandler)
	log.Printf("ready to serve at %s%s", protocol, *addr)
	if *debug {
		log.Println("Verbose output enabled")
	}
	if secure {
		log.Fatal(http.ListenAndServeTLS(*addr, *cert, *key, nil))
	} else {
		log.Fatal(http.ListenAndServe(*addr, nil))
	}
}

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Congratulations, you've reached kel-agent! " +
		"If you can see this, you should be able to connect to the websocket."))
}
