package config

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/adrg/xdg"
	"github.com/imdario/mergo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	appName = "kel-agent"
)

func ParseAllConfigs() (Config, error) {
	// Flags take precedence, and need to be parsed first for logging level
	conf, err := parseFlags()
	if err != nil {
		return Config{}, err
	}
	file, err := parseConfigFile()
	if err != nil {
		return Config{}, err
	}
	if err := mergo.Merge(&conf, file); err != nil {
		return Config{}, err
	}
	if err := mergo.Merge(&conf, defaultConf); err != nil {
		return Config{}, err
	}
	log.Debug().Msgf("effective configuration is %v", conf)
	return conf, nil
}

func parseFlags() (Config, error) {
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

	var ll zerolog.Level
	switch {
	case *trace:
		ll = zerolog.TraceLevel
		log.Trace().Msg("TRACE output enabled")
	case *debug:
		ll = zerolog.DebugLevel
		log.Debug().Msg("DEBUG output enabled")
	default:
		ll = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(ll)

	// hosting address backward compat
	if i := strings.Index(conf.Websocket.Address, ":"); i >= 0 {
		port, err := strconv.Atoi(conf.Websocket.Address[i+1:])
		if err != nil {
			return Config{}, err
		}
		conf.Websocket.Port = uint(port)
		conf.Websocket.Address = conf.Websocket.Address[:i]
	}

	log.Trace().Msgf("flag config: %v", conf)
	return conf, nil
}

func parseConfigFile() (Config, error) {
	path, err := xdg.ConfigFile(filepath.Join(appName, "config.yaml"))
	if err != nil {
		return Config{}, err
	}
	var conf Config
	dat, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		log.Debug().Msgf("no config file found at '%s'", path)
		if dat, err = yaml.Marshal(defaultConf); err != nil {
			return Config{}, err
		}
		if err := os.WriteFile(path, dat, 0o755); err != nil {
			return Config{}, err
		}
		log.Debug().Msgf("wrote default config to '%s'", path)
		return defaultConf, nil
	}
	if err != nil {
		return Config{}, err
	}
	if err := yaml.Unmarshal(dat, &conf); err != nil {
		return Config{}, err
	}
	log.Trace().Msgf("file config: %v", conf)
	return conf, nil
}
