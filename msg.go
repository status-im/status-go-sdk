package sdk

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/crypto/sha3"
)

var (
	// NewContactKeyType message type for newContactKeyFormat
	NewContactKeyType = "~#c1"
	// ContactRequestType message type for contactRequestFormat
	ContactRequestType = "~#c2"
	// ConfirmedContactRequestType message type for confirmedContactRequestFormat
	ConfirmedContactRequestType = "~#c3"
	// StandardMessageType message type for StandardMessageFormat
	StandardMessageType = "~#c4"
	// SeenType message type for SeentType
	SeenType = "~#c5"
	// ContactUpdateType message type for contactUpdateMsg
	ContactUpdateType = "~#c6"
)

// supportedMessage check if the message type is supported
func supportedMessage(msgType string) bool {
	_, ok := map[string]bool{
		NewContactKeyType:           true,
		ContactRequestType:          true,
		ConfirmedContactRequestType: true,
		StandardMessageType:         true,
		SeenType:                    true,
		ContactUpdateType:           true,
	}[msgType]

	return ok
}

// Msg is a structure used by Subscribers and Publish().
type Msg struct {
	From        string `json:"from"`
	Text        string `json:"text"`
	ChannelName string `json:"channel"`
	Timestamp   int64  `json:"ts"`
	Raw         string `json:"-"`
	Type        string `json:"-"`
}

// NewMsg creates a new Msg with a generated UUID
func NewMsg(from, text, channel string) *Msg {
	return &Msg{
		From:        from,
		Text:        text,
		ChannelName: channel,
		Timestamp:   time.Now().Unix(),
	}
}

// ID gets the message id
func (m *Msg) ID() string {
	return fmt.Sprintf("%X", sha3.Sum256([]byte(m.Raw)))
}

// ToPayload converts current struct to a valid payload
func (m *Msg) ToPayload() string {
	message := fmt.Sprintf(messagePayloadMsg,
		m.Text,
		m.Timestamp*100,
		m.Timestamp)
	println(message)

	return rawrChatMessage(message)
}

func rawrChatMessage(raw string) string {
	bytes := []byte(raw)

	return fmt.Sprintf("0x%s", hex.EncodeToString(bytes))
}

func unrawrChatMessage(message string) ([]byte, error) {
	return hex.DecodeString(message[2:])
}

// MessageFromPayload creates a message from a payload
func MessageFromPayload(payload string) (*Msg, error) {
	message, err := unrawrChatMessage(payload)
	if err != nil {
		return nil, err
	}
	var x []interface{}
	if err := json.Unmarshal(message, &x); err != nil {
		return nil, err
	}

	if len(x) < 1 {
		return nil, errors.New("unsupported message type")
	}
	if x[0].(string) != "~#c4" {
		return nil, errors.New("unsupported message type")
	}
	properties := x[1].([]interface{})

	return &Msg{
		// TODO (adriacidre) add from username
		From:      "TODO : someone",
		Text:      properties[0].(string),
		Timestamp: int64(properties[3].(float64)),
		Raw:       string(message),
	}, nil
}
