package sdk

// Contact details for a known user
type Contact struct {
	account *Account
	ch      *Channel
	sub     *Subscription
	SymKey  string
	TopicID string
	Name    string
	Image   string
	Address string
	PubKey  string
}

// Accept accepts a contact as a friend, sends back a ConfirmedContactRequest
func (c *Contact) Accept() error {
	// Add the proposed symkey
	symkey, err := addSymKey(c.account.conn, c.SymKey)
	if err != nil {
		return err
	}

	c.ch, err = c.account.Join(c.Name, c.TopicID, symkey)
	if err != nil {
		return err
	}

	return c.ch.ConfirmedContactRequest(c)
}

// ContactRequestHandler handler for contact requests
type ContactRequestHandler func(*Contact)

// Subscribe to 1 to 1 messages from he current contact.
func (c *Contact) Subscribe(fn MsgHandler) (*Subscription, error) {
	c.sub.AddHook(c.PubKey, fn)
	return c.sub, nil
}

// Publish send a 1 to 1 message to the current contact.
func (c *Contact) Publish(body string) error {
	return c.ch.publish(body, "~:user-message")
}
