package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"github.com/k0swe/wsjtx-go"
	"log"
	"net/http"
	"reflect"
)

const defaultAddr = "localhost:8081"

var allowedOrigins sliceFlag = []string{
	"https://log.k0swe.radio",
	"http://localhost:8080",
	"http://localhost:4200",
}
var debug *bool

func main() {
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
	http.HandleFunc("/websocket", websocketHandler)
	http.HandleFunc("/", indexHandler)
	log.Printf("kel-agent ready to serve at %s%s", protocol, *addr)
	if *debug {
		log.Println("Verbose output enabled")
	}
	if secure {
		log.Fatal(http.ListenAndServeTLS(*addr, *cert, *key, nil))
	} else {
		log.Fatal(http.ListenAndServe(*addr, nil))
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Congratulations, you've reached kel-agent! " +
		"If you can see this, you should be able to connect to the websocket."))
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

	wsjtServ := wsjtx.MakeServer()
	wsjtChan := make(chan interface{}, 5)
	go wsjtServ.ListenToWsjtx(wsjtChan)

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
