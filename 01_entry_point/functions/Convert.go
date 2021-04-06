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

func Convert(inputFileSlice []string, tempFolderName string, tempFileName string, command <-chan string, wg *sync.WaitGroup, WeedsAddressString string, NatsAddressString string) []string {
	//connect to messaging server
	nc, err := nats.Connect(NatsAddressString)
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

	// create struct that is going to be sent to the messaging server
	type Request struct {
		ConvertToPDF     string
		Id               int
		Filename         string
		Identifier       string
		Tempfoldername   string
		Originalfilename string
	}

	// connection to filesystem
	client := weedo.NewClient(WeedsAddressString)
	client.Master()

	// create that represents each particular sent files
	UniqueString := helpers.RandomStringGenerator(12)

	// map to store links to files
	linkMap := make(map[int]string)

	// create channel to send messages
	personChanSend := make(chan *Request)
	ec.BindSendChan("request_converting_to_pdf", personChanSend)

	i := 1

	// iteration on file array that is passed to this function
	for _, inputFileName := range inputFileSlice {

		// rename the files to same name (unique string + 00001,00002, etc.)
		count := "0000" + strconv.Itoa(i)
		tempStringSlice := strings.SplitAfter(inputFileName, ".")
		tempFileSuffix := "." + tempStringSlice[len(tempStringSlice)-1]

		// Read the file into a byte slice
		originalFileByteSlice, err := os.ReadFile(inputFileName)

		if err != nil {
			fmt.Println(err)
		}
		// make sure that the folder exists
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

		// upload the file to server and get url
		fid, _, err := client.AssignUpload(file.Name(), "", file)

		if err != nil {
			fmt.Println("error with uploading the file:", err)
		}

		purl, _, err := client.GetUrl(fid)

		if err != nil {
			fmt.Println("error with getting Public url:", err)
		}
		// send message to the messaging server
		req := Request{Id: i, Filename: purl, Identifier: UniqueString, Tempfoldername: tempFolderName, Originalfilename: tempFileName + count + tempFileSuffix}
		log.Infof("Sending request id: %d with arg: %s", req.Id, req.Filename)
		personChanSend <- &req

		i++

	}
	// create struct that represents response to sent messages
	type Response struct {
		ID                 int
		Identifier         string
		Fid                string
		Originalidentifier string
	}

	// create channel to receive messages
	personChanRecv := make(chan *Response)
	_, _ = ec.BindRecvChan("response_converted_to_pdf", personChanRecv)

	tempString := strconv.Itoa(len(inputFileSlice)) + UniqueString
	fmt.Println(tempString)

	for {
		// wait for incoming messages
		deq := <-personChanRecv

		//time.Sleep(time.Second*2)

		// check if the incoming message belongs to sent message, if so, print it out and add link to map
		if deq.Originalidentifier == UniqueString {
			fmt.Println(deq.Identifier)
			fmt.Println(deq.Fid)
			linkMap[deq.ID] = deq.Fid
		}

		// check if all responses are received, if so, iterate over the map, pull out links from map,
		// sort them and return them
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
