/*
Serve is a very simple static file server in go
Usage:
	-p="8100": port to serve on
	-d=".":    the directory of static files to host
Navigating to http://localhost:8100/fs/ will display the index.html or directory
listing file.
*/
package functions

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func FileServer() {
	port := flag.String("p", "8100", "port to serve on")
	directory := flag.String("d", "/temp/", "the directory of static file to host")
	flag.Parse()

	http.Handle("/fs/", http.StripPrefix("/fs/", http.FileServer(http.Dir("/temp/"))))

	fmt.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
