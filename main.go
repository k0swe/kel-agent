package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"github.com/xylo04/wsjtx-go/wsjtx"
	"log"
	"net/http"
	"reflect"
)

const addr = "localhost:8081"

var allowedOrigins sliceFlag = []string{"http://localhost:8080", "http://localhost:4200"}
var debug *bool

func main() {
	flag.Var(&allowedOrigins, "origins", "comma-separated list of allowed origins")
	debug = flag.Bool("v", false, "Verbose debugging output")
	flag.Parse()
	log.Println("Allowed origins are", allowedOrigins)

	http.HandleFunc("/websocket", websocketHandler)
	log.Print("k0s-agent ready to serve at http://", addr)
	if *debug {
		log.Println("Verbose output enabled")
	}
	log.Fatal(http.ListenAndServe(addr, nil))
}

var upgrader = websocket.Upgrader{}

type WebsocketMessage struct {
	Wsjtx WsjtxMessage `json:"wsjtx"`
}

type WsjtxMessage struct {
	MsgType string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = logbookCheckOrigin
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	defer ws.Close()
	log.Println("Established websocket session with", r.RemoteAddr)

	wsjtChan := make(chan interface{}, 5)
	go wsjtx.ListenToWsjtx(wsjtChan)

	for {
		wsjtMsg := <-wsjtChan
		if *debug {
			log.Println("Sending wsjtx message:", wsjtMsg)
		}
		wsMsg := WebsocketMessage{Wsjtx: WsjtxMessage{
			MsgType: reflect.TypeOf(wsjtMsg).Name(),
			Payload: wsjtMsg,
		}}
		message, _ := json.Marshal(wsMsg)
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
