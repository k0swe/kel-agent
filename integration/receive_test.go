package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/k0swe/kel-agent/internal/ws"
)

func (s *integrationTestSuite) TestReceive() {
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

			_, got, err := s.wsClient.ReadMessage()
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
