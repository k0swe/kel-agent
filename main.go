package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/k0swe/kel-agent/internal/config"
	"github.com/k0swe/kel-agent/internal/ws"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var versionInfo string

func main() {
	versionInfo = fmt.Sprintf("kel-agent %v (%v)", Version, GitCommit)
	fmt.Printf("%v %v %v %v %v\n",
		versionInfo, runtime.Version(), runtime.GOOS, runtime.GOARCH, BuildTime)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	c, err := config.ParseAllConfigs()
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't get configuration")
	}
	c.VersionInfo = versionInfo

	wsServer, err := ws.Start(c)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	<-wsServer.Stop
}
