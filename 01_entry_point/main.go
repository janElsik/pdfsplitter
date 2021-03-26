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
	InputFileName1  = "/temp/temp-folder/bmp.bmp"
	InputFileName2  = "/temp/temp-folder/txt.txt"
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

type JSON struct {
	Number int    `json:"linknumber"`
	Href   string `json:"href"`
	ImgSrc string `json:"imgsrc"`
}

func main() {

	weedoClient := weedo.NewClient("10.0.0.27:9333")
	start := time.Now()

	inputFileSlice := []string{InputFileName1,
		InputFileName2, InputFileName3, InputFileName4,
		InputFileName5, InputFileName6, InputFileName7, InputFileName8, InputFileName9, InputFileName10, InputFileName11,
	}
	tempFileName := helpers.RandomStringGenerator(12)
	tempFolderName := "/temp/" + helpers.RandomStringGenerator(12) + "/"
	var wg sync.WaitGroup
	wg.Add(1)
	command := make(chan string)

	// possible through go routine, but potentionally very costly regarding memory
	var linkSlice []string

	linkSlice = functions.Convert(inputFileSlice, tempFolderName, tempFileName, command, &wg)
	err := os.RemoveAll(tempFolderName)
	if err != nil {
		fmt.Println(err)
	}

	err = os.Mkdir(tempFolderName, 0777)
	if err != nil {
		fmt.Println(err)
	}

	wg.Wait()

	fmt.Println("folder:", tempFolderName)

	for _, link := range linkSlice {

		fmt.Println(link)

	}

	mergedFileLink := functions.Merge(tempFolderName, linkSlice)

	fmt.Println("link to merged file:", mergedFileLink)
	var wg2 sync.WaitGroup
	wg2.Add(1)

	thumbSlice, splitLinkSlice := functions.Split(tempFolderName, mergedFileLink, &wg2)
	wg2.Wait()

	err = os.RemoveAll(tempFolderName)
	if err != nil {
		fmt.Println(err)
	}

	err = os.Mkdir(tempFolderName, 0777)
	if err != nil {
		fmt.Println(err)
	}

	jsonOutput := []JSON{}

	for i, link := range splitLinkSlice {
		//fmt.Println(thumbSlice[i], link)
		fmt.Println("<a href=" + link + "><img src=" + thumbSlice[i] + "></a>")
		jsonOutput = append(jsonOutput, JSON{
			Number: i,
			Href:   link,
			ImgSrc: thumbSlice[i],
		})

	}
	byteArray, err := json.Marshal(jsonOutput)
	if err != nil {
		fmt.Println("Marshaling:", err)
	}
	err = os.WriteFile(tempFolderName+"jsongo.json", byteArray, 0644)
	if err != nil {
		fmt.Println("Writing marshaled file:", err)
	}
	file, err := os.Open(tempFolderName + "jsongo.json")

	if err != nil {
		fmt.Println("Opening file:", err)
	}

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
