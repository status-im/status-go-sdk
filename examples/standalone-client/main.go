package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/whisper/shhclient"
	"github.com/ethereum/go-ethereum/whisper/whisperv5"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func makePayload(text string) string {
	timestamp := time.Now().Unix() * 1000
	format := `["~#c4",["%s","text/plain","~:public-group-user-message",%d,%d]]`
	payload := fmt.Sprintf(format, text, timestamp*100, timestamp)

	return payload
}

func topicFromChatName(chatName string) whisperv5.TopicType {
	h := sha3.NewKeccak256()
	h.Write([]byte(chatName))
	fullTopic := h.Sum(nil)
	return whisperv5.BytesToTopic(fullTopic)
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">> ")
	text, err := reader.ReadString('\n')
	check(err)

	return strings.TrimSpace(text)
}

func main() {
	ipcPath := os.Getenv("IPC_PATH")
	if ipcPath == "" {
		fmt.Printf("you must specify the IPC_PATH environment variable. ")
		fmt.Printf("Try with: export IPC_PATH=path-to-file\n")
		os.Exit(1)
	}

	whisperClient, err := shhclient.Dial(ipcPath)
	check(err)

	ctx := context.Background()

	// public chat
	chatName := "testsdkfoobarbaz"

	// topic
	topic := topicFromChatName(chatName)

	// symkey
	symKeyID, err := whisperClient.GenerateSymmetricKeyFromPassword(ctx, chatName)
	check(err)

	// keypair
	keyPairID, err := whisperClient.NewKeyPair(ctx)
	check(err)

	for {
		text := readLine()

		// payload
		payload := makePayload(text)

		msg := whisperv5.NewMessage{
			SymKeyID:  symKeyID,
			Sig:       keyPairID,
			TTL:       10,
			Topic:     topic,
			Payload:   []byte(payload),
			PowTime:   1,
			PowTarget: 0.001,
		}

		hash, err := whisperClient.Post(ctx, msg)
		check(err)
		fmt.Printf("-- sent message with has %s\n", hash)
	}
}
