package functions

import (
	"fmt"
	"github.com/ginuerzh/weedo"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"os"
	"pdfsplitter_cez_preprod/01_entry_point/helpers"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func Convert(inputFileSlice []string, tempFolderName string, tempFileName string, command <-chan string, wg *sync.WaitGroup) []string {
	nc, err := nats.Connect("10.0.0.27:4222")
	if err != nil {
		panic(err)
	}

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
	}
	defer ec.Close()
	log.Info("Connected to NATS and ready to send messages")

	defer wg.Done()
	type Request struct {
		ConvertToPDF     string
		Id               int
		Filename         string
		Identifier       string
		Tempfoldername   string
		Originalfilename string
	}

	// initialize SeaWeedFS
	client := weedo.NewClient("10.0.0.27:9333")
	client.Master()

	UniqueString := helpers.RandomStringGenerator(12)

	linkMap := make(map[int]string)

	personChanSend := make(chan *Request)
	ec.BindSendChan("request_converting_to_pdf", personChanSend)

	i := 1
	for _, inputFileName := range inputFileSlice {

		count := "0000" + strconv.Itoa(i)
		tempStringSlice := strings.SplitAfter(inputFileName, ".")
		tempFileSuffix := "." + tempStringSlice[len(tempStringSlice)-1]

		// Read the file into a byte slice
		originalFileByteSlice, err := os.ReadFile(inputFileName)

		if err != nil {
			fmt.Println(err)
		}
		// make sure that the folder exists - either creates it or throws error, which we do not handle
		err = os.Mkdir(tempFolderName, 0777)

		if err != nil {
			//	fmt.Println(err)
		}
		// Write the byte slice into new file
		err = os.WriteFile(tempFolderName+tempFileName+count+tempFileSuffix, originalFileByteSlice, 0777)

		if err != nil {
			fmt.Println(err)
		}
		file, _ := os.Open(tempFolderName + tempFileName + count + tempFileSuffix)
		fid, _, err := client.AssignUpload(file.Name(), "", file)

		if err != nil {
			fmt.Println("error with uploading the file:", err)
		}

		purl, _, err := client.GetUrl(fid)

		if err != nil {
			fmt.Println("error with getting Public url:", err)
		}

		req := Request{Id: i, Filename: purl, Identifier: UniqueString, Tempfoldername: tempFolderName, Originalfilename: tempFileName + count + tempFileSuffix}
		log.Infof("Sending request id: %d with arg: %s", req.Id, req.Filename)
		personChanSend <- &req

		i++

	}

	type Response struct {
		ID                 int
		Identifier         string
		Fid                string
		Originalidentifier string
	}

	personChanRecv := make(chan *Response)
	_, _ = ec.BindRecvChan("response_converted_to_pdf", personChanRecv)

	myString := strconv.Itoa(len(inputFileSlice)) + UniqueString
	fmt.Println(myString)

	for {
		// wait for incoming messages
		deq := <-personChanRecv

		//time.Sleep(time.Second*2)
		if deq.Originalidentifier == UniqueString {
			fmt.Println(deq.Identifier)
			fmt.Println(deq.Fid)
			linkMap[deq.ID] = deq.Fid
		}

		if len(linkMap) == len(inputFileSlice) && linkMap[len(linkMap)] != "" {

			keys := make([]int, 0, len(linkMap))
			values := make([]string, 0, len(linkMap))
			for k := range linkMap {
				keys = append(keys, k)

			}
			sort.Ints(keys)

			for _, k := range keys {
				fmt.Println(k, linkMap[k])
				values = append(values, linkMap[k])
			}

			return values
		}
	}

}
