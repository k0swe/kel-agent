package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/k0swe/kel-agent/internal/server"
)

func (s *integrationTestSuite) TestReceive() {
	tests := []string{"heartbeat"}

	for _, tt := range tests {
		s.T().Run(tt, func(t *testing.T) {
			input, _ := ioutil.ReadFile(fmt.Sprintf("receive/%s.bin", tt))
			want, _ := ioutil.ReadFile(fmt.Sprintf("receive/%s.json", tt))
			_, _ = s.fake.SendMessage(input)

			_, got, err := s.wsClient.ReadMessage()
			s.Require().NoError(err)
			wantObj := json.Unmarshal(want, &server.WebsocketMessage{})
			gotObj := json.Unmarshal(got, &server.WebsocketMessage{})
			s.Require().Equal(wantObj, gotObj)
		})
	}
}
