package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	sdk "github.com/status-im/status-go-sdk"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	rpcClient, err := rpc.Dial("http://localhost:8545")
	checkErr(err)
	client := sdk.NewClient(rpcClient)

	address, err := client.Signup("foobar")
	checkErr(err)
	fmt.Printf("Account created: %+v\n", address)

	_, err = client.Login(address, "foobar")
	checkErr(err)

	chatName := "testsdkfoobarbaz"
	symKeyID, err := client.PublicChatSymKey(chatName)
	checkErr(err)

	topic := client.PublicChatTopic(chatName)

	for {
		text := fmt.Sprintf("PING %d", time.Now().UnixNano())
		hash, err := client.Post(symKeyID, topic, text)
		checkErr(err)
		fmt.Printf("Message sent with hash %s\n", hash)
		time.Sleep(time.Second * 5)
	}
}
