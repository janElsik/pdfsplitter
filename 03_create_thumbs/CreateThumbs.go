package main

import (
	"fmt"
	"github.com/ginuerzh/weedo"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"pdfsplitter_cez_preprod/03_create_thumbs/helpers"
	"strconv"
	"strings"
	"time"
)

const thumbX1 = "200"

func main() {

	//connect to messaging server and to the file server

	client := weedo.NewClient("10.0.0.27:9333")

	nc, err := nats.Connect("10.0.0.27:4222")
	if err != nil {
		fmt.Println(err)

	}

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		fmt.Println(err)

	}
	defer ec.Close()

	log.Info("Create Thumbs: Conected to NATS and ready to receive messages")

	// create struct to receive incoming messages

	type ThumbCreateRequest struct {
		Maxnumber    int
		Createthumbs string
		Id           int
		Filelink     string
		Foldername   string
		Identifier   string
	}

	//create struct to send messages

	type ThumbList struct {
		ThumbLink  string
		Id         int
		Identifier string
	}
	// enter listening queue and for loop
	personChanRecv := make(chan *ThumbCreateRequest)
	_, _ = ec.BindRecvQueueChan("request_thumb_creation", "request_thumb_creation_queue", personChanRecv)
	personChanSend := make(chan *ThumbList)
	_ = ec.BindSendChan("create_thumb_response", personChanSend)

	// receive request, download file, create thumbnail, upload and send response with url
	for {
		req := <-personChanRecv
		if req.Id == 0 && req.Filelink == "" && req.Foldername == "" {
			continue
		}
		err = os.Mkdir("/temp", 0777)
		if err != nil {
			//	fmt.Println(err)
		}

		err = os.Mkdir(req.Foldername, 0777)
		if err != nil {
			//	fmt.Println(err)
		}

		err = helpers.DownloadFile(req.Foldername+strconv.Itoa(req.Id)+".pdf", req.Filelink)
		if err != nil {
			fmt.Println("Downloading file", err)
		}
		originalFileName := req.Foldername + strconv.Itoa(req.Id) + ".pdf"
		pngFileName := strings.ReplaceAll(originalFileName, ".pdf", ".png")

		cmd := exec.Command("mutool", "draw", "-N", "-o", pngFileName, "-h", thumbX1, "-F", "png", originalFileName)
		if err := cmd.Run(); err != nil {
			fmt.Printf("error inside exec pdf to png conversion: %v,\n", err)
		}

		file, _ := os.Open(pngFileName)

		fid, _, err := client.AssignUpload(req.Foldername+strconv.Itoa(req.Id), "image/png", file)
		if err != nil {
			fmt.Println(err)
		}

		purl, _, err := client.GetUrl(fid)

		if err != nil {
			fmt.Println(err)
		}

		_ = os.Remove(originalFileName)
		_ = os.Remove(pngFileName)

		time.Sleep(time.Millisecond * 30)

		deq := ThumbList{
			ThumbLink:  purl,
			Id:         req.Id,
			Identifier: req.Identifier,
		}
		personChanSend <- &deq
		fmt.Println("id", req.Id, "sent")

	}

}
