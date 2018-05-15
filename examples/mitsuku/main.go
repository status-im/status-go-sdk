// This package provides an easy example of a conversational chatbot integrated
// with rest version of mitsuku (https://www.pandorabots.com/mitsuku/), a
// three-time winner of the Loebner Prize Turing Test.
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/rpc"
	sdk "github.com/status-im/status-go-sdk"
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
	pwd := "password"

	rpcClient, err := rpc.Dial("http://localhost:8545")
	checkErr(err)

	remoteClient := &remoteClient{rpcClient}
	client := sdk.New(remoteClient)

	a, err := client.SignupAndLogin(pwd)

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
		ch.Publish(askMitsuku(properties.Text))

	})
	runtime.Goexit()
}

func askMitsuku(question string) string {
	hc := http.Client{}
	uri := "https://kakko.pandorabots.com/pandora/talk?botid=87437a824e345a0d&skin=chat"
	form := url.Values{}
	form.Add("message", question)
	req, _ := http.NewRequest("POST", uri, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := hc.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	parts := strings.Split(string(body), "</B>")
	parts = strings.Split(parts[2], "<br>")

	return strings.Trim(parts[0], " ")
}
