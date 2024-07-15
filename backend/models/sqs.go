package models

import (
	"github.com/aws/aws-lambda-go/events"
)

type MissingExpensePeriodMessage struct {
	Period   string `json:"period"`
	Username string `json:"username"`
}

// SQSMessage represents a message from an SQS event. This custom type exists to be able to implement the LoggerObject interface.
type SQSMessage struct {
	events.SQSMessage
}

func (s SQSMessage) LogName() string {
	return "sqs_message"
}

func (s SQSMessage) LogProperties() map[string]interface{} {
	return map[string]interface{}{
		"s_message_id":                s.MessageId,
		"s_receipt_handle":            s.ReceiptHandle,
		"s_body":                      s.Body,
		"s_md5_of_body":               s.Md5OfBody,
		"s_md5_of_message_attributes": s.Md5OfMessageAttributes,
		"o_attributes":                s.Attributes,
		"s_event_source_arn":          s.EventSourceARN,
		"s_event_source":              s.EventSource,
		"s_aws_region":                s.AWSRegion,
	}
}
