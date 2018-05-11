package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/status-im/status-go-sdk"
)

func main() {
	rpcClient, err := rpc.Dial("http://localhost:8545")
	checkErr(err)

	remoteClient := &remoteClient{rpcClient}
	client := sdk.New(remoteClient)

	a, err := client.SignupAndLogin("password")
	checkErr(err)

	ch, err := a.JoinPublicChannel("supu")
	if err != nil {
		checkErr(err)
	}

	fmt.Printf("%+v\n", ch)

	for range time.Tick(3 * time.Second) {
		message := fmt.Sprintf("PING : %d", time.Now().Unix())
		_ = ch.Publish(message)
	}
}

type remoteClient struct {
	c *rpc.Client
}

func (rc *remoteClient) Call(req *sdk.Request, res interface{}) error {
	return rc.c.Call(res, req.Method, req.Params)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
