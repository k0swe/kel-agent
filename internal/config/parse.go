package config

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"dario.cat/mergo"
	"github.com/adrg/xdg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	appName = "kel-agent"
)

var defaultConfigFile *string
var configFile *string

func init() {
	c, _ := xdg.ConfigFile(filepath.Join(appName, "config.yaml"))
	c = filepath.Clean(c)
	defaultConfigFile = &c
}

func ParseAllConfigs(schemaUrl string) (Config, error) {
	// Flags take precedence, and need to be parsed first for logging level
	conf, err := parseFlags()
	if err != nil {
		return Config{}, err
	}
	file, err := parseConfigFile(schemaUrl)
	if err != nil {
		return Config{}, err
	}
	if err := mergo.Merge(&conf, file); err != nil {
		return Config{}, err
	}
	if err := mergo.Merge(&conf, defaultConf); err != nil {
		return Config{}, err
	}
	log.Debug().Object("config", conf).Msgf("coalesced configuration")
	log.Info().Strs("origins", conf.Websocket.AllowedOrigins).Msg("allowed origins")
	return conf, nil
}

func parseFlags() (Config, error) {
	var conf = Config{}

	flag.StringVar(&conf.Websocket.Address, "host", conf.Websocket.Address, "websocket hosting address")
	flag.UintVar(&conf.Websocket.Port, "port", conf.Websocket.Port, "websocket hosting port")
	flag.StringVar(&conf.Websocket.Key, "key", conf.Websocket.Key, "TLS key")
	flag.StringVar(&conf.Websocket.Cert, "cert", conf.Websocket.Cert, "TLS certificate")
	var origins sliceFlag
	flag.Var(&origins, "origins", "comma-separated list of allowed origins")

	configFile = flag.String("config", *defaultConfigFile, "path to the configuration file")
	debug := flag.Bool("v", false, "verbose debugging output")
	trace := flag.Bool("vv", false, "trace debugging output")

	flag.Parse()
	conf.Websocket.AllowedOrigins = origins
	c, _ := filepath.Abs(*configFile)
	configFile = &c

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

	log.Trace().Object("flags", conf).Msgf("flag config")
	return conf, nil
}

func parseConfigFile(schemaUrl string) (Config, error) {
	var conf Config
	dat, err := os.ReadFile(*configFile)
	if os.IsNotExist(err) {
		log.Debug().Msgf("no config file found at '%s'", *configFile)
		if dat, err = yaml.Marshal(defaultConf); err != nil {
			return Config{}, err
		}
		// add schema to top of file
		dat = append([]byte("# yaml-language-server: $schema="+schemaUrl+"\n"), dat...)
		if err := os.WriteFile(*configFile, dat, 0o755); err != nil {
			return Config{}, err
		}
		log.Debug().Msgf("wrote default config to '%s'", *configFile)
		return defaultConf, nil
	}
	if err != nil {
		return Config{}, err
	}
	if err := yaml.Unmarshal(dat, &conf); err != nil {
		return Config{}, err
	}
	log.Debug().Str("path", *configFile).Msg("using config file")
	log.Trace().Object("file", conf).Msg("file config")
	return conf, nil
}
