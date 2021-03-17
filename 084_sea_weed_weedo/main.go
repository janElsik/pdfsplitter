package main

import (
	"fmt"
	"github.com/ginuerzh/weedo"
	"os"
)

func main() {

	client := weedo.NewClient("10.0.0.27:9333")
	//TODO na githubu je funkce, jak si rovnou načíst náhledy
	//TODO ponořit se do fid
	file, _ := os.Open("/home/jelsik/Downloads/temp-folder/gobook.pdf")
	fid, _, err := client.AssignUpload("/testfolder/gobook.pdf", "application/pdf", file)
	fmt.Println(fid)
	if err != nil {
		fmt.Println(err)
	}

	purl, url, err := client.GetUrl(fid)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(purl)
	fmt.Println(url)

	location, err := client.GetUrls(fid)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(location)

	for locPart := range location {
		fmt.Println(locPart)
	}

}
