package sdk

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
)

// SDK is a set of tools to interact with status node
type SDK struct {
	RPCClient  RPCClient
	address    string
	userName   string
	channels   []*Channel
	minimumPoW float64
}

// New creates a default SDK object
func New(address string) *SDK {
	return &SDK{
		RPCClient:  newRPC(address),
		minimumPoW: 0.001,
	}
}

// Close all channels you're subscribed to
func (c *SDK) Close() {
	for _, channel := range c.channels {
		channel.Close()
	}
}

// LoginResponse : json response returned by status_login.
type LoginResponse struct {
	Result struct {
		AddressKeyID string `json:"address_key_id"`
	} `json:"result"`
}

// Login to status with the given credentials
func (c *SDK) Login(addr, pwd string) error {
	var res LoginResponse
	cmd := fmt.Sprintf(statusLoginFormat, addr, pwd)
	if err := c.call(cmd, &res); err != nil {
		return err
	}
	c.address = res.Result.AddressKeyID

	return nil
}

// SignupResponse : json response returned by status_signup.
type SignupResponse struct {
	Result struct {
		Address  string `json:"address"`
		Pubkey   string `json:"pubkey"`
		Mnemonic string `json:"mnemonic"`
	} `json:"result"`
}

// Signup creates a new account with the given credentials
func (c *SDK) Signup(pwd string) (addr string, pubkey string, mnemonic string, err error) {
	var res SignupResponse

	cmd := fmt.Sprintf(statusSignupFormat, pwd)

	if err = c.call(cmd, &res); err != nil {
		return
	}

	return res.Result.Address, res.Result.Pubkey, res.Result.Mnemonic, err
}

// SignupAndLogin sign up and login on status network
func (c *SDK) SignupAndLogin(password string) (addr string, pubkey string, mnemonic string, err error) {
	addr, pubkey, mnemonic, err = c.Signup(password)
	if err != nil {
		return
	}
	err = c.Login(addr, password)
	return
}

/*
// Join a specific channel by name
func (c *SDK) Join(channelName string) (*Channel, error) {
	ch, err := c.joinPublicChannel(channelName)
	if err != nil {
		c.channels = append(c.channels, ch)
	}

	return ch, err
}
*/

// GenerateSymKeyFromPasswordResponse GenerateSymKeyFromPassword json response
type GenerateSymKeyFromPasswordResponse struct {
	Key string `json:"result"`
	ID  int    `json:"id"`
}

// Web3ShaResponse Web3Sha json response
type Web3ShaResponse struct {
	Result string `json:"result"`
}

// NewMessageFilterResponse NewMessageFilter json response
type NewMessageFilterResponse struct {
	Result string `json:"result"`
}

// JoinPublicChannel joins a status public channel
func (c *SDK) JoinPublicChannel(channelName string) (*Channel, error) {
	var symkeyResponse GenerateSymKeyFromPasswordResponse
	cmd := fmt.Sprintf(generateSymKeyFromPasswordFormat, channelName)
	if err := c.call(cmd, &symkeyResponse); err != nil {
		return nil, err
	}
	key := symkeyResponse.Key

	p := "0x" + hex.EncodeToString([]byte(channelName))

	var web3ShaResponse Web3ShaResponse
	cmd = fmt.Sprintf(web3ShaFormat, p, symkeyResponse.ID)
	if err := c.call(cmd, &web3ShaResponse); err != nil {
		return nil, err
	}
	topicID := web3ShaResponse.Result[0:10]

	var newMessageFilterResponse NewMessageFilterResponse
	cmd = fmt.Sprintf(newMessageFilterFormat, topicID, key)
	if err := c.call(cmd, &newMessageFilterResponse); err != nil {
		return nil, err
	}
	filterID := newMessageFilterResponse.Result

	ch := &Channel{
		conn:        c,
		channelName: channelName,
		filterID:    filterID,
		topicID:     topicID,
		channelKey:  key,
	}
	c.channels = append(c.channels, ch)

	return ch, nil
}

func (c *SDK) call(cmd string, res interface{}) error {
	log.Println("[ REQUST ] : " + cmd)
	body, err := c.RPCClient.Call(cmd)
	if err != nil {
		return err
	}
	log.Println("[ RESPONSE ] : " + body.(string))

	return json.Unmarshal([]byte(body.(string)), &res)
}
