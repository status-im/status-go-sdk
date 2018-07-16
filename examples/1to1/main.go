package main

import (
	"log"
	"runtime"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/status-im/status-go-sdk"
	"github.com/status-im/status-go-sdk/examples/lib"
)

const (
	botUsername = "1to1bot"
)

func main() {
	rpcClient, err := rpc.Dial("http://localhost:8545")
	checkErr(err)

	client := sdk.New(rpcClient)
	bot, err := client.SignupAndLogin("password")
	bot.Username = botUsername
	bot.Image = lib.BotImage

	// Wait for contact requests
	err = bot.OnContactRequest(func(ct *sdk.Contact) {
		log.Println("Received contact request from " + ct.Name)
		log.Println("Accepting " + ct.Name + "'s contact request'")
		checkErr(ct.Accept())

		_, err = ct.Subscribe(func(m *sdk.Msg) {
			if m.Type != sdk.StandardMessageType {
				return
			}

			properties := m.Properties.(*sdk.PublishMsg)
			log.Println(properties.Text)

			// Print the response.
			err = ct.Publish(lib.AskMitsuku(properties.Text))
			checkErr(err)
		})
		checkErr(err)
	})
	checkErr(err)

	log.Println("Making the bot visible on a public channel so you can easily get its address")
	ch, err := bot.JoinPublicChannel("supu")
	checkErr(err)
	err = ch.Publish("hey there! I'm " + botUsername + ", and I am accepting contact requests :P")
	checkErr(err)

	runtime.Goexit()
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
