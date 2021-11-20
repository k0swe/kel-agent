package integration

import (
	"net"
	"strconv"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/k0swe/kel-agent/internal/config"
	"github.com/k0swe/kel-agent/internal/server"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

const origin = "https://test.example"

type integrationTestSuite struct {
	suite.Suite
	conf     config.Config
	wsClient *websocket.Conn
	server   *server.Server
	fake     *WsjtxFake
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, &integrationTestSuite{})
}

func (s *integrationTestSuite) SetupSuite() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	s.conf = config.Config{
		Websocket: config.WebsocketConfig{
			Address:        "127.0.0.1",
			Port:           8081, // TODO: use OS-assigned port
			AllowedOrigins: []string{origin},
		},
		Wsjtx: config.WsjtxConfig{
			Enabled: true,
			Address: "127.0.0.1",
			Port:    2237, // TODO: use OS-assigned port
		},
		VersionInfo: "kel-agent v0.0.0 (abcd)",
	}
	var err error
	s.server, err = server.Start(s.conf)
	s.Require().NoError(err)
	<-s.server.Started

	wsAddr := net.JoinHostPort(s.conf.Websocket.Address, strconv.Itoa(int(s.conf.Websocket.Port)))
	wsAddr = "ws://" + wsAddr + "/websocket"
	header := map[string][]string{"Origin": {origin}}
	s.wsClient, _, err = websocket.DefaultDialer.Dial(wsAddr, header)
	s.Require().NoError(err)
}

func (s *integrationTestSuite) SetupTest() {
	var err error
	s.fake, err = NewFake(&net.UDPAddr{Port: int(s.conf.Wsjtx.Port)}, s.T())
	s.Require().NoError(err)
	s.T().Log("suite reports fake is connected")
}

func (s *integrationTestSuite) TearDownTest() {
	s.fake.Stop()
}