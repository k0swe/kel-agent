package integration

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func (s *integrationTestSuite) TestSend() {
	s.primeConnection()
	tests := []string{
		"clear",
		"close",
		"configure",
		"freeText",
		"haltTx",
		"heartbeat",
		"highlightCallsign",
		"location",
		"replay",
		"reply",
		"switchConfiguration",
	}

	for _, tt := range tests {
		s.T().Run(tt, func(t *testing.T) {
			input, _ := os.ReadFile(fmt.Sprintf("send/%s.json", tt))
			want, _ := os.ReadFile(fmt.Sprintf("send/%s.bin", tt))
			err := s.wsClient.WriteMessage(websocket.TextMessage, input)
			s.Require().NoError(err)

			select {
			case got := <-s.fake.ReceiveChan:
				s.Require().NoError(err)
				s.Require().Equal(want, got)
			case <-time.After(500 * time.Millisecond):
				t.Log("timeout")
				t.Fail()
			}
		})
	}
}

func (s *integrationTestSuite) primeConnection() {
	// Because this is UDP, the server doesn't have an address for WSJTX until WSJTX has sent the
	// server a message.
	clearMsg, _ := hex.DecodeString(`adbccbda00000002000000030000000657534a542d58`)
	_, err := s.fake.SendMessage(clearMsg)
	s.Require().NoError(err)
	_, _ = s.fake.SendMessage(clearMsg)
	s.T().Log("connection is primed for a send test")
}
