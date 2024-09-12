package rigctld_wrapper

import (
	"encoding/json"
	"fmt"
	"github.com/k0swe/kel-agent/internal/config"
	"github.com/k0swe/rigctld-go"
	"github.com/rs/zerolog/log"
	"net"
	"time"
)

type Handler struct {
	client rigctld.Client
	conf   config.Config
}

func NewHandler(c *config.Config) (*Handler, error) {
	ipAddr := net.ParseIP(c.Rigctld.Address)
	if ipAddr == nil {
		log.Error().Str("address", c.Rigctld.Address).Msg("couldn't parse rigctld IP address")
		return nil, fmt.Errorf("couldn't parse rigctld IP address %s", c.Rigctld.Address)
	}
	log.Info().
		Str("address", fmt.Sprintf("%v:%d", ipAddr, c.Rigctld.Port)).
		Msg("Connecting to rigctld on TCP")
	var err error
	client, err := rigctld.ConnectTo(ipAddr, c.Rigctld.Port)
	if err != nil {
		log.Error().Err(err).Msg("couldn't connect to rigctld")
		return nil, fmt.Errorf("couldn't connect to rigctld: %s", err)
	}
	client.SetReadDeadline(250 * time.Millisecond)
	h := &Handler{
		client: client,
		conf:   *c,
	}
	return h, nil
}

func (h Handler) HandleClientCommand(payload []byte) (interface{}, error) {
	var msg map[string]interface{}
	if err := json.Unmarshal(payload, &msg); err != nil {
		log.Warn().Err(err).Msg("failed to parse client command; dropping")
		return nil, err
	}
	log.Debug().Interface("msg", msg).Msg("handler received client command")

	if command, ok := msg["command"].(string); ok {
		switch command {
		case "getFrequency":
			return h.handleGetFreq()
		case "getMode":
			return h.handleGetMode()
		default:
			log.Warn().Str("command", command).Msg("unknown command")
		}
	} else {
		log.Warn().Msg("command not present")
	}
	return nil, nil
}

func (h Handler) handleGetFreq() (interface{}, error) {
	freq, err := h.client.GetFreq()
	if err != nil {
		log.Warn().Err(err).Msg("couldn't get frequency")
		resp := map[string]interface{}{
			"command": "getFrequency",
			"error":   fmt.Sprintf("%v", err),
		}
		return resp, err
	}
	resp := map[string]interface{}{
		"command":   "getFrequency",
		"frequency": freq,
	}
	return resp, nil
}

func (h Handler) handleGetMode() (interface{}, error) {
	mode, bandpass, err := h.client.GetMode()
	if err != nil {
		log.Warn().Err(err).Msg("couldn't get mode")
		resp := map[string]interface{}{
			"command": "getMode",
			"error":   fmt.Sprintf("%v", err),
		}
		return resp, err
	}
	resp := map[string]interface{}{
		"command":  "getMode",
		"mode":     mode.String(),
		"bandpass": bandpass,
	}
	return resp, nil
}
