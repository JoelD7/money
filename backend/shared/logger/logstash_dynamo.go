package logger

import (
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/aws/smithy-go/logging"
	"strings"
)

type LogstashDynamo struct {
}

func NewLogstashDynamo() *LogstashDynamo {
	return &LogstashDynamo{}
}

func (l *LogstashDynamo) Logf(classification logging.Classification, format string, v ...interface{}) {
	eventMsg := "dynamo_query_request"
	if strings.Contains(format, "Response") {
		eventMsg = "dynamo_query_response"
	}

	logfields := make([]models.LoggerField, 0)

	for i, value := range v {
		logfields = append(logfields, models.Any(fmt.Sprintf("value_%d", i), value))
	}

	if classification == logging.Debug {
		Info(eventMsg, logfields...)
		return
	}

	Warning(eventMsg, nil, logfields...)
}
