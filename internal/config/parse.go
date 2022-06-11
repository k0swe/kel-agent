package config

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/adrg/xdg"
	"github.com/imdario/mergo"
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

	log.Trace().Msgf("flag config: %v", conf)
	return conf, nil
}

func parseConfigFile() (Config, error) {
	var conf Config
	dat, err := ioutil.ReadFile(*configFile)
	if os.IsNotExist(err) {
		log.Debug().Msgf("no config file found at '%s'", *configFile)
		if dat, err = yaml.Marshal(defaultConf); err != nil {
			return Config{}, err
		}
		if err := ioutil.WriteFile(*configFile, dat, 0o755); err != nil {
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
	log.Debug().Msgf("using config file at '%s'", *configFile)
	log.Trace().Msgf("file config: %v", conf)
	return conf, nil
}
