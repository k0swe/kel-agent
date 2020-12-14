package main

import (
	"github.com/k0swe/wsjtx-go"
	"log"
	"reflect"
)

type WsjtxMessage struct {
	MsgType string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func handleWsjtx(msgChan chan WsjtxMessage) {
	wsjtServ := wsjtx.MakeServer()
	wsjtChan := make(chan interface{}, 5)
	go wsjtServ.ListenToWsjtx(wsjtChan)

	for {
		wsjtMsg := <-wsjtChan
		if *debug {
			log.Println("Sending wsjtx message:", wsjtMsg)
		}
		msgChan <- WsjtxMessage{
			MsgType: reflect.TypeOf(wsjtMsg).Name(),
			Payload: wsjtMsg,
		}
	}
}
