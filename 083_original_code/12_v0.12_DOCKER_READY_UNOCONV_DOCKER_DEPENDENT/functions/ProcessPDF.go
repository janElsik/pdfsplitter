package functions

import (
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func ProcessPDF(w http.ResponseWriter, r *http.Request, directoryPath string, splitToXPages int, dirName string, formFileName string) (tempFile string) {

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
	TempFile, err := ioutil.TempFile(directoryPath+dirName, "*.pdf")
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

	// split the file using pdfcpu and save all the pieces to the same folder as the uploaded file
	//TODO CHANGE FILE OUTPUT OT _XXXXX.PDF. Using this pattern will make it iterate in order when displayed
	err = api.SplitFile(TempFile.Name(), directoryPath+dirName, splitToXPages, nil)

	if err != nil {
		fmt.Printf("error with splitting the pdf file %v \n", err)
	}

	//TODO in future for cleaning - this erases the whole folder ==> first needed to move thumbs somewhere else	err = os.RemoveAll(directoryPath + dirName)
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

		wm, _ := api.TextWatermark(formFileName+"  pg."+string(secondPart), "font:Courier, points:45, sc:.6 abs, rot: 0", false, false, pdfcpu.POINTS) //points:45,col: .1 .5 0,
		err = api.AddWatermarksFile(directoryPath+dirName+"/"+file.Name(), "", nil /*[]string{"odd"}*/, wm, nil)
		if err != nil {
			panic(err)
		}

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

		err := os.Rename(directoryPath+dirName+"/"+file.Name(), directoryPath+dirName+"/"+finalString+".pdf")
		if err != nil {
			fmt.Printf("error with renaming: %v \n", err)
		}

	}

	//put the rename func here!

	// return this string for other processing in main (used in methods that need this file name)
	return TempFile.Name()

}
