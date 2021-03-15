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

	log.Info("Conected to NATS and ready to receive messages")

	// Make sure this type and its properties are exported
	// so the serializer doesn't bork
	type Request struct {
		Id  int
		Arg string
	}
	personChanRecv := make(chan *Request)
	_, _ = ec.BindRecvChan("request_subject", personChanRecv)

	for {
		// wait for incoming messages
		req := <-personChanRecv
		if req.Arg == "convert" {
			log.Infof("Received request with no: %d and argument: %s", req.Id, req.Arg)
		}
		//time.Sleep(time.Second*2)
	}

}
