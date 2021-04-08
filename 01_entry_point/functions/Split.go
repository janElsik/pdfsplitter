package functions

import (
	"fmt"
	"github.com/ginuerzh/weedo"
	"github.com/nats-io/nats.go"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"pdfsplitter_cez_preprod/01_entry_point/helpers"
	"sort"
	"strings"
	"sync"
)

func Split(tempFolderName string, mergedFileLink string, wg *sync.WaitGroup, WeedsAddressString string, NatsAddressString string) ([]string, []string) {
	defer wg.Done()

	// create struct to receive messages
	type ThumbCreateRequest struct {
		Maxnumber    int
		Createthumbs string
		Id           int
		Filelink     string
		Foldername   string
		Identifier   string
	}
	// create struct to send messages
	type ThumbList struct {
		ThumbLink  string
		Id         int
		Identifier string
	}

	// identifier to be sent with request message
	identifier := helpers.RandomStringGenerator(12)
	nc, err := nats.Connect(NatsAddressString)
	if err != nil {
		panic(err)
	}

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
	}
	defer ec.Close()
	log.Info("Split PDF: Connected to NATS and ready to send messages")
	// Download the merged file
	err = helpers.DownloadFile(tempFolderName+"merged.pdf", mergedFileLink)
	if err != nil {
		fmt.Println("error with downloading the file:", err)
	}

	// Split the file into multiple 1 page files
	err = api.SplitFile(tempFolderName+"merged.pdf", tempFolderName, 1, nil)
	if err != nil {
		fmt.Println("error with splitting the file:", err)
	}
	// Delete the downloaded file
	err = os.Remove(tempFolderName + "merged.pdf")
	if err != nil {
		fmt.Println(err)
	}
	// Loop through the directory and correct numbering of pages
	files, err := ioutil.ReadDir(tempFolderName)

	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {

		tempSlice := strings.Split(file.Name(), "_")
		firstStringPart := tempSlice[0]
		secondStringPart := tempSlice[1]
		secondStringPart = strings.ReplaceAll(secondStringPart, ".pdf", "")
		if len(secondStringPart) == 1 {
			secondStringPart = "0000" + secondStringPart + ".pdf"
		} else if len(secondStringPart) == 2 {
			secondStringPart = "000" + secondStringPart + ".pdf"
		} else if len(secondStringPart) == 3 {
			secondStringPart = "00" + secondStringPart + ".pdf"
		} else {
			secondStringPart = "0" + secondStringPart + ".pdf"
		}
		err = os.Rename(tempFolderName+file.Name(), tempFolderName+firstStringPart+secondStringPart)

		if err != nil {
			fmt.Println("error with renaming file:", err)
		}
	}

	// upload to server and return links to a single Array
	client := weedo.NewClient(WeedsAddressString)

	var splitLinkList []string

	files, err = ioutil.ReadDir(tempFolderName)
	if err != nil {
		fmt.Println(err)
	}

	for _, fileNameString := range files {
		fileToUpload, err := os.Open(tempFolderName + fileNameString.Name())
		if err != nil {
			fmt.Println(err)
		}
		fid, _, err := client.AssignUpload(tempFolderName+fileNameString.Name(), "application/pdf", fileToUpload)
		if err != nil {
			fmt.Println("error with upload:", err)
		}
		purl, _, err := client.GetUrl(fid)

		if err != nil {
			fmt.Println("error with getting fid:", err)
		}
		splitLinkList = append(splitLinkList, purl)

	}

	personChanSend := make(chan *ThumbCreateRequest)
	ec.BindSendChan("request_thumb_creation", personChanSend)

	// For each uploaded file, send message with file link, have it create thumbnail and return the link
	// return the link for file and link for thumbnail (dont forget to keep order of the files and thumbnails)

	for i, link := range splitLinkList {
		//		fmt.Println(i, link)
		req := ThumbCreateRequest{
			Maxnumber:  len(splitLinkList),
			Id:         i,
			Filelink:   link,
			Foldername: tempFolderName,
			Identifier: identifier,
		}

		//		log.Infof("Sending request to create thumbs with id %d and link %s", i, link)

		personChanSend <- &req
	}

	// create channel to receive response
	personChanRecv := make(chan *ThumbList)
	_, _ = ec.BindRecvChan("create_thumb_response", personChanRecv)
	linkMap := make(map[int]string)
	for {
		deq := <-personChanRecv

		// check if the incoming message belongs to sent message, if so, print it out and add link to map
		if deq.Identifier == identifier {
			linkMap[deq.Id] = deq.ThumbLink
			//			fmt.Println(deq.Id)
		}

		// check if all responses are received, if so, iterate over the map, pull out links from map,
		// sort them and return them
		if len(linkMap) == len(splitLinkList) && linkMap[len(linkMap)-1] != "" {
			keys := make([]int, 0, len(linkMap))
			values := make([]string, 0, len(linkMap))
			for k := range linkMap {
				keys = append(keys, k)

			}
			sort.Ints(keys)

			for _, k := range keys {
				//				fmt.Println(k, linkMap[k])
				values = append(values, linkMap[k])
			}

			return values, splitLinkList
		}

	}

}
