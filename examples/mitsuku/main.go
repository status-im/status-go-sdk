// This package provides an easy example of a conversational chatbot integrated
// with rest version of mitsuku (https://www.pandorabots.com/mitsuku/), a
// three-time winner of the Loebner Prize Turing Test.
package main

import (
	"log"
	"runtime"

	"github.com/ethereum/go-ethereum/rpc"
	sdk "github.com/status-im/status-go-sdk"
	"github.com/status-im/status-go-sdk/examples/lib"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	pwd := "password"

	rpcClient, err := rpc.Dial("http://localhost:8545")
	checkErr(err)

	client := sdk.New(rpcClient)

	a, err := client.SignupAndLogin(pwd)
	checkErr(err)

	ch, err := a.JoinPublicChannel("mitsuku")
	checkErr(err)

	_, _ = ch.Subscribe(func(m *sdk.Msg) {
		if m.Type != sdk.StandardMessageType {
			return
		}
		if a.PubKey == m.PubKey {
			return
		}
		properties := m.Properties.(*sdk.PublishMsg)

		// Print the response.
		err := ch.Publish(lib.AskMitsuku(properties.Text))
		checkErr(err)

	})
	runtime.Goexit()
}
