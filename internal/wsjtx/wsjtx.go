package wsjtx

import (
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

// HandleWsjtx is a goroutine that listens for WSJT-X messages and puts them on the given channel.
func HandleWsjtx(conf config.Config, msgChan chan Message) {
	ipAddr := net.ParseIP(conf.Wsjtx.Address)
	if ipAddr == nil {
		log.Error().Str("address", conf.Wsjtx.Address).Msg("couldn't parse WSJT-X IP address")
		return
	}
	log.Info().Msgf("Listening to WSJT-X at %v:%d UDP", ipAddr, conf.Wsjtx.Port)
	wsjtServ, err := wsjtx.MakeServerGiven(ipAddr, conf.Wsjtx.Port)
	if err != nil {
		log.Error().Err(err).Msg("couldn't listen to WSJT-X")
		return
	}
	wsjtChan := make(chan interface{}, 5)
	errChan := make(chan error, 5)
	go wsjtServ.ListenToWsjtx(wsjtChan, errChan)

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
