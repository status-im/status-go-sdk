package sdk

//RPCClient is a client to manage all rpc calls
type RPCClient interface {
	Call(result interface{}, method string, args ...interface{}) error
}
