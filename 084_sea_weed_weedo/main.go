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
	//TODO fix sorting bug
	file, _ := os.Open("/home/jelsik/Downloads/temp-folder/a.jpg")
	fid, _, err := client.AssignUpload("a.jpg", "image/jpeg", file)
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

	if err != nil {
		fmt.Println(err)
	}

}
