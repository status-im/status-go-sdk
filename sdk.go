package sdk

import (
	"encoding/hex"
	"encoding/json"
	"log"
)

// SDK is a set of tools to interact with status node
type SDK struct {
	RPCClient  RPCClient
	address    string
	pubkey     string
	mnemonic   string
	userName   string
	channels   []*Channel
	minimumPoW float64
}

// New creates a default SDK object
func New(address string) *SDK {
	return &SDK{
		// RPCClient:  newRPC(address),
		minimumPoW: 0.001,
	}
}

// Close all channels you're subscribed to
func (c *SDK) Close() {
	for _, channel := range c.channels {
		channel.Close()
	}
}

// Login to status with the given credentials
func (c *SDK) Login(addr, pwd string) error {
	res, err := statusLoginRequest(c, addr, pwd)
	if err != nil {
		return err
	}
	c.address = res.Result.AddressKeyID

	return nil
}

// Signup creates a new account with the given credentials
func (c *SDK) Signup(pwd string) (addr string, pubkey string, mnemonic string, err error) {
	res, err := statusSignupRequest(c, pwd)
	if err != nil {
		return "", "", "", err
	}
	c.address = res.Result.Address
	c.pubkey = res.Result.Pubkey
	c.mnemonic = res.Result.Mnemonic

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

// NewMessageFilterResponse NewMessageFilter json response
type NewMessageFilterResponse struct {
	Result string `json:"result"`
}

// JoinPublicChannel joins a status public channel
func (c *SDK) JoinPublicChannel(channelName string) (*Channel, error) {
	symkeyResponse, err := shhGenerateSymKeyFromPasswordRequest(c, []string{channelName})
	if err != nil {
		return nil, err
	}
	key := symkeyResponse.Key

	topicID, err := c.calculatePublicChannelTopicID(channelName, symkeyResponse.ID)
	if err != nil {
		return nil, err
	}

	newMessageFilterResponse, err := newShhMessageFilterFormatRequest(c, []string{topicID}, key)
	if err != nil {
		return nil, err
	}

	filterID := newMessageFilterResponse.FilterID

	ch := &Channel{
		conn:       c,
		name:       channelName,
		filterID:   filterID,
		topicID:    topicID,
		channelKey: key,
	}
	c.channels = append(c.channels, ch)

	return ch, nil
}

func (c *SDK) calculatePublicChannelTopicID(name string, symkey int) (topicID string, err error) {
	p := "0x" + hex.EncodeToString([]byte(name))
	web3ShaResponse, err := web3Sha3Request(c, symkey, []string{p})
	if err != nil {
		return
	}
	topicID = web3ShaResponse.Result[0:10]

	return
}

func (c *SDK) call(cmd string, res interface{}) error {
	log.Println("[ REQUEST ] : " + cmd)
	body, err := c.RPCClient.Call(cmd)
	if err != nil {
		return err
	}
	log.Println("[ RESPONSE ] : " + body.(string))

	return json.Unmarshal([]byte(body.(string)), &res)
}
