package main

import (
	"fmt"
	"github.com/ginuerzh/weedo"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/thecodingmachine/gotenberg-go-client/v7"
	"os"
	"os/exec"
	"pdfsplitter_cez_preprod/02_convert_to_pdf/helpers"
	"strconv"
	"strings"
	"time"
)

func main() {
	//connect to messaging server, to the file server and to conversion container
	weedoClient := weedo.NewClient("10.0.0.27:9333")
	gotenbergClient := &gotenberg.Client{Hostname: "http://10.0.0.27:3000"}

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

	// create struct for received messages
	type Request struct {
		ConvertToPDF     string
		Id               int
		Filename         string
		Identifier       string
		Tempfoldername   string
		Originalfilename string
	}
	//create channel for received messages and join queue
	personChanRecv := make(chan *Request)
	_, _ = ec.BindRecvQueueChan("request_converting_to_pdf", "request_converting_to_pdf_queue", personChanRecv)

	// create struct for outgoing messages (response)
	type Response struct {
		ID                 int
		Identifier         string
		Fid                string
		Originalidentifier string
	}

	// create channel for sending messages
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
		// make sure that directory exists, if not, create it
		err = os.Mkdir("/temp", 0777)
		if err != nil {
			//fmt.Println("1 error with making directory", err)
		}
		err = os.Mkdir(req.Tempfoldername, 0777)
		if err != nil {
			//fmt.Println("2 error with making directory", err)
		}
		//newLink := req.Filename
		//newLink = strings.ReplaceAll(newLink,"172.21.0.3","0.0.0.0")

		//log.Infof("Received request with no: %d and argument: %s", req.Id, req.Filename)

		// download file
		err = helpers.DownloadFile(req.Tempfoldername+req.Originalfilename, req.Filename)
		if err != nil {
			fmt.Println("Downloading file:", err)
		}
		fileName := req.Tempfoldername + req.Originalfilename

		// check the file suffix - if pdf, do not convert
		if strings.HasSuffix(fileName, "pdf") == false {

			// check the file suffix, if it is office doc, send for conversion to container service (gotenberg)
			// if not, convert using unoconv
			if strings.HasSuffix(fileName, ".docx") == true ||
				strings.HasSuffix(fileName, ".doc") == true ||
				strings.HasSuffix(fileName, ".ods") == true ||
				strings.HasSuffix(fileName, ".odt") == true ||
				strings.HasSuffix(fileName, ".rtf") == true ||
				strings.HasSuffix(fileName, ".xls") == true ||
				strings.HasSuffix(fileName, ".xlsx") == true ||
				strings.HasSuffix(fileName, ".txt") == true {
				doc, err := gotenberg.NewDocumentFromPath(fileName, fileName)
				if err != nil {
					fmt.Println(fileName, err)
				}
				fileNameSlice := strings.Split(fileName, ".")
				suffixOld := fileNameSlice[len(fileNameSlice)-1]
				req := gotenberg.NewOfficeRequest(doc)
				dest := strings.ReplaceAll(fileName, "."+suffixOld, ".pdf")
				err = gotenbergClient.Store(req, dest)
				if err != nil {
					fmt.Println(fileName, err)
				}
			} else {

				cmd := exec.Command("unoconv", "-f", "pdf", fileName)

				if err := cmd.Run(); err != nil {
					fmt.Printf("error with converting to pdf: %v \n", err)
					fmt.Println(fileName)
				}
			}
		}

		// in the messages, the received file name is with original suffix (ex. ".jpg" etc)
		// this makes sure that we upload the correct converted file, with ".pdf" suffix
		pdfFileName := strings.Split(fileName, ".")
		pdfFileName[len(pdfFileName)-1] = ".pdf"
		stringToUpload := strings.Join(pdfFileName, "")

		file, _ := os.Open(stringToUpload)
		fid, _, err := weedoClient.AssignUpload(pdfFileName[len(pdfFileName)-2]+pdfFileName[len(pdfFileName)-1], "application/pdf", file)
		//fmt.Println(fid)
		if err != nil {
			fmt.Println(err)
		}
		// get url
		purl, _, err := weedoClient.GetUrl(fid)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(purl)
		if err != nil {
			fmt.Println(err)
		}

		os.Remove(fileName)
		s = strconv.Itoa(req.Id) + req.Identifier
		// send response
		deq := Response{ID: req.Id, Identifier: s, Fid: purl, Originalidentifier: req.Identifier}
		personChanSend <- &deq
		time.Sleep(time.Microsecond * 20)

	}

}
