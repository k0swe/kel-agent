package main

import (
	rigctl "github.com/ftl/rigproxy/pkg/client"
	"log"
	"reflect"
	"time"
)

type HamlibMessage struct {
	MsgType string      `json:"type,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

type RigState struct {
	Model     string `json:"model"`
	Frequency int64  `json:"frequency"`
	Mode      string `json:"mode"`
	Width     int    `json:"passbandWidthHz"`
}

const maxWaitInterval = 10 * time.Second

var lastState = RigState{}
var websocketChannel chan HamlibMessage
var maximumWait *time.Ticker

func handleHamlib(msgChan chan HamlibMessage) {
	websocketChannel = msgChan
	conn, err := rigctl.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	err = conn.StartPolling(200*time.Millisecond, 100*time.Millisecond,
		rigctl.PollCommand(rigctl.OnFrequency(onFrequency)),
		rigctl.PollCommand(rigctl.OnModeAndPassband(onModeAndPassband)))
	if err != nil {
		log.Fatal(err)
	}

	// If the poller has no changes after maxWaitInterval, send a heartbeat
	maximumWait = time.NewTicker(maxWaitInterval)
	for range maximumWait.C {
		sendState()
	}
}

func onFrequency(f rigctl.Frequency) {
	frequency := int64(f)
	if lastState.Frequency != frequency {
		lastState.Frequency = frequency
		sendState()
	}
}

func onModeAndPassband(m rigctl.Mode, f rigctl.Frequency) {
	mode := string(m)
	width := int(f)
	if lastState.Mode != mode || lastState.Width != width {
		lastState.Mode = mode
		lastState.Width = width
		sendState()
	}
}

func sendState() {
	websocketChannel <- HamlibMessage{
		MsgType: reflect.TypeOf(lastState).Name(),
		Payload: lastState,
	}
	maximumWait.Reset(maxWaitInterval)
}
