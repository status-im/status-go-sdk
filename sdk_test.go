package sdk

import (
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

func TestNewSDK(t *testing.T) {
	rpcClient, err := rpc.Dial("http://localhost:8545")
	assert.NoError(t, err)

	client := New(rpcClient)
	assert.Equal(t, 0.001, client.minimumPoW)
	assert.Equal(t, rpcClient, client.RPCClient)
}
