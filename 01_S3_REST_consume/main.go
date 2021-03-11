package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"io/ioutil"
	"os/exec"
	"strings"
)

var (
	s3Client *s3.S3
)

const (
	BUCKET_NAME = "testbucket"
)

func init() {
	S3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials("minioadmin", "minioadmin", ""),
		Endpoint:         aws.String("http://localhost:9000"),
		Region:           aws.String("eu-central-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession := session.New(S3Config)

	s3Client = s3.New(newSession)

}

func main() {

	var objectList []string

	i := 1

	for _, object := range listObjects(BUCKET_NAME).Contents {
		objectList = append(objectList, *object.Key)
		getObject(*object.Key, BUCKET_NAME, "/temp/temp-folder/"+*object.Key)
		if strings.Contains(*object.Key, ".pdf") {
			err := api.SplitFile("/temp/temp-folder/"+*object.Key, "/temp/temp-folder/split", 1, nil)
			if err != nil {
				panic(err)
			}
		}
		if strings.Contains(*object.Key, ".xls") || strings.Contains(*object.Key, ".xlsx") {
			cmd := exec.Command("unoconv", "-f", "pdf", "-d", "s", "/temp/temp-folder/"+*object.Key)
			if err := cmd.Run(); err != nil {
				fmt.Println("here", err)
				panic(err)
			}
		}

		i++

	}

}

func getObject(fileNameToDownload string, bucketName string, fileNameToWrite string) {
	fmt.Println("Downloading:", fileNameToDownload)

	resp, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileNameToDownload),
	})

	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	err = ioutil.WriteFile(fileNameToWrite, body, 0644)
	if err != nil {
		panic(err)
	}

}

func listObjects(bucketName string) (resp *s3.ListObjectsV2Output) {
	resp, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		panic(err)
	}
	return resp
}
