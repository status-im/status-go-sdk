package gfycat

import (
	"encoding/hex"
	"math/big"
	"math/rand"
	"strings"

	"github.com/seehuhn/mt19937"
)

// AddressBasedUsername gets a "random" username based on an address.
func AddressBasedUsername(address string) string {
	pubKey := fromHex(address)

	i := new(big.Int)
	i.SetBytes(pubKey)

	rng := rand.New(mt19937.New())
	seed := i.Int64()

	rng.Seed(seed)

	adj1 := adjectives[rng.Intn(len(adjectives))]
	adj2 := adjectives[rng.Intn(len(adjectives))]
	animal := animals[rng.Intn(len(animals))]

	return strings.Title(adj1 + " " + adj2 + " " + animal)
}

func fromHex(s string) []byte {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			println("-------------------")
			s = s[2:]
		}
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return hex2Bytes(s)
}
func hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}
