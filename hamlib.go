package main

import (
	"github.com/dh1tw/goHamlib"
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

const pollInterval = 500 * time.Millisecond
const maxWaitInterval = 10 * time.Second

var lastState RigState
var websocketChannel chan HamlibMessage
var poller *time.Ticker
var maximumWait *time.Ticker

func handleHamlib(msgChan chan HamlibMessage) {
	websocketChannel = msgChan
	rig := goHamlib.Rig{}
	goHamlib.SetDebugLevel(goHamlib.DebugNone)
	if err := rig.Init(214); err != nil {
		panic(err)
	}
	if err := rig.SetPort(goHamlib.Port{
		RigPortType: goHamlib.RigPortSerial,
		Portname:    "/dev/ttyUSB0",
		Baudrate:    9600,
		Databits:    8,
		Stopbits:    1,
		Parity:      goHamlib.ParityNone,
		Handshake:   goHamlib.HandshakeNone,
	}); err != nil {
		panic(err)
	}
	if err := rig.Open(); err != nil {
		panic(err)
	}
	poller = time.NewTicker(pollInterval)
	maximumWait = time.NewTicker(maxWaitInterval)

	for {
		select {
		case <-poller.C:
			state := getCurrentState(rig)
			if state != lastState {
				sendState(state)
			}
		case <-maximumWait.C:
			state := getCurrentState(rig)
			sendState(state)
		}
	}
}

func sendState(state RigState) {
	websocketChannel <- HamlibMessage{
		MsgType: reflect.TypeOf(state).Name(),
		Payload: state,
	}
	lastState = state
	poller.Reset(pollInterval)
	maximumWait.Reset(maxWaitInterval)
}

func getCurrentState(rig goHamlib.Rig) RigState {
	freq, err := rig.GetFreq(goHamlib.VFOCurrent)
	if err != nil {
		panic(err)
	}
	mode, width, err := rig.GetMode(goHamlib.VFOCurrent)
	state := RigState{
		Model:     rig.Caps.ModelName,
		Frequency: int64(freq),
		Mode:      goHamlib.ModeName[mode],
		Width:     width,
	}
	return state
}
