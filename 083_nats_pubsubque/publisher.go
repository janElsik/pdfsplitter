package main

import (
	"time"

	"github.com/nats-io/nats.go"
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

	log.Info("Connected to NATS and ready to send messages")

	type Request struct {
		ID  int
		Arg string
	}
	requestChanSend := make(chan *Request)
	ec.BindSendChan("request_subject", requestChanSend)

	i := 0
	for {

		req := Request{ID: i, Arg: "convert"}
		// Create instance of type Request with Id and set to
		// the current value of i
		if i%2 == 0 {
			req = Request{ID: i, Arg: "trevnoc"}
		}

		// Just send to the channel! :)
		log.Infof("Sending request: %d with arg: %s", req.ID, req.Arg)
		requestChanSend <- &req

		// Pause and increment counter
		time.Sleep(time.Second * 1)
		i = i + 13
	}
}
