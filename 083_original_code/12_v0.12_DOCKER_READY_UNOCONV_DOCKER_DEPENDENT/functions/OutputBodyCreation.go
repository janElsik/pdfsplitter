package functions

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func OutputBodyCreation(directoryPath string, dirName string, imagesFolder string, thumbsFolder string, w http.ResponseWriter, fileServerRoot string, tempFileName string) {

	files, _ := ioutil.ReadDir(directoryPath + dirName)

	// this loop creates image links in the body of /upload
	for _, file := range files {

		imagesFolderAmended := strings.ReplaceAll(imagesFolder, "/", "")
		thumbsFolderAmended := strings.ReplaceAll(thumbsFolder, "/", "")

		if file.Name() == imagesFolderAmended || file.Name() == thumbsFolderAmended || file.Name() == tempFileName {
			continue
		}

		var thumbPath string = strings.ReplaceAll(file.Name(), ".pdf", ".png")
		var picLink string = `
		<a href=` + fileServerRoot + dirName + "/" + file.Name() + `><img src=` + fileServerRoot + dirName + thumbsFolder + thumbPath + `></a>
	`

		fmt.Fprint(w, picLink)

	}

}
