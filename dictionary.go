package sdk

func shhGenerateSymKeyFromPasswordRequest(sdk *SDK, password string) (string, error) {
	// `{"jsonrpc":"2.0","id":2950,"method":"shh_generateSymKeyFromPassword","params":["%s"]}`
	var resp string
	return resp, sdk.RPCClient.Call(&resp, "shh_generateSymKeyFromPassword", password)
}

type shhFilterFormatParam struct {
	AllowP2P     bool     `json:"allowP2P"`
	Topics       []string `json:"topics"`
	Type         string   `json:"type"`
	SymKeyID     string   `json:"symKeyID"`
	PrivateKeyID string   `json:"privateKeyID"`
}

func newShhMessageFilterFormatRequest(sdk *SDK, topics []string, symKey, privateKeyID string) (string, error) {
	// `{"jsonrpc":"2.0","id":2,"method":"shh_newMessageFilter","params":[{"allowP2P":true,"topics":["%s"],"type":"sym","symKeyID":"%s"}]}`
	var res string
	params := &shhFilterFormatParam{
		AllowP2P:     true,
		Topics:       topics,
		Type:         "sym",
		SymKeyID:     symKey,
		PrivateKeyID: privateKeyID,
	}
	if len(symKey) > 0 {
		params.SymKeyID = symKey
	}
	if len(privateKeyID) > 0 {
		params.PrivateKeyID = privateKeyID
	}

	return res, sdk.RPCClient.Call(&res, "shh_newMessageFilter", params)
}

func web3Sha3Request(sdk *SDK, data string) (string, error) {
	// `{"jsonrpc":"2.0","method":"web3_sha3","params":["%s"],"id":%d}`
	var res string
	return res, sdk.RPCClient.Call(&res, "web3_sha3", data)
}

type statusLoginParam struct {
	Address  string `json:"address"`
	Password string `json:"password"`
}

type loginResponse struct {
	AddressKeyID string `json:"address_key_id"`
}

func statusLoginRequest(sdk *SDK, address, password string) (*loginResponse, error) {
	// `{"jsonrpc":"2.0","method":"status_login","params":[{"address":"%s","password":"%s"}]}`
	var res loginResponse

	params := &statusLoginParam{
		Address:  address,
		Password: password,
	}
	return &res, sdk.RPCClient.Call(&res, "status_login", params)
}

type statusSignupParam struct {
	Password string `json:"password"`
}

type signupResponse struct {
	Address  string `json:"address"`
	Pubkey   string `json:"pubkey"`
	Mnemonic string `json:"mnemonic"`
}

func statusSignupRequest(sdk *SDK, password string) (*signupResponse, error) {
	// `{"jsonrpc":"2.0","method":"status_signup","params":[{"password":"%s"}]}`
	var res signupResponse

	params := &statusSignupParam{
		Password: password,
	}

	return &res, sdk.RPCClient.Call(&res, "status_signup", params)
}

func shhGetFilterMessagesRequest(sdk *SDK, filter string) (interface{}, error) {
	// `{"jsonrpc":"2.0","id":2968,"method":"shh_getFilterMessages","params":["%s"]}`
	var res interface{}

	return res, sdk.RPCClient.Call(&res, "shh_getFilterMessages", filter)
}

// Message message to be sent though ssh_post calls
type Message struct {
	Signature string  `json:"sig"`
	SymKeyID  string  `json:"symKeyID,omitempty"`
	PubKey    string  `json:"pubKey,omitempty"`
	Payload   string  `json:"payload"`
	Topic     string  `json:"topic"`
	TTL       uint32  `json:"ttl"`
	PowTarget float64 `json:"powTarget"`
	PowTime   uint32  `json:"powTime"`
}

func shhPostRequest(sdk *SDK, msg *Message) (string, error) {
	var res string
	return res, sdk.RPCClient.Call(&res, "shhext_post", msg)
}

func addSymKey(sdk *SDK, symKey string) (string, error) {
	// {"jsonrpc":"2.0","method":"shh_addSymKey", "params":["` + symkey + `"], "id":1}
	var res string
	return res, sdk.RPCClient.Call(&res, "shh_addSymKey", symKey)
}
