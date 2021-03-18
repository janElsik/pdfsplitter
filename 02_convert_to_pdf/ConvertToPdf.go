package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strconv"
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
		Id         int
		Filename   string
		Identifier string
	}
	personChanRecv := make(chan *Request)
	_, _ = ec.BindRecvChan("request_subject", personChanRecv)

	type Response struct {
		Identifier string
	}

	personChanSend := make(chan *Response)
	ec.BindSendChan("request_subject", personChanSend)
	s := ""
	fmt.Println(s)
	for {
		// wait for incoming messages

		req := <-personChanRecv
		if req.Filename == "" {
			continue
		}

		log.Infof("Received request with no: %d and argument: %s", req.Id, req.Filename)

		cmd := exec.Command("unoconv", "-f", "pdf", req.Filename)

		if err := cmd.Run(); err != nil {
			fmt.Printf("error with converting to pdf: %v \n", err)
			fmt.Println(req.Filename)
		}

		os.Remove(req.Filename)
		s = strconv.Itoa(req.Id) + req.Identifier

		deq := Response{Identifier: s}
		personChanSend <- &deq

		//time.Sleep(time.Second*2)
	}

}
