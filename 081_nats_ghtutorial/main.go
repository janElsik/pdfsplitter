package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

func main() {

	// Connect to a server
	nc, _ := nats.Connect(nats.DefaultURL)

	// Simple Publisher
	err := nc.Publish("foo", []byte("Hello World"))
	if err != nil {
		fmt.Printf("following error occured with //simple publisher: %s", err)
	}

	// Simple Async Subscriber
	_, _ = nc.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})

	// Responding to a request message
	_, _ = nc.Subscribe("reqiest", func(m *nats.Msg) {
		_ = m.Respond([]byte("answer is 42"))
	})

	// Simple Sync Subscriber
	sub, err := nc.SubscribeSync("foo")
	_, err = sub.NextMsg(2 * time.Second)

	// Channel Subscriber
	ch := make(chan *nats.Msg, 64)
	sub, err = nc.ChanSubscribe("foo", ch)

	// Unscubscribe
	sub.Unsubscribe()

	// Drain
	sub.Drain()

	// Requests
	_, err = nc.Request("help", []byte("help me"), 10*time.Millisecond)

	// Replies
	nc.Subscribe("help", func(m *nats.Msg) {
		nc.Publish(m.Reply, []byte("I can help!"))
	})

	// Drain connection (Preferred for responders)
	// Close() not needed if this is called.
	nc.Drain()

	// Close connection
	nc.Close()

}
