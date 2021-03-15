package functions

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

func ThumbsCreation(directoryPath string, dirName string, thumbX1 string) {

	files, _ := ioutil.ReadDir(directoryPath + dirName)
	start := time.Now()

	for _, file := range files {
		fmt.Println(directoryPath + dirName + "/" + file.Name())
		pdfFileName := directoryPath + dirName + "/" + file.Name()
		pdfFileNameWithoutPDF := directoryPath + dirName + "/THUMBS" + file.Name()
		pdfFileNameWithoutPDF = strings.ReplaceAll(pdfFileNameWithoutPDF, ".pdf", ".png")

		cmd := exec.Command("mutool", "draw", "-N", "-o", pdfFileNameWithoutPDF, "-h", thumbX1, "-F", "png", pdfFileName)
		if err := cmd.Run(); err != nil {
			//log.Fatal(err)
			fmt.Printf("error inside exec pdf to png conversion: %v,\n", err)
		}
		//f, err := os.Create(pdfFileNameWithoutPDF + ".png")
		fmt.Println(pdfFileNameWithoutPDF + " created.")

		/*	err = png.Encode(f, img)
			if err != nil {
				fmt.Printf("fitz error4: %v \n", err)
			}

			f.Close()


		*/
	}
	elapsed := time.Since(start)
	fmt.Println("thumbs creation:", elapsed)

}
