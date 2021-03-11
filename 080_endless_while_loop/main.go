package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {

	// endless loop that checks, if said folder is empty. If not, it deletes files there (one by one, not all)
	for {
		files, _ := ioutil.ReadDir("/temp/temp-folder/subfolder")

		if len(files) != 0 {
			for _, file := range files {
				fmt.Println(file.Name())
				err := os.Remove("/temp/temp-folder/subfolder/" + file.Name())
				if err != nil {
					log.Println("error with deleting: ", err)
				}
				log.Println(file.Name(), "deleted!")
			}
		}
		time.Sleep(2 * time.Second)
	}
}
