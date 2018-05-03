package sdk

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/whisper/whisperv6"
)

type Client struct {
	cw         *clientWrapper
	keyID      string
	address    string
	minimumPoW float64
}

func NewClient(rpcclient RPCClient) *Client {
	cw := newClientWrapper(rpcclient)
	return &Client{cw: cw}
}

func (c *Client) Signup(password string) (string, error) {
	resp, err := c.cw.StatusSignup(password)
	if err != nil {
		return "", err
	}

	return resp.Address, nil
}

func (c *Client) Login(address, password string) (string, error) {
	resp, err := c.cw.StatusLogin(address, password)
	if err != nil {
		return "", err
	}

	c.keyID = resp.AddressKeyID

	return resp.AddressKeyID, nil
}

func (c *Client) PublicChatSymKey(chatName string) (string, error) {
	return c.cw.ShhGenerateSymKeyFromPassword(chatName)
}

func (c *Client) PublicChatTopic(chatName string) whisperv6.TopicType {
	h := sha3.NewKeccak256()
	h.Write([]byte(chatName))
	fullTopic := h.Sum(nil)

	return whisperv6.BytesToTopic(fullTopic)
}

func (c *Client) MakePayload(text string) string {
	timestamp := time.Now().Unix() * 1000
	format := `["~#c4",["%s","text/plain","~:public-group-user-message",%d,%d]]`
	payload := fmt.Sprintf(format, text, timestamp*100, timestamp)

	return payload
}

func (c *Client) Post(symKeyID string, topic whisperv6.TopicType, text string) (string, error) {
	payload := c.MakePayload(text)
	msg := &whisperv6.NewMessage{
		SymKeyID:  symKeyID,
		Sig:       c.keyID,
		TTL:       10,
		Topic:     topic,
		Payload:   []byte(payload),
		PowTime:   1,
		PowTarget: 0.001,
	}

	return c.cw.ShhPost(msg)
}

func (c *Client) NewFilter(topic whisperv6.TopicType, symKeyID string) (string, error) {
	cr := whisperv6.Criteria{
		AllowP2P: true,
		Topics:   []whisperv6.TopicType{topic},
		SymKeyID: symKeyID,
	}

	return c.cw.ShhNewMessageFilter(cr)
}

func (c *Client) FilterMessages(filterID string) ([]*whisperv6.Message, error) {
	return c.cw.ShhFilterMessages(filterID)
}
