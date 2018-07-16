package sdk

import (
	"time"
)

// MsgHandler is a callback function that processes messages delivered to
// asynchronous subscribers.
type MsgHandler func(msg *Msg)

// Subscription is a polling helper for a specific channel
type Subscription struct {
	unsubscribe chan bool
	channel     *Channel
	hooks       map[string]MsgHandler
}

// Subscribe polls on specific channel topic and executes given function if
// any message is received
func (s *Subscription) Subscribe(channel *Channel, fn MsgHandler) {
	s.channel = channel
	for {
		select {
		case <-s.unsubscribe:
			return
		default:
			if m := channel.pollMessages(); m != nil {
				if properties, ok := m.Properties.(*PublishMsg); ok {
					if properties.Visibility == "~:user-message" {
						for k := range s.hooks {
							println(k)
						}
						if hook, ok := s.hooks[m.PubKey]; ok {
							hook(m)
							continue
						}
					}
				}
				fn(m)
			}
		}
		// TODO(adriacidre) : move this period to configuration
		time.Sleep(time.Second * 3)
	}
}

// Unsubscribe stops polling on the current subscription channel
func (s *Subscription) Unsubscribe() {
	s.unsubscribe <- true
	s.channel.removeSubscription(s)
}

// AddHook adds a hook to the current subscription to be exectuted when the
// message is received from a specific address.
func (s *Subscription) AddHook(pubkey string, fn MsgHandler) {
	if s.hooks == nil {
		s.hooks = make(map[string]MsgHandler)
	}
	s.hooks[pubkey] = fn
}
