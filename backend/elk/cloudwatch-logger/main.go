package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io"
)

func handler(event events.CloudwatchLogsEvent) error {
	logsData, err := event.AWSLogs.Parse()
	if err != nil {
		return err
	}

	decodedString, err := base64.StdEncoding.DecodeString(event.AWSLogs.Data)
	if err != nil {
		return err
	}

	gzipReader, err := gzip.NewReader(bytes.NewReader(decodedString))
	if err != nil {
		return err
	}

	data, err := io.ReadAll(gzipReader)
	if err != nil {
		return err
	}

	fmt.Println("Cloudwatch event: ", string(data))

	return nil
}

func main() {
	lambda.Start(handler)
}
