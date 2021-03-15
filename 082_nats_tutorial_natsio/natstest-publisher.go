package main

import (
	nats "github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"time"
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
		Id  int
		Arg string
	}

	personChanSend := make(chan *Request)
	ec.BindSendChan("request_subject", personChanSend)
	i := 8

	for {
		req := Request{Id: i, Arg: "convert"}
		// Create instance of type Request with Id and set to
		// the current value of i
		if i%2 == 0 {
			req = Request{Id: i, Arg: "trevnoc"}
		}

		// Send to the channel

		//log.Infof("Sending request %d", req.Id)
		log.Infof("Sending request id: %d with arg: %s", req.Id, req.Arg)
		personChanSend <- &req

		// Pause and increment counter
		time.Sleep(time.Second * 1)

		i += 15

	}

}
