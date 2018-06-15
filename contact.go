package sdk

// Contact details for a known user
type Contact struct {
	account *Account
	ch      *Channel
	SymKey  string
	TopicID string
	Name    string
	Image   string
	Address string
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
