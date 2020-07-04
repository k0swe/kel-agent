package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"github.com/xylo04/wsjtx-go/wsjtx"
	"log"
	"net/http"
	"strings"
)

const addr = "localhost:8081"

var allowedOrigins sliceFlag = []string{"http://localhost:8080"}

func main() {
	flag.Var(&allowedOrigins, "origins", "comma-separated list of allowed origins")
	flag.Parse()
	log.Println("Allowed origins are", allowedOrigins)

	http.HandleFunc("/websocket", websocketHandler)
	log.Print("k0s-agent ready to serve at http://", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

var upgrader = websocket.Upgrader{}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = logbookCheckOrigin
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	defer ws.Close()
	log.Println("Established websocket session with", r.RemoteAddr)

	c := make(chan interface{}, 5)
	go wsjtx.ListenToWsjtx(c)

	for {
		message, _ := json.Marshal(<-c)
		_ = ws.WriteMessage(websocket.TextMessage, []byte(message))
	}
}

func logbookCheckOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}
	log.Println("Rejecting websocket request from origin", origin)
	return false
}

type sliceFlag []string

func (i *sliceFlag) String() string {
	return "my string representation"
}

func (i *sliceFlag) Set(value string) error {
	tokens := strings.Split(value, ",")
	for _, t := range tokens {
		*i = append(*i, t)
	}
	return nil
}
