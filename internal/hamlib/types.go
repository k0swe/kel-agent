package hamlib

// Message is the envelope for Hamlib data sent over the websocket.
type Message struct {
	MsgType string      `json:"type,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

// RigState holds the current state of a rig as reported by Hamlib.
type RigState struct {
	Model     string `json:"model"`
	Frequency int64  `json:"frequency"`
	Mode      string `json:"mode"`
	Width     int    `json:"passbandWidthHz"`
}
