package main

import (
	"encoding/json"
	"fmt"
	"github.com/ginuerzh/weedo"
	"io/ioutil"
	"net/http"
	"os"
	"pdfsplitter_cez_preprod/01_entry_point/functions"
	"pdfsplitter_cez_preprod/01_entry_point/helpers"
	"strings"
	"sync"
	"time"
)

const (
	WeedsAddressString        = "10.4.237.28:9333"
	NatsAddressString         = "10.4.220.151:4222"
	WeedsVolumeVisibleAddress = "http://10.4.56.20:8080/"
)

//ALL LINKS THAT ARE TO BE VISIBLE ARE CHANGED FROM INNER CLUSTER ADDRESS TO OUTER CLUSTER ADDRESS (CONST WeedsVolumeVisibleAddress) SO THAT THEY ARE CLICKABLE
type JSON struct {
	Number int    `json:"linknumber"`
	Href   string `json:"href"`
	ImgSrc string `json:"imgsrc"`
}

func setupRoutes() {
	// functions that handle index and localhost:[port]/process
	go http.HandleFunc("/", CallIndex)
	go http.HandleFunc("/process", Organizer)
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Printf("error with ListenAndServe: %v \n", err)
	}

}

func CallIndex(w http.ResponseWriter, r *http.Request) {

	//write out all the logic of html index, instead of body
	_, err := w.Write([]byte(helpers.HtmlHeader2))

	if err != nil {
		fmt.Println("Response writer:", err)
	}
	//write out html body
	_, err = fmt.Fprint(w, `

<form
        enctype="multipart/form-data"
        action="/process"
        method="post"
>
	<label>split your file</label>
    <input type="file" name="myfiles" multiple=multiple/>
    <input type="submit" value="upload" />
</form>



</body>
</html>
`)

	if err != nil {
		fmt.Println("Fprint:", err)
	}

}

func main() {
	fmt.Println("program started")
	setupRoutes()
}

func Organizer(w http.ResponseWriter, r *http.Request) {

	// connection to filesystem
	weedoClient := weedo.NewClient(WeedsAddressString)
	// start to track time since start of program
	start := time.Now()
	_ = os.Mkdir("/temp", 0777)

	// array with filepaths
	var inputFileSlice []string

	// randomly generated string used to rename the files to unique names
	tempFileName := helpers.RandomStringGenerator(12)

	// randomly generated string used to create folder with unique name
	tempFolderName := "/temp/" + helpers.RandomStringGenerator(12)
	fmt.Println("tempfoldername:", tempFolderName)
	err := os.Mkdir(tempFolderName, 0777)
	if err != nil {
		fmt.Println("error with making dir")
	}

	tempFolderName = tempFolderName + "/"

	// 32MB is the default used by FormFile
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Println(err)
	}

	// iterate over data from html form and save them in form of files
	for _, fileTypeMap := range r.MultipartForm.File {
		for _, value := range fileTypeMap {

			f, err := value.Open()
			if err != nil {
				fmt.Println("opening value:", err)
			}

			fmt.Println(value.Filename)
			fmt.Println(tempFolderName + value.Filename)
			tempFile, err := ioutil.TempFile(tempFolderName, value.Filename)

			if err != nil {
				fmt.Println("temp file initializing:", err)
			}
			fileBytes, err := ioutil.ReadAll(f)
			if err != nil {
				fmt.Println("reading filebytes:", err)
			}

			_, err = tempFile.Write(fileBytes)

			if err != nil {
				fmt.Println("writing temp file:", err)
			}

			_ = os.Rename(tempFile.Name(), tempFolderName+value.Filename)

			tempFile.Close()
			_ = f.Close()

		}
	}

	// read in all file names in temp dir and put them into an array
	dirSlice, _ := os.ReadDir(tempFolderName)

	for _, file := range dirSlice {
		inputFileSlice = append(inputFileSlice, tempFolderName+file.Name())
	}

	// this block makes sure that conversion (functions.Convert) is completed before continuing with the
	// execution of the program (pointer to wg variable)
	var wg sync.WaitGroup
	wg.Add(1)
	command := make(chan string)

	// possible through go routine, but potentionally very costly regarding memory
	var linkSlice []string

	// conversion call on input files, returns Array with links converted documents
	linkSlice = functions.Convert(inputFileSlice, tempFolderName, tempFileName, command, &wg, WeedsAddressString, NatsAddressString)
	err = os.RemoveAll(tempFolderName)
	if err != nil {
		fmt.Println(err)
	}

	err = os.Mkdir(tempFolderName, 0777)
	if err != nil {
		fmt.Println(err)
	}

	// waits for completion of conversion call
	wg.Wait()

	// prints name of temp folder
	fmt.Println("folder:", tempFolderName)

	// prints the links to converted documents
	for _, link := range linkSlice {

		fmt.Println(link)

	}

	// call to merge converted files into one, returns link to merged file
	mergedFileLink := functions.Merge(tempFolderName, linkSlice, WeedsAddressString)

	// prints the link to merged file
	fmt.Println("link to merged file:", mergedFileLink)

	// this block makes sure that split (functions.Split) is completed before continuing with the
	// execution of the program (pointer to wg2 variable)
	var wg2 sync.WaitGroup
	wg2.Add(1)

	// call to split the merged file into single pages, returns links to split pdfs and to thumbnails of
	// the split pdfs
	thumbSlice, splitLinkSlice := functions.Split(tempFolderName, mergedFileLink, &wg2, WeedsAddressString, NatsAddressString)
	wg2.Wait()

	err = os.Mkdir(tempFolderName, 0777)
	if err != nil {
		fmt.Println(err)
	}

	// create new struct for JSON
	jsonOutput := []JSON{}

	// iterate over pdf and thumbnail links and add to JSON struct
	for i, link := range splitLinkSlice {

		tempLinkSlice1 := strings.SplitAfterN(link, "/", 4)
		newLink := WeedsVolumeVisibleAddress + tempLinkSlice1[len(tempLinkSlice1)-1]

		//strings.ReplaceAll(link, "10.4.129.95", "10.4.56.20")

		tempLinkSlice2 := strings.SplitAfterN(thumbSlice[i], "/", 4)
		newThumbLink := WeedsVolumeVisibleAddress + tempLinkSlice2[len(tempLinkSlice2)-1]
		//strings.ReplaceAll(thumbLinks[count], "10.4.129.95", "10.4.56.20")

		//fmt.Println(thumbSlice[i], link)
		fmt.Println("<a href=" + newLink + "><img src=" + newThumbLink + "></a>")
		jsonOutput = append(jsonOutput, JSON{
			Number: i,
			Href:   link,
			ImgSrc: thumbSlice[i],
		})

	}

	// create byteArray of the JSON struct
	byteArray, err := json.Marshal(jsonOutput)
	if err != nil {
		fmt.Println("Marshaling:", err)
	}

	// write the JSON struct to tempfolder
	err = os.WriteFile(tempFolderName+"jsongo.json", byteArray, 0644)
	if err != nil {
		fmt.Println("Writing marshaled file:", err)
	}
	file, err := os.Open(tempFolderName + "jsongo.json")

	if err != nil {
		fmt.Println("Opening file:", err)
	}

	// upload JSON file, get url and print the url
	fid, _, err := weedoClient.AssignUpload("jsongo.json", "application/json", file)
	if err != nil {
		fmt.Println("Opening file:", err)
	}
	purl, _, err := weedoClient.GetUrl(fid)
	if err != nil {
		fmt.Println("Getting url:", err)
	}

	err = os.RemoveAll(tempFolderName)
	if err != nil {
		fmt.Println("Removing file:", err)
	}

	tempStringSlice := strings.SplitAfterN(purl, "/", 4)
	newPurl := WeedsVolumeVisibleAddress + tempStringSlice[len(tempStringSlice)-1]

	fmt.Println("link to JSON:", newPurl)

	tempStringSlice = strings.SplitAfterN(mergedFileLink, "/", 4)
	newMergedFileLink := WeedsVolumeVisibleAddress + tempStringSlice[len(tempStringSlice)-1]

	fmt.Println("link to merged file:", newMergedFileLink)

	functions.OutputBodyCreation(w, splitLinkSlice, thumbSlice, WeedsVolumeVisibleAddress)
	elapsed := time.Since(start)
	fmt.Println("process took:", elapsed)

}
