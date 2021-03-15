package functions

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func OutputBodyCreation(directoryPath string, dirName string, w http.ResponseWriter, S3Config *aws.Config) {
	bucket := aws.String(dirName)

	newSession := session.New(S3Config)

	s3Client := s3.New(newSession)

	files, _ := ioutil.ReadDir(directoryPath + dirName + "/")

	for _, file := range files {
		if strings.Contains(file.Name(), ".png") {
			continue
		}
		key := aws.String(file.Name())
		f, err := os.Open(directoryPath + dirName + "/" + file.Name())

		if err != nil {
			panic(err)
		}

		_, err = s3Client.PutObject(&s3.PutObjectInput{
			Body:   f,
			Bucket: bucket,
			Key:    key,
			ACL:    aws.String(s3.BucketCannedACLPublicRead),
		})

		req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: bucket,
			Key:    key,
		})
		urlStr, err := req.Presign(7 * 24 * time.Hour)

		newKeyNormalString := strings.ReplaceAll(file.Name(), ".pdf", ".png")
		newKeyAWSString := aws.String("THUMBS" + newKeyNormalString)

		req, _ = s3Client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: bucket,
			Key:    newKeyAWSString,
		})
		imgUrlStr, err := req.Presign(7 * 24 * time.Hour)

		if err != nil {
			panic(err)
		}

		var picLink string = `
		<a href=` + urlStr + `><img src=` + imgUrlStr + `></a>
	`

		fmt.Fprint(w, picLink)

	}

}
