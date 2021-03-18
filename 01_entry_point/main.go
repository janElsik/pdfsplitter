package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"os"
	"pdfsplitter_cez_preprod/01_entry_point/functions"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	InputFileName1  = "/temp/temp-folder/bmp.bmp"
	InputFileName2  = "/temp/temp-folder/csv.csv"
	InputFileName3  = "/temp/temp-folder/docx.docx"
	InputFileName4  = "/temp/temp-folder/gobook.pdf"
	InputFileName5  = "/temp/temp-folder/jpg.jpg"
	InputFileName6  = "/temp/temp-folder/ods.ods"
	InputFileName7  = "/temp/temp-folder/odt.odt"
	InputFileName8  = "/temp/temp-folder/png.png"
	InputFileName9  = "/temp/temp-folder/rtf.rtf"
	InputFileName10 = "/temp/temp-folder/xls.xls"
	InputFileName11 = "/temp/temp-folder/xlsx.xlsx"
)

func Convert(inputFileSlice []string, tempFolderName string, tempFileName string, ec *nats.EncodedConn, command <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	type Request struct {
		Id         int
		Filename   string
		Identifier string
	}

	UniqueString := functions.RandomStringGenerator(12)

	personChanSend := make(chan *Request)
	ec.BindSendChan("request_subject", personChanSend)

	i := 1
	os.Mkdir(tempFolderName+"converted", 0777)
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
		if tempFileSuffix != ".pdf" {

			req := Request{Id: i, Filename: tempFolderName + tempFileName + count + tempFileSuffix, Identifier: UniqueString}
			log.Infof("Sending request id: %d with arg: %s", req.Id, req.Filename, req.Identifier)
			personChanSend <- &req

		}

		i++

	}

	type Response struct {
		Identifier string
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
		if myString == deq.Identifier {

			return
		}
	}

}

func main() {
	start := time.Now()
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

	inputFileSlice := []string{InputFileName1, InputFileName2, InputFileName3, InputFileName4,
		InputFileName5, InputFileName6, InputFileName7, InputFileName8, InputFileName9, InputFileName10, InputFileName11}
	tempFileName := functions.RandomStringGenerator(12)
	tempFolderName := "/temp/" + functions.RandomStringGenerator(12) + "/"
	var wg sync.WaitGroup
	wg.Add(1)
	command := make(chan string)

	go Convert(inputFileSlice, tempFolderName, tempFileName, ec, command, &wg)

	wg.Wait()

	fmt.Println("folder:", tempFolderName)

	elapsed := time.Since(start)
	fmt.Println("process took:", elapsed)
}
