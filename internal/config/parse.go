package config

import (
	"flag"
	"strconv"
	"strings"

	"github.com/imdario/mergo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ParseAllConfigs() Config {
	conf := parseFlags()
	if err := mergo.Merge(&conf, defaultConf); err != nil {
		panic("problem merging default config values with flags")
	}
	// TODO: use this once I figure out why mergo is overwriting conf.LogLevel
	// zerolog.SetGlobalLevel(conf.LogLevel)
	log.Debug().Msgf("final configuration is %v", conf)
	return conf
}

func parseFlags() Config {
	var conf = Config{}

	flag.StringVar(&conf.Websocket.Address, "host", conf.Websocket.Address, "websocket address")
	flag.UintVar(&conf.Websocket.Port, "port", conf.Websocket.Port, "websocket port")
	flag.StringVar(&conf.Websocket.Key, "key", conf.Websocket.Key, "TLS key")
	flag.StringVar(&conf.Websocket.Cert, "cert", conf.Websocket.Cert, "TLS certificate")
	var origins sliceFlag
	flag.Var(&origins, "origins", "comma-separated list of allowed origins")

	debug := flag.Bool("v", false, "verbose debugging output")
	trace := flag.Bool("vv", false, "trace debugging output")

	flag.Parse()
	conf.Websocket.AllowedOrigins = origins

	switch {
	case *trace:
		conf.LogLevel = zerolog.TraceLevel
	case *debug:
		conf.LogLevel = zerolog.DebugLevel
	default:
		conf.LogLevel = zerolog.InfoLevel
	}
	// TODO: remove this
	zerolog.SetGlobalLevel(conf.LogLevel)

	// hosting address backward compat
	if i := strings.Index(conf.Websocket.Address, ":"); i >= 0 {
		port, _ := strconv.Atoi(conf.Websocket.Address[i+1:])
		conf.Websocket.Port = uint(port)
		conf.Websocket.Address = conf.Websocket.Address[:i]
	}

	return conf
}
