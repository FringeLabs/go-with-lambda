package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"log"
)

// Motivation stored in the payload
type Motivation struct {
	DoYouLove bool   `json:"doyoulove"`
	What      string `json:"what"`
	Why       string `json:"why"`
	When      int    `json:"when"`
	Where     struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"where"`
}

// Handler of Lambda function
func Handler(request events.S3Event) error {
	log.Println("Received event", request)

	// getting S3 object references from event
	var bucket, key string = request.Records[0].S3.Bucket.Name, request.Records[0].S3.Object.Key
	log.Println("Bucket: ", bucket, "\nKey: ", key)

	// initialize aws session and creating downloader
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// reading S3 content
	downloader := s3manager.NewDownloader(sess)
	payload, err := downloader.S3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key)})

	if err != nil {
		log.Fatal(err)
		return err
	}

	// reading data buffer
	log.Println("Payload:", payload)
	raw, err := ioutil.ReadAll(payload.Body)

	if err != nil {
		log.Fatal(err)
		return err
	}

	// unmarshalling
	parsed := string(raw[:])
	var m Motivation
	json.Unmarshal(raw, &m)
	log.Println("Content:", parsed)

	return nil
}

func main() {
	lambda.Start(Handler)
}
