package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/status-im/status-go-sdk"
)

func main() {
	client := sdk.New("localhost:30303")

	if err := client.SignupOrLogin("supu", "password"); err != nil {
		panic(err)
	}

	ch, err := client.Join("supu")
	if err != nil {
		panic("Couldn't connect to status")
	}

	_, _ = ch.Subscribe(func(m *sdk.Msg) {
		log.Println("Message from ", m.From, " with body: ", m.Text)

		if strings.Contains(m.Text, "PING :") {
			time.Sleep(5 * time.Second)
			message := fmt.Sprintf("PONG : %d", time.Now().Unix())
			_ = ch.Publish(message)
		}

	})

	runtime.Goexit()
}
