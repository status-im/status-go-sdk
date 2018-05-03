package sdk

import (
	"github.com/ethereum/go-ethereum/whisper/whisperv6"
)

type clientWrapper struct {
	rpcClient RPCClient
}

func newClientWrapper(c RPCClient) *clientWrapper {
	return &clientWrapper{c}
}

func (c *clientWrapper) StatusSignup(password string) (*SignupResponse, error) {
	var resp SignupResponse
	return &resp, c.rpcClient.Call(&resp, "status_signup", map[string]string{"password": password})
}

func (c *clientWrapper) StatusLogin(address, password string) (*LoginResponse, error) {
	var resp LoginResponse
	return &resp, c.rpcClient.Call(&resp, "status_login", map[string]string{"address": address, "password": password})
}

func (c *clientWrapper) ShhGenerateSymKeyFromPassword(password string) (string, error) {
	var symKeyID string
	return symKeyID, c.rpcClient.Call(&symKeyID, "shh_generateSymKeyFromPassword", password)
}

func (c *clientWrapper) ShhPost(msg *whisperv6.NewMessage) (string, error) {
	var hash string
	return hash, c.rpcClient.Call(&hash, "shh_post", msg)
}

func (c *clientWrapper) ShhNewMessageFilter(criteria whisperv6.Criteria) (string, error) {
	var id string
	return id, c.rpcClient.Call(&id, "shh_newMessageFilter", criteria)
}

func (c *clientWrapper) ShhFilterMessages(id string) ([]*whisperv6.Message, error) {
	var msgs []*whisperv6.Message
	return msgs, c.rpcClient.Call(&msgs, "shh_getFilterMessages", id)
}
