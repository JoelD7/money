package logger

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"net"
	"os"
	"regexp"
	"runtime/debug"
	"sync"
	"time"
)

type logLevel string

const (
	infoLevel    logLevel = "info"
	errLevel     logLevel = "error"
	warningLevel logLevel = "warning"
	panicLevel   logLevel = "panic"

	connectionTimeout = time.Second * 3
	//leave this here just in case you decide to add custom log timestamps
	timestampLayout = "2006-01-02T15:04:05.999999999Z"
)

var (
	logstashServerType = env.GetString("LOGSTASH_TYPE", "tcp")
	logstashHost       = env.GetString("LOGSTASH_HOST", "ec2-54-158-20-203.compute-1.amazonaws.com")
	logstashPort       = env.GetString("LOGSTASH_PORT", "5044")

	stackCleaner = regexp.MustCompile(`[^\t]*:\d+`)

	connection net.Conn
)

type LogAPI interface {
	Info(eventName string, objects []models.LoggerObject)
	Warning(eventName string, err error, objects []models.LoggerObject)
	Error(eventName string, err error, objects []models.LoggerObject)
	Critical(eventName string, objects []models.LoggerObject)
	LogLambdaTime(startingTime time.Time, err error, panic interface{})
	Close() error
	MapToLoggerObject(name string, m map[string]interface{}) models.LoggerObject
	SetHandler(handler string)
}

type Log struct {
	Service   string `json:"service,omitempty"`
	useBackup bool
	bw        *bufio.Writer
	wg        sync.WaitGroup
}

type LogData struct {
	Service   string                            `json:"service,omitempty"`
	Level     string                            `json:"level,omitempty"`
	Error     string                            `json:"error,omitempty"`
	Event     string                            `json:"event,omitempty"`
	LogObject map[string]map[string]interface{} `json:"properties,omitempty"`
}

func NewLogger() LogAPI {
	log := &Log{
		Service: env.GetString("AWS_LAMBDA_FUNCTION_NAME", "unknown"),
		bw:      new(bufio.Writer),
	}

	log.establishConnection()

	return log
}

func NewLoggerWithHandler(handler string) LogAPI {
	log := &Log{
		Service: env.GetString("AWS_LAMBDA_FUNCTION_NAME", "unknown"),
		bw:      new(bufio.Writer),
	}

	if handler != "" && log.Service != "unknown" {
		log.Service += "-" + handler
	}

	log.establishConnection()

	return log
}

func (l *Log) SetHandler(handler string) {
	if handler != "" && l.Service != "unknown" {
		l.Service += "-" + handler
	}
}

func (l *Log) Info(eventName string, objects []models.LoggerObject) {
	l.wg.Add(1)
	go l.sendLog(infoLevel, eventName, nil, objects)
}

func (l *Log) Warning(eventName string, err error, objects []models.LoggerObject) {
	l.wg.Add(1)
	go l.sendLog(warningLevel, eventName, err, objects)
}

func (l *Log) Error(eventName string, err error, objects []models.LoggerObject) {
	l.wg.Add(1)
	go l.sendLog(errLevel, eventName, err, objects)
}

func (l *Log) LogLambdaTime(startingTime time.Time, err error, panicErr interface{}) {
	duration := time.Since(startingTime).Seconds()
	durationData := l.MapToLoggerObject("duration_data", map[string]interface{}{
		"f_duration": duration,
	})

	if panicErr != nil {
		panicObject := getPanicObject(panicErr)

		l.Critical("lambda_panicked", []models.LoggerObject{durationData, panicObject})
		return
	}

	if err != nil {
		l.Error("lambda_execution_finished", err, []models.LoggerObject{durationData})
	}

	l.Info("lambda_execution_finished", []models.LoggerObject{durationData})
}

func (l *Log) Critical(eventName string, objects []models.LoggerObject) {
	l.wg.Add(1)
	go l.sendLog(panicLevel, eventName, nil, objects)
}

func (l *Log) sendLog(level logLevel, eventName string, errToLog error, objects []models.LoggerObject) {
	defer l.wg.Done()

	data := l.getLogDataAsBytes(level, eventName, errToLog, objects)

	err := l.write(data)
	if err != nil {
		//The lambda function shouldn't terminate because logs weren't sent. A future way of handling this
		//could be setting Cloudwatch alarms to monitor this kind of failures.
		errorLogger.Println(fmt.Errorf("logger: error writing data to logstash: %w", err))
	}
}

func (l *Log) getLogDataAsBytes(level logLevel, eventName string, errToLog error, objects []models.LoggerObject) []byte {
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

	err := json.NewEncoder(dataAsBytes).Encode(logData)
	if err != nil {
		panic(fmt.Errorf("logger: error encoding log data: %w", err))
	}

	return dataAsBytes.Bytes()
}

func (l *Log) establishConnection() {
	if connection != nil {
		return
	}

	conn, err := net.DialTimeout("tcp", logstashHost+":"+logstashPort, connectionTimeout)
	if err != nil {
		errorLogger.Println(fmt.Errorf("connection to Logstash failed: %w", err))
		//Write to std out if connection to Logstash fails
		l.bw = bufio.NewWriter(os.Stdout)

		return
	}

	connection = conn
	l.bw = bufio.NewWriter(connection)

	return
}

func (l *Log) write(data []byte) error {
	if connection != nil {
		err := connection.SetDeadline(time.Now().Add(connectionTimeout))
		if err != nil {
			return fmt.Errorf("error setting deadline: %w", err)
		}
	}

	_, err := l.bw.Write(data)
	if err != nil {
		return fmt.Errorf("error writing log: %w", err)
	}

	if errors.Is(err, os.ErrDeadlineExceeded) {
		err = connection.SetDeadline(time.Now().Add(connectionTimeout))
		if err != nil {
			return fmt.Errorf("error setting deadline: %w", err)
		}
	}

	return err
}

// Close closes the connection to the Logstash server
func (l *Log) Close() error {
	l.wg.Wait()

	err := l.bw.Flush()
	if err != nil {
		return fmt.Errorf("error flushing logger buffer: %w", err)
	}

	// this will be nil if a connection can never be made, in which case, there is no connection to close.
	if connection == nil {
		return nil
	}

	err = connection.Close()
	if err != nil {
		return fmt.Errorf("error closing connection to Logstash server: %w", err)
	}

	connection = nil

	return nil
}

func (l *Log) MapToLoggerObject(name string, m map[string]interface{}) models.LoggerObject {
	return &ObjectWrapper{
		name:       name,
		properties: m,
	}
}

// getLogObjects transforms the logger objects to a serializable representation.
func getLogObjects(objects []models.LoggerObject) map[string]map[string]interface{} {
	lObjects := make(map[string]map[string]interface{})

	for _, object := range objects {
		lObjects[object.LogName()] = object.LogProperties()
	}

	return lObjects
}

func getPanicObject(panic interface{}) models.LoggerObject {
	clean := stackCleaner.FindAll(debug.Stack(), -1)

	return &ObjectWrapper{
		name: "panic",
		properties: map[string]interface{}{
			"s_message": fmt.Sprintf("%v", panic),
			"s_trace":   string(bytes.Join(clean, []byte("\n\n"))),
		},
	}
}
