package functions

import (
	"fmt"
	"net/http"
	"strings"
)

func OutputBodyCreation(w http.ResponseWriter, fileLinks []string, thumbLinks []string) {

	for count, link := range fileLinks {

		newLink := strings.ReplaceAll(link, "10.4.129.95", "10.4.56.20")
		newThumbLink := strings.ReplaceAll(thumbLinks[count], "10.4.129.95", "10.4.56.20")
		//http://10.4.129.95:8080/2,065b69a46996

		var picLink string = `
		<a href=` + newLink + `><img src=` + newThumbLink + `></a>
	`

		fmt.Fprint(w, picLink)

	}

}
