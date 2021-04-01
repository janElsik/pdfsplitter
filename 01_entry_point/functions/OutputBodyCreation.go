package functions

import (
	"fmt"
	"net/http"
)

func OutputBodyCreation(w http.ResponseWriter, fileLinks []string, thumbLinks []string) {

	for count, link := range fileLinks {

		var picLink string = `
		<a href=` + link + `><img src=` + thumbLinks[count] + `></a>
	`

		fmt.Fprint(w, picLink)

	}

}
