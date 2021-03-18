package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"io/ioutil"
	"strconv"
	"time"
)

var (
	Ak       = "1"                     // account assigned by the file service
	Sk       = "2"                     // file service assigned key
	endPoint = "http://127.0.0.1:8333" //The address of the weed S3 service on the above image
	Region   = "test.com"              //Scope of application
	svc      *s3.S3
)

func init() {
	cres := credentials.NewStaticCredentials(Ak, Sk, "")
	cfg := aws.NewConfig().WithRegion(Region).WithEndpoint(endPoint).WithCredentials(cres).WithS3ForcePathStyle(true)
	sess, err := session.NewSession(cfg)
	if err != nil {
		fmt.Println(err)
	}
	svc = s3.New(sess)
}

func main() {

	// Create a bucket
	bucketName := "weed-test-buck" //The name of the bucket is also the unique identifier for accessing the data below this bucket.
	createBucket(bucketName)
	// Upload image data to the weed file service

	_ = api.SplitFile("/home/jelsik/Downloads/temp-folder/gobook.pdf", "/home/jelsik/Downloads/temp-folder/", 1, nil)

	files, _ := ioutil.ReadDir("/home/jelsik/Downloads/temp-folder/")
	i := 1
	for _, file := range files {
		objectID := "pdfSplit/myFile" + strconv.Itoa(i) + ".pdf"

		dataImage, err := ioutil.ReadFile("/home/jelsik/Downloads/temp-folder/" + file.Name())
		if err != nil {
			fmt.Println(err.Error())
		}
		contentType := "application/pdf"
		putS3Object(dataImage, bucketName, contentType, objectID)
		//Get the file
		req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectID),
		})
		urlStr, err := req.Presign(7 * 24 * time.Hour)
		fmt.Println(urlStr)
		i++
	}
}

func createBucket(bucketName string) {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}
	result, err := svc.CreateBucket(input)
	fmt.Println(result)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println(s3.ErrCodeBucketAlreadyExists, aerr.Error())
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}
}

func putS3Object(dataImage []byte, bucketName, contentType, objectID string) {

	inputObject := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectID),
		ContentType: aws.String(contentType),
		Body:        bytes.NewReader(dataImage),
	}
	_, err := svc.PutObject(inputObject) // first throw-out was "resp"
	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Println(resp)
}

func getS3Object(bucketName, objectID string) []byte {

	inputObject := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectID),
	}
	out, err := svc.GetObject(inputObject)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	res, err := ioutil.ReadAll(out.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return res
}

func deleteS3Object(bucketName, objectID string) {
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectID),
	}

	resp, err := svc.DeleteObject(params)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(resp)
}
