package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/shared/env"
	"log"
	"net"
	"os"
	"regexp"
	"runtime/debug"
	"time"
)

type logLevel string

const (
	infoLevel    logLevel = "info"
	errLevel     logLevel = "error"
	warningLevel logLevel = "warning"
	panicLevel   logLevel = "panic"

	retries           = 3
	backoffFactor     = 2
	connectionTimeout = time.Second * 5
)

var (
	logstashServerType = env.GetString("LOGSTASH_TYPE", "tcp")
	logstashHost       = env.GetString("LOGSTASH_HOST", "ec2-18-215-157-194.compute-1.amazonaws.com")
	logstashPort       = env.GetString("LOGSTASH_PORT", "5044")

	LogClient LogAPI

	stackCleaner = regexp.MustCompile(`[^\t]*:\d+`)
	errorLogger  = log.New(os.Stderr, "ERROR ", log.Llongfile)
)

type Object interface {
	LogName() string
	LogProperties() map[string]interface{}
}

type LogAPI interface {
	Info(eventName string, objects []Object)
	Warning(eventName string, err error, objects []Object)
	Error(eventName string, err error, objects []Object)
	Critical(eventName string, objects []Object)
	LogLambdaTime(startingTime time.Time, err error, panic interface{})
	Close() error
}

type Log struct {
	Service    string `json:"service,omitempty"`
	Handler    string `json:"handler,omitempty"`
	connection net.Conn
	worker     *asyncWorker
}

type asyncWorker struct {
	doneCh chan bool
	buffer chan [][]byte
}

type LogData struct {
	Service   string                            `json:"service,omitempty"`
	Handler   string                            `json:"handler,omitempty"`
	Level     string                            `json:"level,omitempty"`
	Error     string                            `json:"error,omitempty"`
	Event     string                            `json:"event,omitempty"`
	LogObject map[string]map[string]interface{} `json:"properties,omitempty"`
}

func NewLogger() LogAPI {
	if LogClient != nil {
		return LogClient
	}

	LogClient = &Log{
		Service: env.GetString("AWS_LAMBDA_FUNCTION_NAME", "unknown"),
		worker: &asyncWorker{
			doneCh: make(chan bool),
			buffer: make(chan [][]byte),
		},
	}

	return LogClient
}

func NewLoggerWithHandler(handler string) LogAPI {
	if LogClient != nil {
		return LogClient
	}

	LogClient = &Log{
		Service: env.GetString("AWS_LAMBDA_FUNCTION_NAME", "unknown"),
		Handler: handler,
		worker: &asyncWorker{
			doneCh: make(chan bool),
			buffer: make(chan [][]byte),
		},
	}

	return LogClient
}

func (l *Log) Info(eventName string, objects []Object) {
	l.sendLog(infoLevel, eventName, nil, objects)
}

func (l *Log) Warning(eventName string, err error, objects []Object) {
	l.sendLog(warningLevel, eventName, err, objects)
}

func (l *Log) Error(eventName string, err error, objects []Object) {
	l.sendLog(errLevel, eventName, err, objects)
}

func (l *Log) LogLambdaTime(startingTime time.Time, err error, panic interface{}) {
	duration := time.Since(startingTime).Seconds()
	durationData := MapToLoggerObject("duration_data", map[string]interface{}{
		"f_duration": duration,
	})

	if panic != nil {
		panicObject := getPanicObject(panic)

		l.Critical("lambda_panicked", []Object{durationData, panicObject})
		return
	}

	if err != nil {
		l.Error("lambda_execution_finished", err, []Object{durationData})
	}

	l.Info("lambda_execution_finished", []Object{durationData})
}

func (l *Log) Critical(eventName string, objects []Object) {
	l.sendLog(panicLevel, eventName, nil, objects)
}

func (l *Log) sendLog(level logLevel, eventName string, errToLog error, objects []Object) {
	err := l.connect()
	if err != nil {
		errorLogger.Println(fmt.Errorf("error connecting to Logstash server: %w", err))

		return
	}

	logData := &LogData{
		Service:   l.Service,
		Handler:   l.Handler,
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

	err = l.write(dataAsBytes.Bytes())
	if err != nil {
		panic(fmt.Errorf("logger: error writing data to logstash: %w", err))
	}
}

func (l *Log) connect() error {
	if l.connection != nil {
		return nil
	}

	connection, err := net.DialTimeout("tcp", logstashHost+":"+logstashPort, connectionTimeout)

	backoff := time.Second * 2

	for i := 0; i < retries && err != nil; i++ {
		time.Sleep(backoff)

		connection, err = net.DialTimeout("tcp", logstashHost+":"+logstashPort, time.Second*1)
		backoff *= backoffFactor
	}

	l.connection = connection

	return err
}

func (l *Log) write(data []byte) error {
	_, err := l.connection.Write(data)
	backoff := time.Second * 2

	for i := 0; i < retries && err != nil; i++ {
		time.Sleep(backoff)

		_, err = l.connection.Write(data)
		backoff *= backoffFactor
	}

	return err
}

// Close closes the connection to the Logstash server
func (l *Log) Close() error {
	err := l.connection.Close()
	if err != nil {
		return fmt.Errorf("error closing connection to Logstash server: %w", err)
	}

	return nil
}

func MapToLoggerObject(name string, m map[string]interface{}) Object {
	return &ObjectWrapper{
		name:       name,
		properties: m,
	}
}

func getLogObjects(objects []Object) map[string]map[string]interface{} {
	lObjects := make(map[string]map[string]interface{})

	for _, object := range objects {
		lObjects[object.LogName()] = object.LogProperties()
	}

	return lObjects
}

func getPanicObject(panic interface{}) Object {
	clean := stackCleaner.FindAll(debug.Stack(), -1)

	return &ObjectWrapper{
		name: "panic",
		properties: map[string]interface{}{
			"s_message": fmt.Sprintf("%v", panic),
			"s_trace":   string(bytes.Join(clean, []byte("\n\n"))),
		},
	}
}
