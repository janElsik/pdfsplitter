package main

import (
	"fmt"
	"github.com/ginuerzh/weedo"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"pdfsplitter_cez_preprod/02_convert_to_pdf/helpers"
	"strconv"
	"strings"
	"time"
)

func main() {
	// initialize SeaWeedFS
	client := weedo.NewClient("10.0.0.27:9333")

	nc, err := nats.Connect("10.0.0.27:4222")
	if err != nil {
		panic(err)
	}

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
	}
	defer ec.Close()

	log.Info("Convert to pdf: Conected to NATS and ready to receive messages")

	// Make sure this type and its properties are exported
	// so the serializer doesn't bork
	type Request struct {
		ConvertToPDF     string
		Id               int
		Filename         string
		Identifier       string
		Tempfoldername   string
		Originalfilename string
	}
	personChanRecv := make(chan *Request)
	_, _ = ec.BindRecvQueueChan("request_converting_to_pdf", "request_converting_to_pdf_queue", personChanRecv)

	// Response struct - it is needed to listen to the incoming messages
	type Response struct {
		ID                 int
		Identifier         string
		Fid                string
		Originalidentifier string
	}

	// create channel to listen  for the messages
	personChanSend := make(chan *Response)
	ec.BindSendChan("response_converted_to_pdf", personChanSend)
	s := ""
	fmt.Println(s)

	for {
		// wait for incoming messages

		req := <-personChanRecv
		if req.Filename == "" {
			continue
		}
		err = os.Mkdir("/temp", 0777)
		if err != nil {
			fmt.Println("1 error with making directory", err)
		}
		err = os.Mkdir(req.Tempfoldername, 0777)
		if err != nil {
			fmt.Println("2 error with making directory", err)
		}
		//newLink := req.Filename
		//newLink = strings.ReplaceAll(newLink,"172.21.0.3","0.0.0.0")

		log.Infof("Received request with no: %d and argument: %s", req.Id, req.Filename)
		err = helpers.DownloadFile(req.Tempfoldername+req.Originalfilename, req.Filename)
		if err != nil {
			fmt.Println("Downloading file:", err)
		}
		fmt.Println("file downloaded")
		fileName := req.Tempfoldername + req.Originalfilename
		fmt.Println("fileNamecreated")

		//TODO office files convert through gotenberg
		if strings.HasSuffix(fileName, "pdf") == false {
			cmd := exec.Command("unoconv", "-f", "pdf", fileName)

			if err := cmd.Run(); err != nil {
				fmt.Printf("error with converting to pdf: %v \n", err)
				fmt.Println(fileName)
			}
		}

		// in the messages, the received file name is with original suffix (ex. ".jpg" etc)
		// this here makes sure that we upload the correct converted file, with ".pdf" suffix
		pdfFileName := strings.Split(fileName, ".")
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

		os.Remove(fileName)
		s = strconv.Itoa(req.Id) + req.Identifier

		deq := Response{ID: req.Id, Identifier: s, Fid: purl, Originalidentifier: req.Identifier}
		personChanSend <- &deq
		time.Sleep(time.Microsecond * 20)

	}

}
