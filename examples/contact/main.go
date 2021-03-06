package main

import (
	"log"
	"runtime"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/status-im/status-go-sdk"
	"github.com/status-im/status-go-sdk/examples/lib"
)

const (
	botUsername = "TestBot"
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
		// And we simply accept them all (you can add any conditional you want here)
		err = ct.Accept()
		checkErr(err)
	})
	checkErr(err)

	log.Println("Makeing the bot visible on a public channel so you can easily get its address")
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
