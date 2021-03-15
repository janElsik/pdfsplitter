package functions

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"os"
)

func FolderUpload(directoryPath string, dirName string, S3Config *aws.Config) {
	bucket := aws.String(dirName)

	newSession := session.New(S3Config)

	s3Client := s3.New(newSession)

	cparams := &s3.CreateBucketInput{
		Bucket: bucket, //required
	}

	// Create a new bucket using the CreateBucket call
	_, err := s3Client.CreateBucket(cparams)
	if err != nil {
		// Message from an error.
		fmt.Println(err.Error())
	}
	files, _ := ioutil.ReadDir(directoryPath + dirName + "/")

	for _, file := range files {
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

	}

}
