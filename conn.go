package sdk

import (
	"encoding/hex"
	"fmt"

	"github.com/valyala/gorpc"
)

// Conn : TODO ...
type Conn struct {
	rpc        *gorpc.Client
	address    string
	userName   string
	channels   []*Channel
	minimumPoW string
}

func NewConn(address string) *Conn {
	rpc := &gorpc.Client{
		Addr: address, // "rpc.server.addr:12345",
	}
	rpc.Start()

	return &Conn{
		rpc:        rpc,
		minimumPoW: "0.01",
	}
}

func (c *Conn) Close() {
	for _, channel := range c.channels {
		channel.Close()
	}
}

// Login logs in to the network with the given credentials
func (c *Conn) Login(addr, pwd string) error {
	cmd := fmt.Sprintf(statusLoginFormat, addr, pwd)
	res, err := c.rpc.Call(cmd)
	if err != nil {
		return err
	}
	// TODO(adriacidre) unmarshall and treat the response
	println(res)

	return nil
}

// Signup creates a new account with the given credentials
func (c *Conn) Signup(pwd string) error {
	cmd := fmt.Sprintf(statusSignupFormat, pwd)
	res, err := c.rpc.Call(cmd)
	println("------")
	println(res)
	println(err.Error())
	println("------")
	if err != nil {
		return err
	}
	// TODO(adriacidre) unmarshall and treat the response
	println(res)

	return nil
}

// SignupOrLogin will attempt to login with given credentials, in first instance
// or will sign up in case login does not work
func (c *Conn) SignupOrLogin(user, password string) error {
	if err := c.Login(user, password); err != nil {
		c.Signup(password)
		return c.Login(user, password)
	}

	return nil
}

// Join a specific channel by name
func (c *Conn) Join(channelName string) (*Channel, error) {
	ch, err := c.joinPublicChannel(channelName)
	if err != nil {
		c.channels = append(c.channels, ch)
	}

	return ch, err
}

func (c *Conn) joinPublicChannel(channelName string) (*Channel, error) {
	cmd := fmt.Sprintf(generateSymKeyFromPasswordFormat, channelName)
	res, _ := c.rpc.Call(cmd)
	f := unmarshalJSON(res.(string))

	key := f.(map[string]interface{})["result"].(string)
	id := int(f.(map[string]interface{})["id"].(float64))

	src := []byte(channelName)
	p := "0x" + hex.EncodeToString(src)

	cmd = fmt.Sprintf(web3ShaFormat, p, id)
	res, _ = c.rpc.Call(cmd)
	topic := res.(map[string]interface{})["result"].(string)
	topic = topic[0:10]

	cmd = fmt.Sprintf(newMessageFilterFormat, topic, key)
	res, _ = c.rpc.Call(cmd)
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
