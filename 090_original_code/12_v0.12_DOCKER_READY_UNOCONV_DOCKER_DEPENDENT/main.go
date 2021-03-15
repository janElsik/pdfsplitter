package main

import (
	"fmt"
	"net/http"
	"os"
	"pdf/12_v0.12_SERVER_READY_direct_file_serving/functions"
	"strings"
	"time"
)

const (
	directoryPath  string = "/temp/"
	imagesFolder   string = "/pictures/"
	thumbsFolder   string = "/thumbs/"
	thumbX1        string = "200"
	fileServerRoot string = "http://localhost:8100/fs/" //"http://127.0.0.1:8100/fs/"
	splitToXPages  int    = 1
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	os.Mkdir("/temp", 777)
	start := time.Now()

	dirName := functions.RandomStringGenerator(8)
	tempFileName := ""

	//parse the multipart form from index page ==> 10 << 20 specifies a maximum upload of 10MB files.
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println(err)
	}

	//for loop that checks file type (doc/docx or pdf) and then call the appropriate process function
	for _, fileTypeMap := range r.MultipartForm.File {

		for _, value := range fileTypeMap {
			if strings.Contains(value.Filename, ".pdf") {
				tempFileName = functions.ProcessPDF(w, r, directoryPath, splitToXPages, dirName, value.Filename)

			} else if strings.Contains(value.Filename, ".docx") || strings.Contains(value.Filename, ".doc") {
				tempFileName = functions.ProcessDocxDocker(w, r, directoryPath, splitToXPages, dirName)

			} else {
				fmt.Fprint(w, "Unsupported file")
			}
		}
	}
	functions.ThumbsCreation(directoryPath, dirName, thumbsFolder, thumbX1)

	functions.Marshal(directoryPath, dirName, imagesFolder, thumbsFolder, fileServerRoot, tempFileName)
	functions.OutputBodyCreation(directoryPath, dirName, imagesFolder, thumbsFolder, w, fileServerRoot, tempFileName)

	elapsed := time.Since(start)
	fmt.Println("time to execute program:", elapsed)
	defer fmt.Println("finished")

}

func CallIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(functions.HtmlHeader2))
	fmt.Fprint(w, `

<form
        enctype="multipart/form-data"
        action="/upload"
        method="post"
>
	<label>split your file</label>
    <input type="file" name="myFile" />
    <input type="submit" value="upload" />
</form>



</body>
</html>
`)

}

func setupRoutes() {

	go functions.FileServer()
	go http.HandleFunc("/upload", UploadFile)
	go http.HandleFunc("/", CallIndex)
	err := http.ListenAndServe(":8080", nil)

	fmt.Printf("error with ListenAndServer: %v \n", err)

}

func main() {

	fmt.Println("program started")
	setupRoutes()
}
