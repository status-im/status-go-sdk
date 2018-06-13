package gfycat_test

import (
	"testing"

	"github.com/status-im/status-go-sdk/gfycat"
)

func TestHappyPath(t *testing.T) {
	testCases := []struct {
		address  string
		expected string
	}{
		/*
			{
				address:  "57348975ff9199ca636207a396b915d6b6a675b4",
				expected: "Winged Fitting Mosquito",
			},
		*/
		{
			address:  "0x04c940125c0b746c44dad3ef29a15a567fe63f291fa70bb06f0167c1711de2f13bd2dffb3932755775e06886ea3da91ccc5b9b44c8568c62d8e8b08a8bcdc168b3",
			expected: "Winding Murky Riograndeescuerzo",
		},
		{
			address:  "0x04f712bac774eb490bbe834556c18b64ee3878d3fd822cf5a584829aec69e4740132a2ae8e2dbe4d688d1a32069d542d95d203d28f1b6859e20a49610ddd2dea15",
			expected: "Adept Wasteful Northernfurseal",
		},
		{
			address:  "0x04bd61981bada15b60d272b488f3as123ec953afc496ddbec19f470aba35f35bff1234566bc394e8c6e1106b1284e7c2223122dbf3212fbcbbc9f1a41d06e029375",
			expected: "Quixotic Lined Moorhen",
		},
		{
			address:  "0x0111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			expected: "Boring Red Ivorybackedwoodswallow",
		},
		{
			address:  "",
			expected: "Round Infantile Killifish",
		},
	}

	for _, tc := range testCases {
		username := gfycat.AddressBasedUsername(tc.address)
		if tc.expected != username {
			t.Fatalf("Expecting " + tc.expected + " got " + username)
		}
	}
}
