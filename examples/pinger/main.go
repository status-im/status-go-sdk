package main

import (
	"fmt"
	"time"

	"github.com/status-im/status-go-sdk"
)

func main() {
	sdk := sdk.New("localhost:30303")
	if err := sdk.Signup("111222333"); err != nil {
		panic("Couldn't create an account")
	}

	ch, err := sdk.Join("supu")
	if err != nil {
		panic("Couldn't connect to status")
	}

	for range time.Tick(10 * time.Second) {
		message := fmt.Sprintf("PING : %d", time.Now().Unix())
		_ = ch.Publish(message)
	}
}
