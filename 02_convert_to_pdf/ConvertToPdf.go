package main

import (
	"fmt"
	"github.com/ginuerzh/weedo"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	// initialize SeaWeedFS
	client := weedo.NewClient("10.0.0.27:9333")

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
	_, _ = ec.BindRecvQueueChan("request_subject", "hello_queue", personChanRecv)

	// Response struct - it is needed to listen to the incoming messages
	type Response struct {
		ID         int
		Identifier string
		Fid        string
	}

	// create channel to listen  for the messages
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

		//log.Infof("Received request with no: %d and argument: %s", req.Id, req.Filename)

		if strings.HasSuffix(req.Filename, "pdf") == false {
			cmd := exec.Command("unoconv", "-f", "pdf", req.Filename)

			if err := cmd.Run(); err != nil {
				fmt.Printf("error with converting to pdf: %v \n", err)
				fmt.Println(req.Filename)
			}
		}

		// in the messages, the received file name is with original suffix (ex. ".jpg" etc)
		// this here makes sure that we upload the correct converted file, with ".pdf" suffix
		pdfFileName := strings.Split(req.Filename, ".")
		pdfFileName[len(pdfFileName)-1] = ".pdf"
		stringToUpload := strings.Join(pdfFileName, "")

		file, _ := os.Open(stringToUpload)
		fid, _, err := client.AssignUpload(pdfFileName[len(pdfFileName)-2]+pdfFileName[len(pdfFileName)-1], "application/pdf", file)
		//fmt.Println(fid)
		if err != nil {
			fmt.Println(err)
		}

		purl, _, err := client.GetUrl(fid)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(purl)
		if err != nil {
			fmt.Println(err)
		}

		os.Remove(req.Filename)
		s = strconv.Itoa(req.Id) + req.Identifier

		deq := Response{ID: req.Id, Identifier: s, Fid: purl}
		personChanSend <- &deq
		time.Sleep(time.Microsecond * 20)

	}

}
