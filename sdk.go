package sdk

import (
	"encoding/hex"
	"fmt"
)

// SDK is a set of tools to interact with status node
type SDK struct {
	RPCClient  RPCClient
	address    string
	userName   string
	channels   []*Channel
	minimumPoW string
}

// New creates a default SDK object
func New(address string) *SDK {
	return &SDK{
		RPCClient:  newRPC(address),
		minimumPoW: "0.01",
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
	cmd := fmt.Sprintf(statusLoginFormat, addr, pwd)
	res, err := c.call(cmd)
	if err != nil {
		return err
	}
	// TODO(adriacidre) unmarshall and treat the response
	println(res)

	return nil
}

// Signup creates a new account with the given credentials
func (c *SDK) Signup(pwd string) error {
	cmd := fmt.Sprintf(statusSignupFormat, pwd)
	res, err := c.call(cmd)
	if err != nil {
		return err
	}
	// TODO(adriacidre) unmarshall and treat the response
	println(res)

	return nil
}

// SignupOrLogin will attempt to login with given credentials, in first instance
// or will sign up in case login does not work
func (c *SDK) SignupOrLogin(user, password string) error {
	if err := c.Login(user, password); err != nil {
		// TODO (adriacidre) handle this error
		_ = c.Signup(password)
		return c.Login(user, password)
	}

	return nil
}

// Join a specific channel by name
func (c *SDK) Join(channelName string) (*Channel, error) {
	ch, err := c.joinPublicChannel(channelName)
	if err != nil {
		c.channels = append(c.channels, ch)
	}

	return ch, err
}

func (c *SDK) joinPublicChannel(channelName string) (*Channel, error) {
	cmd := fmt.Sprintf(generateSymKeyFromPasswordFormat, channelName)
	res, _ := c.call(cmd)
	f := unmarshalJSON(res.(string))

	key := f.(map[string]interface{})["result"].(string)
	id := int(f.(map[string]interface{})["id"].(float64))

	src := []byte(channelName)
	p := "0x" + hex.EncodeToString(src)

	cmd = fmt.Sprintf(web3ShaFormat, p, id)
	res, _ = c.call(cmd)
	topic := res.(map[string]interface{})["result"].(string)
	topic = topic[0:10]

	cmd = fmt.Sprintf(newMessageFilterFormat, topic, key)
	res, _ = c.call(cmd)
	f3 := unmarshalJSON(res.(string))
	filterID := f3.(map[string]interface{})["result"].(string)

	return &Channel{
		conn:        c,
		channelName: channelName,
		filterID:    filterID,
		topic:       topic,
		channelKey:  key,
	}, nil
}

func (c *SDK) call(cmd string) (interface{}, error) {
	return c.RPCClient.Call(cmd)
}
