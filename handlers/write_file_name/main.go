package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

func main() {
	lambda.Start(handler)
}

func handler(sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		ext, err := getExtFromMessage(message)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("extension of the file is %s", ext)
	}

	return nil
}

func getExtFromMessage(e events.SQSMessage) (string, error) {
	var s3event events.S3Event
	if err := json.Unmarshal([]byte(e.Body), &s3event); err != nil {
		return "", errors.Wrapf(err, "failed to unmarshal: %s", e.Body)
	}
	str, err := url.QueryUnescape(s3event.Records[0].S3.Object.Key)
	if err != nil {
		return "", errors.Wrapf(err, "failed to unescape file name: %s", s3event.Records[0].S3.Object.Key)
	}
	return filepath.Ext(str), nil
}
