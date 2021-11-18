package wsjtx_wrapper

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"

	"github.com/k0swe/kel-agent/internal/config"
	"github.com/k0swe/wsjtx-go/v3"
	"github.com/rs/zerolog/log"
)

type Message struct {
	MsgType string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Handler struct {
	wsjtxServ wsjtx.Server
	conf      config.Config
}

func NewHandler(c config.Config) (*Handler, error) {
	ipAddr := net.ParseIP(c.Wsjtx.Address)
	if ipAddr == nil {
		log.Error().Str("address", c.Wsjtx.Address).Msg("couldn't parse WSJT-X IP address")
		return nil, fmt.Errorf("couldn't parse WSJT-X IP address %s", c.Wsjtx.Address)
	}
	log.Info().Msgf("Listening to WSJT-X at %v:%d UDP", ipAddr, c.Wsjtx.Port)
	var err error
	serv, err := wsjtx.MakeServerGiven(ipAddr, c.Wsjtx.Port)
	if err != nil {
		log.Error().Err(err).Msg("couldn't listen to WSJT-X")
		return nil, fmt.Errorf("couldn't listen to WSJT-X: %s", err)
	}
	return &Handler{
		wsjtxServ: serv,
		conf:      c,
	}, nil
}

// ListenToWsjtx is a goroutine that listens for WSJT-X messages and puts them on the given channel.
func (h *Handler) ListenToWsjtx(msgChan chan Message) {
	defer func() { h.wsjtxServ = wsjtx.Server{} }()
	wsjtChan := make(chan interface{}, 5)
	errChan := make(chan error, 5)
	go h.wsjtxServ.ListenToWsjtx(wsjtChan, errChan)

	for {
		select {
		case wsjtMsg := <-wsjtChan:
			log.Trace().Msgf("Received message from wsjtx: %v", wsjtMsg)
			msgChan <- Message{
				MsgType: reflect.TypeOf(wsjtMsg).Name(),
				Payload: wsjtMsg,
			}
		case err := <-errChan:
			log.Debug().Err(err).Msgf("wsjtx error")
		}
	}
}

func (h *Handler) HandleClientCommand(msgType string, payload []byte) error {
	switch msgType {
	case reflect.TypeOf(wsjtx.HeartbeatMessage{}).Name():
		var heartbeatMsg = &wsjtx.HeartbeatMessage{}
		err := json.Unmarshal(payload, heartbeatMsg)
		if err != nil {
			return err
		}
		return h.wsjtxServ.Heartbeat(*heartbeatMsg)
	case reflect.TypeOf(wsjtx.ClearMessage{}).Name():
		var clearMsg = &wsjtx.ClearMessage{}
		err := json.Unmarshal(payload, clearMsg)
		if err != nil {
			return err
		}
		return h.wsjtxServ.Clear(*clearMsg)
	case reflect.TypeOf(wsjtx.ReplyMessage{}).Name():
		var replyMsg = &wsjtx.ReplyMessage{}
		err := json.Unmarshal(payload, replyMsg)
		if err != nil {
			return err
		}
		return h.wsjtxServ.Reply(*replyMsg)
	case reflect.TypeOf(wsjtx.CloseMessage{}).Name():
		var closeMsg = &wsjtx.CloseMessage{}
		err := json.Unmarshal(payload, closeMsg)
		if err != nil {
			return err
		}
		return h.wsjtxServ.Close(*closeMsg)
	case reflect.TypeOf(wsjtx.ReplayMessage{}).Name():
		var replayMsg = &wsjtx.ReplayMessage{}
		err := json.Unmarshal(payload, replayMsg)
		if err != nil {
			return err
		}
		return h.wsjtxServ.Replay(*replayMsg)
	case reflect.TypeOf(wsjtx.HaltTxMessage{}).Name():
		var haltMsg = &wsjtx.HaltTxMessage{}
		err := json.Unmarshal(payload, haltMsg)
		if err != nil {
			return err
		}
		return h.wsjtxServ.HaltTx(*haltMsg)
	case reflect.TypeOf(wsjtx.FreeTextMessage{}).Name():
		var freeTextMsg = &wsjtx.FreeTextMessage{}
		err := json.Unmarshal(payload, freeTextMsg)
		if err != nil {
			return err
		}
		return h.wsjtxServ.FreeText(*freeTextMsg)
	case reflect.TypeOf(wsjtx.LocationMessage{}).Name():
		var locationMsg = &wsjtx.LocationMessage{}
		err := json.Unmarshal(payload, locationMsg)
		if err != nil {
			return err
		}
		return h.wsjtxServ.Location(*locationMsg)
	case reflect.TypeOf(wsjtx.HighlightCallsignMessage{}).Name():
		var highlightMsg = &wsjtx.HighlightCallsignMessage{}
		err := json.Unmarshal(payload, highlightMsg)
		if err != nil {
			return err
		}
		return h.wsjtxServ.HighlightCallsign(*highlightMsg)
	case reflect.TypeOf(wsjtx.SwitchConfigurationMessage{}).Name():
		var switchConfigMsg = &wsjtx.SwitchConfigurationMessage{}
		err := json.Unmarshal(payload, switchConfigMsg)
		if err != nil {
			return err
		}
		return h.wsjtxServ.SwitchConfiguration(*switchConfigMsg)
	case reflect.TypeOf(wsjtx.ConfigureMessage{}).Name():
		var configMsg = &wsjtx.ConfigureMessage{}
		err := json.Unmarshal(payload, configMsg)
		if err != nil {
			return err
		}
		return h.wsjtxServ.Configure(*configMsg)

	default:
		return fmt.Errorf("implemented wsjtx message type %s", msgType)
	}
}
