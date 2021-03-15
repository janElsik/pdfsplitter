package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"net/http"
	"os"
	"pdf/99_old_code/10_v0.07_MINIO_S3/functions"
	"strings"
	"time"
)

const (
	directoryPath string = "/temp/"
	thumbX1       string = "200"
	splitToXPages int    = 1
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	S3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials("minioadmin", "minioadmin", ""),
		Endpoint:         aws.String("http://localhost:9000"),
		Region:           aws.String("eu-central-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	var dirName = functions.RandomStringGenerator(8)
	var tempFileName string = ""
	fmt.Print(tempFileName)

	//parse the multipart form from index page ==> 10 << 20 specifies a maximum upload of 10MB files.
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println(err)
	}

	//for loop that checks file type (doc/docx or pdf) and then call the appropriate process function
	for _, fileTypeMap := range r.MultipartForm.File {

		for _, value := range fileTypeMap {
			if strings.Contains(value.Filename, ".pdf") {
				tempFileName = functions.ProcessPDF(w, r, directoryPath, splitToXPages, dirName)

			} else if strings.Contains(value.Filename, ".docx") || strings.Contains(value.Filename, ".doc") {
				tempFileName = functions.ProcessDocx(w, r, directoryPath, splitToXPages, dirName)

			} else {
				fmt.Fprint(w, "Unsupported file")
			}
		}
	}
	functions.ThumbsCreation(directoryPath, dirName, thumbX1)
	functions.FolderUpload(directoryPath, dirName, S3Config)
	functions.OutputBodyCreation(directoryPath, dirName, w, S3Config)

	err = os.RemoveAll(directoryPath + dirName + "/")
	fmt.Println(err)
	elapsed := time.Since(start)
	fmt.Println("time to execute program:", elapsed)
	defer fmt.Println("finished")

	fmt.Fprint(w, `</body>`)

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
	go http.HandleFunc("/upload", UploadFile)
	go http.HandleFunc("/", CallIndex)
	err := http.ListenAndServe(":8080", nil)

	fmt.Printf("error with ListenAndServer: %v \n", err)

}

func main() {

	fmt.Println("program started")
	setupRoutes()
}
