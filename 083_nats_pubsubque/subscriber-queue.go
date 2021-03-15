package main

import (
	nats "github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

func main() {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
	}
	defer ec.Close()

	log.Info("Connected to NATS and ready to receive messages")

	type Request struct {
		ID  int
		Arg string
	}
	requestChanRecv := make(chan *Request)

	// This allows us to subscribe to a queue within a subject
	// for load balancing messages among subscribers.
	// https://godoc.org/github.com/nats-io/go-nats#EncodedConn.BindRecvQueueChan
	_, _ = ec.BindRecvQueueChan("request_subject", "hello_queue", requestChanRecv)

	for {
		// Wait for incoming messages
		req := <-requestChanRecv
		if req.Arg == "convert" {
			log.Infof("Received request with no: %d and argument: %s", req.ID, req.Arg)
		}
	}
}
