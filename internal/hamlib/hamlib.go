package hamlib

import (
	"reflect"
	"time"

	"github.com/k0swe/kel-agent/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/xylo04/goHamlib"
)

type Message struct {
	MsgType string      `json:"type,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

type RigState struct {
	Model     string `json:"model"`
	Frequency int64  `json:"frequency"`
	Mode      string `json:"mode"`
	Width     int    `json:"passbandWidthHz"`
}

const pollInterval = 100 * time.Millisecond
const maxWaitInterval = 10 * time.Second

var lastState RigState
var websocketChannel chan Message
var poller *time.Ticker
var maximumWait *time.Ticker

func HandleHamlib(conf *config.Config, msgChan chan Message) {
	websocketChannel = msgChan
	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Info().
				Int("rigModel", conf.Hamlib.RigModel).
				Str("port", conf.Hamlib.PortName).
				Msg("connecting to hamlib")
			err := connectToRig(conf)
			if err != nil {
				log.Warn().
					Err(err).
					Msg("problem connecting to hamlib")
				ticker.Reset(time.Duration(conf.Hamlib.RetrySeconds) * time.Second)
				continue
			}
		}
	}
}

func connectToRig(conf *config.Config) error {
	rig := goHamlib.Rig{}
	goHamlib.SetDebugLevel(goHamlib.DebugNone)
	if err := rig.Init(goHamlib.RigModelID(conf.Hamlib.RigModel)); err != nil {
		return err
	}
	portConf := goHamlib.Port{
		RigPortType: goHamlib.RigPortValue[conf.Hamlib.RigPort],
		Portname:    conf.Hamlib.PortName,
		Baudrate:    conf.Hamlib.BaudRate,
		Databits:    conf.Hamlib.DataBits,
		Stopbits:    conf.Hamlib.StopBits,
		Parity:      goHamlib.Parity(conf.Hamlib.Parity),
		Handshake:   goHamlib.Handshake(conf.Hamlib.Handshake),
	}
	if err := rig.SetPort(portConf); err != nil {
		return err
	}
	if err := rig.Open(); err != nil {
		return err
	}
	defer sendDisconnect()
	log.Info().Msg("connected to hamlib")
	poller = time.NewTicker(pollInterval)
	maximumWait = time.NewTicker(maxWaitInterval)

	for {
		select {
		case <-poller.C:
			state, err := getCurrentState(rig)
			if err != nil {
				return err
			}
			if state != lastState {
				sendState(state)
			}
		case <-maximumWait.C:
			state, err := getCurrentState(rig)
			if err != nil {
				return err
			}
			sendState(state)
		}
	}
}

func sendDisconnect() {
	websocketChannel <- Message{
		MsgType: "disconnect",
		Payload: nil,
	}
}

func sendState(state RigState) {
	websocketChannel <- Message{
		MsgType: reflect.TypeOf(state).Name(),
		Payload: state,
	}
	lastState = state
	poller.Reset(pollInterval)
	maximumWait.Reset(maxWaitInterval)
}

func getCurrentState(rig goHamlib.Rig) (RigState, error) {
	freq, err := rig.GetFreq(goHamlib.VFOCurrent)
	if err != nil {
		return RigState{}, err
	}
	mode, width, err := rig.GetMode(goHamlib.VFOCurrent)
	state := RigState{
		Model:     rig.Caps.ModelName,
		Frequency: int64(freq),
		Mode:      goHamlib.ModeName[mode],
		Width:     width,
	}
	return state, err
}
