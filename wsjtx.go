package main

import (
	"github.com/k0swe/wsjtx-go/v2"
	"log"
	"reflect"
)

type WsjtxMessage struct {
	MsgType string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func handleWsjtx(msgChan chan WsjtxMessage) {
	wsjtServ, _ := wsjtx.MakeServer()
	wsjtChan := make(chan interface{}, 5)
	errChan := make(chan error, 5)
	go wsjtServ.ListenToWsjtx(wsjtChan, errChan)

	for {
		select {
		case wsjtMsg := <-wsjtChan:
			if *debug {
				log.Println("Sending wsjtx message:", wsjtMsg)
			}
			msgChan <- WsjtxMessage{
				MsgType: reflect.TypeOf(wsjtMsg).Name(),
				Payload: wsjtMsg,
			}
		case err := <-errChan:
			if *debug {
				log.Println("error: ", err)
			}
		}
	}
}
