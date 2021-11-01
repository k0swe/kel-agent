package wsjtx

import (
	"reflect"
	"strconv"

	"github.com/k0swe/kel-agent/internal/config"
	"github.com/k0swe/wsjtx-go/v2"
	"github.com/rs/zerolog/log"
)

type Message struct {
	MsgType string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// HandleWsjtx is a goroutine that listens for WSJT-X messages and puts them on the given channel.
func HandleWsjtx(conf config.Config, msgChan chan Message) {
	log.Info().Msgf("Listening to WSJT-X at %s:%d UDP", conf.Wsjtx.Address, conf.Wsjtx.Port)
	wsjtServ, _ := wsjtx.MakeMulticastServer(
		conf.Wsjtx.Address, strconv.Itoa(int(conf.Wsjtx.Port)))
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
