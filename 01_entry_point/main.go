package main

import (
	"encoding/json"
	"fmt"
	"github.com/ginuerzh/weedo"
	"os"
	"pdfsplitter_cez_preprod/01_entry_point/functions"
	"pdfsplitter_cez_preprod/01_entry_point/helpers"
	"sync"
	"time"
)

const (
	// constants with possible file types
	InputFileName1  = "files/bmp.bmp"
	InputFileName2  = "files/txt.txt"
	InputFileName3  = "files/docx.docx"
	InputFileName4  = "files/gobook.pdf"
	InputFileName5  = "files/jpg.jpg"
	InputFileName6  = "files/ods.ods"
	InputFileName7  = "files/odt.odt"
	InputFileName8  = "files/png.png"
	InputFileName9  = "files/rtf.rtf"
	InputFileName10 = "files/xls.xls"
	InputFileName11 = "files/xlsx.xlsx"
)

type JSON struct {
	Number int    `json:"linknumber"`
	Href   string `json:"href"`
	ImgSrc string `json:"imgsrc"`
}

func main() {

	// start to track time since start of program
	start := time.Now()

	// connection to filesystem
	weedoClient := weedo.NewClient("10.0.0.27:9333")

	// array with filepaths
	inputFileSlice := []string{InputFileName1,
		InputFileName2, InputFileName3, InputFileName4,
		InputFileName5, InputFileName6, InputFileName7, InputFileName8, InputFileName9, InputFileName10, InputFileName11,
	}

	// randomly generated string used to rename the files to unique names
	tempFileName := helpers.RandomStringGenerator(12)

	// randomly generated string used to create folder with unique name
	tempFolderName := "/temp/" + helpers.RandomStringGenerator(12) + "/"

	// this block makes sure that conversion (functions.Convert) is completed before continuing with the
	// execution of the program (pointer to wg variable)
	var wg sync.WaitGroup
	wg.Add(1)
	command := make(chan string)

	// possible through go routine, but potentionally very costly regarding memory
	var linkSlice []string

	// conversion call on input files, returns Array with links converted documents
	linkSlice = functions.Convert(inputFileSlice, tempFolderName, tempFileName, command, &wg)
	err := os.RemoveAll(tempFolderName)
	if err != nil {
		fmt.Println(err)
	}

	err = os.Mkdir(tempFolderName, 0777)
	if err != nil {
		fmt.Println(err)
	}

	// waits for completion of conversion call
	wg.Wait()

	// prints name of temp folder
	fmt.Println("folder:", tempFolderName)

	// prints the links to converted documents
	for _, link := range linkSlice {

		fmt.Println(link)

	}

	// call to merge converted files into one, returns link to merged file
	mergedFileLink := functions.Merge(tempFolderName, linkSlice)

	// prints the link to merged file
	fmt.Println("link to merged file:", mergedFileLink)

	// this block makes sure that split (functions.Split) is completed before continuing with the
	// execution of the program (pointer to wg2 variable)
	var wg2 sync.WaitGroup
	wg2.Add(1)

	// call to split the merged file into single pages, returns links to split pdfs and to thumbnails of
	// the split pdfs
	thumbSlice, splitLinkSlice := functions.Split(tempFolderName, mergedFileLink, &wg2)
	wg2.Wait()

	err = os.Mkdir(tempFolderName, 0777)
	if err != nil {
		fmt.Println(err)
	}

	// create new struct for JSON
	jsonOutput := []JSON{}

	// iterate over pdf and thumbnail links and add to JSON struct
	for i, link := range splitLinkSlice {
		//fmt.Println(thumbSlice[i], link)
		fmt.Println("<a href=" + link + "><img src=" + thumbSlice[i] + "></a>")
		jsonOutput = append(jsonOutput, JSON{
			Number: i,
			Href:   link,
			ImgSrc: thumbSlice[i],
		})

	}

	// create byteArray of the JSON struct
	byteArray, err := json.Marshal(jsonOutput)
	if err != nil {
		fmt.Println("Marshaling:", err)
	}

	// write the JSON struct to tempfolder
	err = os.WriteFile(tempFolderName+"jsongo.json", byteArray, 0644)
	if err != nil {
		fmt.Println("Writing marshaled file:", err)
	}
	file, err := os.Open(tempFolderName + "jsongo.json")

	if err != nil {
		fmt.Println("Opening file:", err)
	}

	// upload JSON file, get url and print the url
	fid, _, err := weedoClient.AssignUpload("jsongo.json", "application/json", file)
	if err != nil {
		fmt.Println("Opening file:", err)
	}
	purl, _, err := weedoClient.GetUrl(fid)
	if err != nil {
		fmt.Println("Getting url:", err)
	}

	err = os.Remove(tempFolderName + "jsongo.json")
	if err != nil {
		fmt.Println("Removing file:", err)
	}

	fmt.Println("link to JSON:", purl)
	fmt.Println("link to merged file:", mergedFileLink)

	elapsed := time.Since(start)
	fmt.Println("process took:", elapsed)
}
