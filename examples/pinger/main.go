package main

import (
	"fmt"
	"time"

	"github.com/status-im/status-go-sdk"
)

func main() {
	client := sdk.New("localhost:30303")

	addr, _, _, err := client.Signup("password")
	if err != nil {
		return
	}

	if err := client.Login(addr, "password"); err != nil {
		panic(err)
	}

	ch, err := client.JoinPublicChannel("supu")
	if err != nil {
		panic("Couldn't connect to status")
	}

	for range time.Tick(10 * time.Second) {
		message := fmt.Sprintf("PING : %d", time.Now().Unix())
		_ = ch.Publish(message)
	}
}
