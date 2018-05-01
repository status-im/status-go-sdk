package sdk

import (
	"fmt"
	"log"
	"time"
)

// MsgHandler is a callback function that processes messages delivered to
// asynchronous subscribers.
type MsgHandler func(msg *Msg)

// Subscription is a polling helper for a specific channel
type Subscription struct {
	unsubscribe chan bool
	channel     *Channel
}

// SubscriptionResponse json response for shh_getFilterMessages requests
type SubscriptionResponse struct {
	Result interface{} `json:"result"`
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
			var res SubscriptionResponse
			cmd := fmt.Sprintf(getFilterMessagesFormat, channel.filterID)
			if err := channel.conn.call(cmd, &res); err != nil {
				log.Fatalf("Error when sending request to server: %s", err)
				continue
			}

			switch vv := res.Result.(type) {
			case []interface{}:
				for _, u := range vv {
					payload := u.(map[string]interface{})["payload"]
					message, err := MessageFromPayload(payload.(string))
					if err != nil {
						log.Println(err)
					} else {
						if supportedMessage(message.Type) {
							fn(message)
						}
					}
				}
			default:
				log.Println(res.Result, "is of a type I don't know how to handle")
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
