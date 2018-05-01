package sdk

import "github.com/valyala/gorpc"

// RPCClient is a client to manage all rpc calls
type RPCClient interface {
	Call(request interface{}) (response interface{}, err error)
}

func newRPC(address string) RPCClient {
	rpc := &gorpc.Client{
		Addr: address, // "rpc.server.addr:12345",
	}
	rpc.Start()

	return rpc
}
