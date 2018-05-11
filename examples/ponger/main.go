package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/status-im/status-go-sdk"
)

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

func main() {
	rpcClient, err := rpc.Dial("http://localhost:8545")
	checkErr(err)

	remoteClient := &remoteClient{rpcClient}
	client := sdk.New(remoteClient)

	a, err := client.SignupAndLogin("password")
	checkErr(err)

	ch, err := a.JoinPublicChannel("supu")
	checkErr(err)

	_, _ = ch.Subscribe(func(m *sdk.Msg) {
		log.Println("Message from ", m.From, " with body: ", m.Raw)

		if strings.Contains(m.Raw, "PING :") {
			time.Sleep(5 * time.Second)
			message := fmt.Sprintf("PONG : %d", time.Now().Unix())
			_ = ch.Publish(message)
		}
	})

	runtime.Goexit()
}
