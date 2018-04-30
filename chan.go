package sdk

import (
	"encoding/json"
	"fmt"
	"log"
)

// Channel : ...
type Channel struct {
	conn          *SDK
	channelName   string
	filterID      string
	channelKey    string
	topic         string
	subscriptions []*Subscription
}

// Publish : Publishes a message with the given body on the current channel
func (c *Channel) Publish(body string) error {
	msg := NewMsg(c.conn.userName, body, c.channelName)
	cmd := fmt.Sprintf(standardMessageFormat,
		c.conn.address,
		c.channelKey,
		msg.ToPayload(),
		c.topic,
		c.conn.minimumPoW,
	)
	_, err := c.conn.call(cmd)

	return err
}

// Subscribe : ...
func (c *Channel) Subscribe(fn MsgHandler) (*Subscription, error) {
	log.Println("Subscribed to channel '", c.channelName, "'")
	subscription := &Subscription{}
	go subscription.Subscribe(c, fn)
	c.subscriptions = append(c.subscriptions, subscription)

	return subscription, nil
}

// Close current channel and all its subscriptions
func (c *Channel) Close() {
	for _, sub := range c.subscriptions {
		c.removeSubscription(sub)
	}
}

// NewContactKeyRequest : First message that is sent to a future contact. At that
// point the only topic we know that the contact is filtering is the
// discovery-topic with his public key so that is what NewContactKey will
// be sent to.
// It contains the sym-key and topic that will be used for future communications
// as well as the actual message that we want to send.
// The sym-key and topic are generated randomly because we don’t want to have
// any correlation between a topic and its participants to avoid leaking
// metadata.
// When one of the contacts recovers his account, a NewContactKey message is
// sent as well to change the symmetric key and topic.
func (c *Channel) NewContactKeyRequest(username string) {
	contactRequest := fmt.Sprintf(contactRequestMsg, username, "", "", "")
	msg := fmt.Sprintf(newContactKeyMsg, c.conn.address, c.topic, contactRequest)

	c.callStandardMsg(msg)
}

// ContactRequest : Wrapped in a NewContactKey message when initiating a contact request.
func (c *Channel) ContactRequest(username, image string) {
	msg := fmt.Sprintf(contactRequestMsg, username, image, c.conn.address, "")
	c.callStandardMsg(msg)
}

// ConfirmedContactRequest : This is the message that will be sent when the
// contact accepts the contact request. It will be sent on the topic that
// was provided in the NewContactKey message and use the sym-key.
// Both users will therefore have the same filter.
func (c *Channel) ConfirmedContactRequest(username, image string) {
	msg := fmt.Sprintf(confirmedContactRequestMsg, username, image, c.conn.address, "")
	c.callStandardMsg(msg)
}

// ContactUpdateRequest : Sent when the user changes his name or profile-image.
func (c *Channel) ContactUpdateRequest(username, image string) {
	msg := fmt.Sprintf(contactUpdateMsg, username, image)
	c.callStandardMsg(msg)
}

// SeenRequest : Sent when a user sees a message (opens the chat and loads the
// message). Can acknowledge multiple messages at the same time
func (c *Channel) SeenRequest(ids []string) error {
	body, err := json.Marshal(ids)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(seenMsg, body)
	c.callStandardMsg(msg)

	return nil
}

func (c *Channel) callStandardMsg(body string) {
	msg := rawrChatMessage(body)

	cmd := fmt.Sprintf(standardMessageFormat,
		c.conn.address,
		c.channelKey,
		msg,
		c.topic,
		c.conn.minimumPoW)

	// TODO (adriacidre) manage this error
	_, _ = c.conn.call(cmd)
}

func (c *Channel) removeSubscription(sub *Subscription) {
	var subs []*Subscription
	for _, s := range c.subscriptions {
		if s != sub {
			subs = append(subs, s)
		}
	}
	c.subscriptions = subs
}
