//go:build !hamlib

package hamlib

import (
	"github.com/k0swe/kel-agent/internal/config"
	"github.com/rs/zerolog/log"
)

// HandleHamlib is a no-op stub when kel-agent is built without Hamlib support.
func HandleHamlib(conf *config.Config, msgChan chan Message) {
	log.Warn().Msg("Hamlib support is not compiled into this build; ignoring hamlib configuration")
}
