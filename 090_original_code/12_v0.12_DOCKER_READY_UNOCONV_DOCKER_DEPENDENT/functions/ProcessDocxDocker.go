package functions

import (
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/thecodingmachine/gotenberg-go-client/v7"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func ProcessDocxDocker(w http.ResponseWriter, r *http.Request, directoryPath string, splitToXPages int,
	dirName string) (tempFile string) {
	// returns file from the provided form key
	file, _, err := r.FormFile("myFile")
	if err != nil {
		fmt.Printf("error with reading the form key: %v \n", err)
	}

	defer file.Close()

	// writes out html code
	w.Write([]byte(HtmlHeader))

	// makes directory with random string
	err = os.Mkdir(directoryPath+dirName, 0777)

	if err != nil {
		fmt.Printf("error with making the directory: %v \n", err)
	}

	// create the empty file and name it so that it follows the pattern
	// (in this case it will be named upload-somerandomnumber)
	TempFile, err := ioutil.TempFile(directoryPath+dirName, "*.docx")
	if err != nil {
		fmt.Printf("error with making the empty file: %v \n", err)
	}

	// read all of the contents of the uploaded form file into a byte slice
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error with reading the contents of the uploaded form file: %v \n", err)
	}

	// write this byte slice to the blank file
	_, err = TempFile.Write(fileBytes)
	if err != nil {
		fmt.Printf("error with writing the byte slice into the blank file: %v \n", err)
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	c := &gotenberg.Client{Hostname: "http://10.2.2.15:3000"}
	doc, err := gotenberg.NewDocumentFromPath(TempFile.Name(), TempFile.Name())
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("TempFile.Name():",TempFile.Name())
	//fmt.Println("dest",directoryPath + dirName + strings.ReplaceAll(TempFile.Name(), ".docx", ".pdf"))
	req := gotenberg.NewOfficeRequest(doc)
	dest := strings.ReplaceAll(TempFile.Name(), ".docx", ".pdf")
	err = c.Store(req, dest)
	if err != nil {
		fmt.Println(err)
	}

	/*
		cmd := exec.Command("unoconv", "-f", "pdf", TempFile.Name())

		if err := cmd.Run(); err != nil {
			fmt.Printf("error with converting docx to pdf: %v \n", err)
			fmt.Println(TempFile.Name())
		}

	*/

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	tempFileName := strings.ReplaceAll(TempFile.Name(), ".docx", ".pdf")

	// split the file using pdfcpu and save all the pieces to the same folder as the uploaded file
	err = api.SplitFile(tempFileName, directoryPath+dirName, splitToXPages, nil)

	if err != nil {
		fmt.Printf("error with splitting the pdf file %v \n", err)
	}
	err = os.Remove(tempFileName)
	err = os.Remove(TempFile.Name())

	files, _ := ioutil.ReadDir(directoryPath + dirName + "/")
	for _, file := range files {

		if file.IsDir() {
			continue
		}

		originalFileName := file.Name()

		tempString := strings.SplitAfter(originalFileName, "_")

		firstPart := strings.ReplaceAll(tempString[0], "_", "")
		secondPart := strings.ReplaceAll(tempString[1], ".pdf", "")

		secondPartFinal := ""

		var partNumberLength int = len(secondPart)

		switch partNumberLength {
		case 1:
			secondPartFinal = "0000" + secondPart
		case 2:
			secondPartFinal = "000" + secondPart
		case 3:
			secondPartFinal = "00" + secondPart
		case 4:
			secondPartFinal = "0" + secondPart
		case 5:
			secondPartFinal = secondPart
		}

		finalString := firstPart + secondPartFinal

		fmt.Println(finalString + ".pdf")

		err := os.Rename(directoryPath+dirName+"/"+file.Name(), directoryPath+dirName+"/"+finalString+".pdf")
		if err != nil {
			fmt.Printf("error with renaming: %v \n", err)
		}

	}

	// return this string for other processing in main (used in methods that need this file name)

	return tempFileName

}
