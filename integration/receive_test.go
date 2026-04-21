package integration

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/k0swe/kel-agent/internal/ws"
)

func (s *integrationTestSuite) TestReceive() {
	s.waitForPipelineReady()

	tests := []string{
		"clear",
		"close",
		"decode",
		"heartbeat",
		"logged-adif",
		"qso-logged",
		"status-222",
		"status-231",
		"wspr-decode",
	}

	for _, tt := range tests {
		s.T().Run(tt, func(t *testing.T) {
			input, _ := os.ReadFile(fmt.Sprintf("receive/%s.bin", tt))
			want, _ := os.ReadFile(fmt.Sprintf("receive/%s.json", tt))
			_, _ = s.fake.SendMessage(input)

			_ = s.wsClient.SetReadDeadline(time.Now().Add(5 * time.Second))
			_, got, err := s.wsClient.ReadMessage()
			_ = s.wsClient.SetReadDeadline(time.Time{})
			s.Require().NoError(err)

			wantObj := &ws.WebsocketMessage{}
			err = json.Unmarshal(want, &wantObj)
			s.Require().NoError(err)
			gotObj := &ws.WebsocketMessage{}
			err = json.Unmarshal(got, &gotObj)
			s.Require().NoError(err)
			s.Require().Equal(wantObj, gotObj)
		})
	}
}

// waitForPipelineReady ensures the WSJTX→hub→websocket message pipeline is
// working before the first subtest assertion runs. Under slow or
// QEMU-emulated environments the pipeline may not have forwarded the first
// UDP datagram yet, causing ReadMessage to block indefinitely.
func (s *integrationTestSuite) waitForPipelineReady() {
	// Use the same clear message used by primeConnection in send_test.go.
	clearMsg, err := hex.DecodeString(`adbccbda00000002000000030000000657534a542d58`)
	s.Require().NoError(err)
	const maxRetries = 10
	for i := 0; i < maxRetries; i++ {
		_, _ = s.fake.SendMessage(clearMsg)
		_ = s.wsClient.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		_, _, err := s.wsClient.ReadMessage()
		_ = s.wsClient.SetReadDeadline(time.Time{})
		if err == nil {
			s.T().Logf("pipeline ready after %d probe(s)", i+1)
			return
		}
	}
	s.Require().Fail("WSJTX→websocket pipeline did not become ready in time")
}
