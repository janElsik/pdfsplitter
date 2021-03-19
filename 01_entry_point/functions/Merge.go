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

func Merge(tempFolderName string, fileLinks []string) string {
	//1. Download files to temp folder with randomly generated name+a number
	fileCount := 1
	for _, tempString := range fileLinks {
		err := helpers.DownloadFile(tempFolderName+strconv.Itoa(fileCount)+".pdf", tempString)

		if err != nil {
			fmt.Printf("Downloading error with %s : %s\n", tempString, err)

		}

		fileCount++
	}

	//2. Merge files into one
	var sliceToMerge []string
	downloadedFilesSlice, _ := ioutil.ReadDir(tempFolderName)
	for _, downloadedFile := range downloadedFilesSlice {
		sliceToMerge = append(sliceToMerge, tempFolderName+downloadedFile.Name())
	}
	err := api.MergeCreateFile(sliceToMerge, tempFolderName+"merged.pdf", nil)
	if err != nil {
		fmt.Printf("error with merging file: %s \n", err)
	}

	//3. Upload files to SeaweedFS

	file, err := os.Open(tempFolderName + "merged.pdf")
	if err != nil {
		fmt.Println("error with opening file for uploading:", err)
	}
	client := weedo.NewClient("10.0.0.27:9333")
	fid, _, err := client.AssignUpload(tempFolderName+"merged.pdf", "application/pdf", file)
	if err != nil {
		fmt.Println("error with uploading file to seaweed:", err)
	}

	publicUrl, _, err := client.GetUrl(fid)
	if err != nil {
		fmt.Println("error with getting url from seaweed:", err)
	}
	//4. Erase everything in the folder, including merged file
	err = os.RemoveAll(tempFolderName)

	if err != nil {
		fmt.Println("error with removing the temp folder:", err)
	}
	// create the folder for future purposes
	err = os.Mkdir(tempFolderName, 0777)
	if err != nil {
		fmt.Println("error with removing the temp folder:", err)
	}

	//5. Return link to the merged file
	return publicUrl

}
