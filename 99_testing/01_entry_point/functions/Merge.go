package functions

import (
	"fmt"
	"github.com/ginuerzh/weedo"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"io/ioutil"
	"os"
	"pdfsplitter_cez_preprod/01_entry_point/helpers"
	"strconv"
)

func Merge(tempFolderName string, fileLinks []string, WeedsAddressString string) string {
	// Download files to temp folder with randomly generated name and a number
	fileCount := 0
	for _, tempString := range fileLinks {
		var newFileCount string
		if fileCount < 10 {
			newFileCount = "000" + strconv.Itoa(fileCount)
		} else if fileCount < 100 {
			newFileCount = "00" + strconv.Itoa(fileCount)
		} else if fileCount < 1000 {
			newFileCount = "0" + strconv.Itoa(fileCount)
		} else {
			newFileCount = strconv.Itoa(fileCount)
		}
		err := helpers.DownloadFile(tempFolderName+newFileCount+".pdf", tempString)

		if err != nil {
			fmt.Printf("Downloading error with %s : %s\n", tempString, err)

		}

		fileCount++
	}

	// Merge files into one
	var sliceToMerge []string
	downloadedFilesSlice, _ := ioutil.ReadDir(tempFolderName)
	for _, downloadedFile := range downloadedFilesSlice {
		sliceToMerge = append(sliceToMerge, tempFolderName+downloadedFile.Name())
	}
	err := api.MergeCreateFile(sliceToMerge, tempFolderName+"merged.pdf", nil)
	if err != nil {
		fmt.Printf("error with merging file: %s \n", err)
	}

	// Upload files to server
	file, err := os.Open(tempFolderName + "merged.pdf")
	if err != nil {
		fmt.Println("error with opening file for uploading:", err)
	}
	client := weedo.NewClient(WeedsAddressString)
	fid, _, err := client.AssignUpload(tempFolderName+"merged.pdf", "application/pdf", file)
	if err != nil {
		fmt.Println("error with uploading file to seaweed:", err)
	}

	// get url
	publicUrl, _, err := client.GetUrl(fid)
	if err != nil {
		fmt.Println("error with getting url from seaweed:", err)
	}
	// Erase everything in the folder, including merged file
	err = os.RemoveAll(tempFolderName)

	if err != nil {
		fmt.Println("error with removing the temp folder:", err)
	}
	// create the folder for future purposes
	err = os.Mkdir(tempFolderName, 0777)
	if err != nil {
		fmt.Println("error with removing the temp folder:", err)
	}

	// Return link to the merged file
	return publicUrl

}
