package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net"
	"strings"
)

var (
	logstashServerType = env.GetString("LOGSTASH_TYPE", "tcp")
	logstashHost       = env.GetString("LOGSTASH_HOST", "ec2-54-209-181-211.compute-1.amazonaws.com")
	logstashPort       = env.GetString("LOGSTASH_PORT", "5044")
)

type logstashEvent struct {
	Message tokenEvent `json:"message"`
}

type tokenEvent struct {
	Token string `json:"token"`
}

func handler(event events.CloudwatchLogsEvent) error {
	fmt.Println("in handler")

	rawData := event.AWSLogs.Data
	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	zr, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer zr.Close()

	var logsData events.CloudwatchLogsData
	dec := json.NewDecoder(zr)

	err = dec.Decode(&logsData)
	if err != nil {
		fmt.Println("error parsing logs")
		return err
	}

	//establish connection
	connection, err := net.Dial(logstashServerType, logstashHost+":"+logstashPort)
	if err != nil {
		fmt.Println("error establish connection")
		panic(err)
	}

	cleanedMessage := strings.TrimSpace(strings.Split(logsData.LogEvents[0].Message, "money_app_log:  ")[1])

	tk := &tokenEvent{}
	err = json.Unmarshal([]byte(cleanedMessage), tk)
	if err != nil {
		fmt.Println(err)
	}

	lEvent := &logstashEvent{
		Message: *tk,
	}

	str, err := utils.GetJsonString(lEvent)
	if err != nil {
		fmt.Println(err)
		return err
	}

	///send some data
	_, err = connection.Write([]byte(str + "\n"))
	if err != nil {
		fmt.Println("error writing logs")
	}
	fmt.Println("Wrote message:", str)

	fmt.Println("done")
	defer connection.Close()

	return nil
}

func main() {
	lambda.Start(handler)
}
