package functions

import (
	"fmt"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"os"
	"pdfsplitter_cez_preprod/01_entry_point/helpers"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func Convert(inputFileSlice []string, tempFolderName string, tempFileName string, ec *nats.EncodedConn, command <-chan string, wg *sync.WaitGroup) []string {
	defer wg.Done()
	type Request struct {
		Id         int
		Filename   string
		Identifier string
	}

	UniqueString := helpers.RandomStringGenerator(12)

	linkMap := make(map[int]string)

	personChanSend := make(chan *Request)
	ec.BindSendChan("request_subject", personChanSend)

	i := 1
	for _, inputFileName := range inputFileSlice {

		count := "0000" + strconv.Itoa(i)
		//TODO if the number of documents to merge si higher than 9, add check to amend the count variable -- ex. 00001, 00020, etc.
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

		// check - if suffix is not .pdf, converts the file to pdf

		req := Request{Id: i, Filename: tempFolderName + tempFileName + count + tempFileSuffix, Identifier: UniqueString}
		log.Infof("Sending request id: %d with arg: %s", req.Id, req.Filename, req.Identifier)
		personChanSend <- &req

		i++

	}

	type Response struct {
		ID         int
		Identifier string
		Fid        string
	}

	personChanRecv := make(chan *Response)
	_, _ = ec.BindRecvChan("request_subject", personChanRecv)

	myString := strconv.Itoa(len(inputFileSlice)) + UniqueString
	fmt.Println(myString)

	for {
		// wait for incoming messages
		deq := <-personChanRecv

		//time.Sleep(time.Second*2)
		fmt.Println(deq.Identifier)
		fmt.Println(deq.Fid)
		linkMap[deq.ID] = deq.Fid

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
