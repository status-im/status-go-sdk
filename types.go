package sdk

type SignupResponse struct {
	Address  string `json:"address"`
	PubKey   string `json:"pubkey"`
	Mnemonic string `json:"mnemonic"`
}

type LoginResponse struct {
	AddressKeyID string `json:"address_key_id"`
}
