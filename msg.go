package sdk

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

const (
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
	// PNBroadcastAvailabilityType message type for push notification broadcast
	// availability
	PNBroadcastAvailabilityType = "~#c90"
	// PNRegistrationType message type for sending a registration request to
	// a push notification server
	PNRegistrationType = "~#c91"
	// PNRegistrationConfirmationType message type to allow a push notification
	// server confirm a registration
	PNRegistrationConfirmationType = "~#c92"
)

// supportedMessage check if the message type is supported
func supportedMessage(msgType string) bool {
	_, ok := map[string]bool{
		NewContactKeyType:              true,
		ContactRequestType:             true,
		ConfirmedContactRequestType:    true,
		StandardMessageType:            true,
		SeenType:                       true,
		ContactUpdateType:              true,
		PNBroadcastAvailabilityType:    true,
		PNRegistrationType:             true,
		PNRegistrationConfirmationType: true,
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

func rawrChatMessage(raw string) string {
	bytes := []byte(raw)

	return fmt.Sprintf("0x%s", hex.EncodeToString(bytes))
}

func unrawrChatMessage(message string) ([]byte, error) {
	return hex.DecodeString(message[2:])
}

func messageFromEnvelope(u interface{}) (msg *Msg, err error) {
	payload := u.(map[string]interface{})["payload"]
	return messageFromPayload(payload.(string))
}

func messageFromPayload(payload string) (*Msg, error) {
	var msg []interface{}

	rawMsg, err := unrawrChatMessage(payload)
	if err != nil {
		return nil, err
	}
	spew.Dump(rawMsg)

	if err = json.Unmarshal(rawMsg, &msg); err != nil {
		return nil, err
	}

	if len(msg) < 1 {
		return nil, errors.New("unknown message format")
	}

	msgType := msg[0].(string)
	if !supportedMessage(msgType) {
		return nil, errors.New("unsupported message type")
	}

	properties := msg[1].([]interface{})

	timestamp := time.Now().Unix() * 100
	if len(properties) > 2 {
		timestamp = int64(properties[3].(float64))
	}

	return &Msg{
		Type:      msgType,
		From:      "TODO : someone",
		Text:      properties[0].(string),
		Timestamp: timestamp,
		Raw:       string(rawMsg),
	}, nil
}
