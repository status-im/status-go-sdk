package sdk

type shhRequest struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type generateSymKeyFromPasswordResponse struct {
	Key string `json:"result"`
	ID  int    `json:"id"`
}

func shhGenerateSymKeyFromPasswordRequest(sdk *SDK, password string) (string, error) {
	// `{"jsonrpc":"2.0","id":2950,"method":"shh_generateSymKeyFromPassword","params":["%s"]}`
	var resp string
	return resp, sdk.RPCClient.Call(&resp, "shh_generateSymKeyFromPassword", password)
}

type shhFilterFormatParam struct {
	AllowP2P bool     `json:"allowP2P"`
	Topics   []string `json:"topics"`
	Type     string   `json:"type"`
	SymKeyID string   `json:"symKeyID"`
}

type newMessageFilterResponse struct {
	FilterID string `json:"result"`
}

func newShhMessageFilterFormatRequest(sdk *SDK, topics []string, symKey string) (string, error) {
	// `{"jsonrpc":"2.0","id":2,"method":"shh_newMessageFilter","params":[{"allowP2P":true,"topics":["%s"],"type":"sym","symKeyID":"%s"}]}`
	var res string
	params := &shhFilterFormatParam{
		AllowP2P: true,
		Topics:   topics,
		Type:     "sym",
		SymKeyID: symKey,
	}

	return res, sdk.RPCClient.Call(&res, "shh_newMessageFilter", params)
}

type web3ShaResponse struct {
	Result string `json:"result"`
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

type getFilterMessagesResponse struct {
	Result interface{} `json:"result"`
}

func shhGetFilterMessagesRequest(sdk *SDK, filter string) (interface{}, error) {
	// `{"jsonrpc":"2.0","id":2968,"method":"shh_getFilterMessages","params":["%s"]}`
	var res interface{}

	return res, sdk.RPCClient.Call(&res, "shh_getFilterMessages", filter)
}

type Message struct {
	Signature string  `json:"sig"`
	SymKeyID  string  `json:"symKeyID"`
	Payload   string  `json:"payload"`
	Topic     string  `json:"topic"`
	TTL       uint32  `json:"ttl"`
	PowTarget float64 `json:"powTarget"`
	PowTime   uint32  `json:"powTime"`
}

// error response {"jsonrpc":"2.0","id":633,"error":{"code":-32000,"message":"message rejected, PoW too low"}}
type sshPostError struct {
	Code    float64 `json:"code"`
	Message string  `json:"message"`
}

type shhPostResponse struct {
	Error *sshPostError `json:"error"`
}

func shhPostRequest(sdk *SDK, msg *Message) (string, error) {
	var res string
	return res, sdk.RPCClient.Call(&res, "shh_post", msg)
}
