package main

import (
	"encoding/json"
	"github.com/k0swe/wsjtx-go"
	"log"
	"reflect"
)

type WsjtxMessage struct {
	MsgType string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func handleWsjtx(msgChan chan []byte) {
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
		msgChan <- message
	}
}
