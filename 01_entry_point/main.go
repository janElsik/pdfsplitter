package main

import (
	"fmt"
	"os"
	"pdfsplitter_cez_preprod/01_entry_point/functions"
	"pdfsplitter_cez_preprod/01_entry_point/helpers"
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

func main() {
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

	/*	err = os.Mkdir(tempFolderName, 0777)
		if err != nil {
			fmt.Println(err)
		}

	*/

	for i, link := range splitLinkSlice {
		//fmt.Println(thumbSlice[i], link)
		fmt.Println("<a href=" + link + "><img src=" + thumbSlice[i] + "></a>")
	}
	elapsed := time.Since(start)
	fmt.Println("process took:", elapsed)
}
