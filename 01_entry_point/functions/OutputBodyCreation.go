package functions

import (
	"fmt"
	"net/http"
	"strings"
)

func OutputBodyCreation(w http.ResponseWriter, fileLinks []string, thumbLinks []string, WeedsVolumeVisibleAddress string) {

	//////////////////////////////////////////////////////////*
	//
	//sampleString := "http://10.4.129.95:8080/2,065b69a46996"
	//
	//
	//
	// sampleSlice := strings.SplitAfterN(sampleString,"/",4)
	//
	// fmt.Println(sampleSlice[len(sampleSlice)-1])
	//
	//*////////////////////////////////////////////////////////

	for count, link := range fileLinks {

		tempLinkSlice1 := strings.SplitAfterN(link, "/", 4)
		newLink := WeedsVolumeVisibleAddress + tempLinkSlice1[len(tempLinkSlice1)-1]

		//strings.ReplaceAll(link, "10.4.129.95", "10.4.56.20")

		tempLinkSlice2 := strings.SplitAfterN(thumbLinks[count], "/", 4)
		newThumbLink := WeedsVolumeVisibleAddress + tempLinkSlice2[len(tempLinkSlice2)-1]
		//strings.ReplaceAll(thumbLinks[count], "10.4.129.95", "10.4.56.20")

		var picLink string = `
		<a href=` + newLink + `><img src=` + newThumbLink + `></a>
	`

		fmt.Fprint(w, picLink)

	}

}
