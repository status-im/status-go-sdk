package sdk

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDictionaryMethods(t *testing.T) {
	var str string
	var genericRes interface{}

	testCases := []struct {
		Description string
		Response    interface{}
		Method      string
		Params      interface{}
		Callback    func(*SDK)
	}{
		{
			Description: "statusLoginRequest",
			Response:    &loginResponse{},
			Method:      "status_login",
			Params: &statusLoginParam{
				Address:  "ADDRESS",
				Password: "PASSWORD",
			},
			Callback: func(sdk *SDK) {
				response, err := statusLoginRequest(sdk, "ADDRESS", "PASSWORD")
				assert.NoError(t, err)
				assert.NotNil(t, response)
			},
		},
		{
			Description: "statusSignupRequest",
			Response:    &signupResponse{},
			Method:      "status_signup",
			Params: &statusSignupParam{
				Password: "PASSWORD",
			},
			Callback: func(sdk *SDK) {
				response, err := statusSignupRequest(sdk, "PASSWORD")
				assert.NoError(t, err)
				assert.NotNil(t, response)
			},
		},
		{
			Description: "shhGenerateSymKeyFromPasswordRequest",
			Response:    &str,
			Method:      "shh_generateSymKeyFromPassword",
			Params:      "PASSWORD",
			Callback: func(sdk *SDK) {
				response, err := shhGenerateSymKeyFromPasswordRequest(sdk, "PASSWORD")
				assert.NoError(t, err)
				assert.NotNil(t, response)
			},
		},
		{
			Description: "shhPostRequest",
			Response:    &str,
			Method:      "shhext_post",
			Params:      &Message{},
			Callback: func(sdk *SDK) {
				response, err := shhPostRequest(sdk, &Message{})
				assert.NoError(t, err)
				assert.NotNil(t, response)
			},
		},
		{
			Description: "shhGetFilterMessagesRequest",
			Response:    &genericRes,
			Method:      "shh_getFilterMessages",
			Params:      "",
			Callback: func(sdk *SDK) {
				response, err := shhGetFilterMessagesRequest(sdk, "")
				assert.NoError(t, err)
				assert.Nil(t, response)
			},
		},
		{
			Description: "web3Sha3Request",
			Response:    &str,
			Method:      "web3_sha3",
			Params:      "",
			Callback: func(sdk *SDK) {
				response, err := web3Sha3Request(sdk, "")
				assert.NoError(t, err)
				assert.NotNil(t, response)
			},
		},
		{
			Description: "newShhMessageFilterFormatRequest",
			Response:    &str,
			Method:      "shh_newMessageFilter",
			Params: &shhFilterFormatParam{
				AllowP2P: true,
				Topics:   []string{},
				Type:     "sym",
				SymKeyID: "",
			},
			Callback: func(sdk *SDK) {
				response, err := newShhMessageFilterFormatRequest(sdk, []string{}, "", "")
				assert.NoError(t, err)
				assert.NotNil(t, response)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			rpcClient := NewMockRPCClient(mockCtrl)
			sdk := New(rpcClient)
			rpcClient.EXPECT().Call(tc.Response, tc.Method, tc.Params)
			tc.Callback(sdk)
		})
	}
}
