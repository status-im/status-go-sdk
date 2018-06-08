package sdk

// RPCClient is a client to manage all rpc calls
type RPCClient interface {
	Call(result interface{}, method string, args ...interface{}) error
}

// SDK is a set of tools to interact with status node
type SDK struct {
	RPCClient  RPCClient
	minimumPoW float64
}

// New creates a default SDK object
func New(c RPCClient) *SDK {
	return &SDK{
		RPCClient:  c,
		minimumPoW: 0.001,
	}
}

// Login to status with the given credentials
func (c *SDK) Login(addr, pwd string) (a *Account, err error) {
	res, err := statusLoginRequest(c, addr, pwd)
	if err != nil {
		return a, err
	}
	return &Account{
		conn:         c,
		AddressKeyID: res.AddressKeyID,
	}, err
}

// Signup creates a new account with the given credentials
func (c *SDK) Signup(pwd string) (a *Account, err error) {
	res, err := statusSignupRequest(c, pwd)

	if err != nil {
		return a, err
	}
	return &Account{
		conn:     c,
		Address:  res.Address,
		PubKey:   res.Pubkey,
		Mnemonic: res.Mnemonic,
	}, err

}

// SignupAndLogin sign up and login on status network
func (c *SDK) SignupAndLogin(password string) (a *Account, err error) {
	a, err = c.Signup(password)
	if err != nil {
		return
	}
	la, err := c.Login(a.Address, password)
	a.AddressKeyID = la.AddressKeyID
	return
}

// Request json request.
type Request struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// NewMessageFilterResponse NewMessageFilter json response
type NewMessageFilterResponse struct {
	Result string `json:"result"`
}
