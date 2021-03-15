package functions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type Links struct {
	Number int    `json:"linknumber"`
	Href   string `json:"href"`
	ImgSrc string `json:"imgsrc"`
}

func Marshal(directoryPath string, dirName string, imagesFolder string, thumbsFolder string, fileServerRoot string, tempFileName string) {
	// Empty slice
	links := []Links{}

	// Loop over filesystem
	files, _ := ioutil.ReadDir(directoryPath + dirName)

	linkCount := 1

	// this loop creates image links in the body of /upload
	for _, file := range files {

		imagesFolderAmended := strings.ReplaceAll(imagesFolder, "/", "")
		thumbsFolderAmended := strings.ReplaceAll(thumbsFolder, "/", "")

		if file.Name() == imagesFolderAmended || file.Name() == thumbsFolderAmended || file.Name() == tempFileName {
			continue
		}

		thumbPath := strings.ReplaceAll(file.Name(), ".pdf", ".png")

		href := fileServerRoot + dirName + "/" + file.Name()
		imgSrc := fileServerRoot + dirName + thumbsFolder + thumbPath

		links = append(links, Links{Number: linkCount, Href: href, ImgSrc: imgSrc})

		linkCount++
	}

	ByteArray, err := json.Marshal(links)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(directoryPath+dirName+"/"+`jsongo.json`, ByteArray, 0644)
	if err != nil {
		fmt.Println(err)
	}

}
