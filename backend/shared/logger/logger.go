package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/shared/env"
	"net"
)

type logLevel string

const (
	infoLevel    logLevel = "info"
	errLevel     logLevel = "error"
	warningLevel logLevel = "warning"
)

var (
	logstashServerType = env.GetString("LOGSTASH_TYPE", "tcp")
	logstashHost       = env.GetString("LOGSTASH_HOST", "ec2-54-82-72-174.compute-1.amazonaws.com")
	logstashPort       = env.GetString("LOGSTASH_PORT", "5044")
)

type Object interface {
	LogName() string
	LogProperties() map[string]interface{}
}

type Logger struct {
	Service string `json:"service,omitempty"`
}

type LogData struct {
	Service   string                            `json:"service,omitempty"`
	Level     string                            `json:"level,omitempty"`
	Error     string                            `json:"error,omitempty"`
	Event     string                            `json:"event,omitempty"`
	LogObject map[string]map[string]interface{} `json:"properties,omitempty"`
}

func NewLogger() *Logger {
	return &Logger{env.GetString("AWS_LAMBDA_FUNCTION_NAME", "unknown")}
}

func (l *Logger) Info(eventName string, objects []Object) {
	l.sendLog(infoLevel, eventName, nil, objects)
}

func (l *Logger) Warning(eventName string, err error, objects []Object) {
	l.sendLog(warningLevel, eventName, err, objects)
}

func (l *Logger) Error(eventName string, err error, objects []Object) {
	l.sendLog(errLevel, eventName, err, objects)
}

func (l *Logger) sendLog(level logLevel, eventName string, errToLog error, objects []Object) {
	connection, err := net.Dial(logstashServerType, logstashHost+":"+logstashPort)
	if err != nil {
		panic(fmt.Errorf("error connecting to Logstash server: %w", err))
	}

	logData := &LogData{
		Service:   l.Service,
		Event:     eventName,
		Level:     string(level),
		LogObject: getLogObjects(objects),
	}

	if errToLog != nil {
		logData.Error = errToLog.Error()
	}

	dataAsBytes := new(bytes.Buffer)

	err = json.NewEncoder(dataAsBytes).Encode(logData)
	if err != nil {
		panic(fmt.Errorf("logger: error encoding log data: %w", err))
	}

	_, err = connection.Write(dataAsBytes.Bytes())
	if err != nil {
		panic(fmt.Errorf("logger: error writing data to logstash: %w", err))
	}
}

func getLogObjects(objects []Object) map[string]map[string]interface{} {
	lObjects := make(map[string]map[string]interface{})

	for _, object := range objects {
		lObjects[object.LogName()] = object.LogProperties()
	}

	return lObjects
}
