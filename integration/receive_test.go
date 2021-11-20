package integration

import (
	"encoding/json"
	"os"

	"github.com/k0swe/kel-agent/internal/server"
)

func (s *integrationTestSuite) TestReceiveHeartbeat() {
	input, _ := os.ReadFile("receive/heartbeat.bin")
	want, _ := os.ReadFile("receive/heartbeat.json")
	_, _ = s.fake.SendMessage(input)

	_, got, err := s.wsClient.ReadMessage()
	s.Require().NoError(err)
	wantObj := json.Unmarshal(want, &server.WebsocketMessage{})
	gotObj := json.Unmarshal(got, &server.WebsocketMessage{})
	s.Require().Equal(wantObj, gotObj)
}
