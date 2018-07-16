package sdk

import (
	"encoding/hex"
)

// Account represents a logged in user on statusd node
type Account struct {
	conn         *SDK
	Address      string
	AddressKeyID string
	PubKey       string
	Mnemonic     string
	Username     string
	Image        string
	channels     []*Channel
}

// JoinPublicChannel joins a status public channel
func (a *Account) JoinPublicChannel(name string) (*Channel, error) {
	return a.createAndJoin(name, name)
}

// CreatePrivateChannel creates and joins a private channel
func (a *Account) CreatePrivateChannel(name, password string) (*Channel, error) {
	return a.createAndJoin(name, password)
}

func (a *Account) createAndJoin(name, password string) (*Channel, error) {
	symKey, err := shhGenerateSymKeyFromPasswordRequest(a.conn, password)
	if err != nil {
		return nil, err
	}

	topicID, err := a.calculatePublicChannelTopicID(name)
	if err != nil {
		return nil, err
	}

	return a.Join(name, topicID, symKey)
}

// OnContactRequest executes specified ContactRequestHandler logic when a contact
// request is received
func (a *Account) OnContactRequest(fn ContactRequestHandler) error {
	topicID, err := a.calculatePublicChannelTopicID(discoveryChannelName)
	if err != nil {
		return err
	}

	filterID, err := newShhMessageFilterFormatRequest(a.conn, []string{topicID}, "", a.PubKey)
	if err != nil {
		println(err.Error())
		return err
	}

	discoveryChannel := &Channel{
		account:       a,
		name:          discoveryChannelName,
		filterID:      filterID,
		TopicID:       topicID,
		ChannelPubKey: a.PubKey,
	}

	var sub *Subscription
	sub, err = discoveryChannel.Subscribe(func(msg *Msg) {
		switch c := msg.Properties.(type) {
		case *NewContactKeyMsg:
			fn(&Contact{
				account: a,
				SymKey:  c.Address,
				TopicID: c.Topic,
				Name:    c.Contact.Name,
				Address: c.Contact.Address,
				Image:   c.Contact.Image,
				ch:      discoveryChannel,
				sub:     sub,
				PubKey:  msg.PubKey,
			})
		}
	})

	return err
}

// Join joins a status channel
func (a *Account) Join(channelName, topicID, symKey string) (*Channel, error) {
	filterID, err := newShhMessageFilterFormatRequest(a.conn, []string{topicID}, symKey, "")
	if err != nil {
		return nil, err
	}

	ch := &Channel{
		account:       a,
		name:          channelName,
		filterID:      filterID,
		TopicID:       topicID,
		ChannelSymKey: symKey,
	}
	a.channels = append(a.channels, ch)

	return ch, nil
}

func (a *Account) calculatePublicChannelTopicID(name string) (topicID string, err error) {
	p := "0x" + hex.EncodeToString([]byte(name))
	hash, err := web3Sha3Request(a.conn, p)
	if err != nil {
		return
	}
	topicID = hash[0:10]

	return
}

// Close all channels you're subscribed to
func (a *Account) Close() {
	for _, ch := range a.channels {
		ch.Close()
	}
}
