package logger

import (
	"bufio"
	"bytes"
	"encoding/json"
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

	//leave this here just in case you decide to add custom log timestamps
	timestampLayout = "2006-01-02T15:04:05.999999999Z"
)

var (
	stackCleaner = regexp.MustCompile(`[^\t]*:\d+`)

	once              sync.Once
	connectionTimeout = time.Second * 3
	connDeadlineIncr  = time.Minute * 5
)

type LogAPI interface {
	Info(eventName string, fields ...models.LoggerField)
	Warning(eventName string, err error, fields ...models.LoggerField)
	Error(eventName string, err error, fields ...models.LoggerField)
	Critical(eventName string, fields ...models.LoggerField)
	LogLambdaTime(startingTime time.Time, err error, panic interface{})
	Finish() error
	MapToLoggerObject(name string, m map[string]interface{}) models.LoggerField
	SetHandler(handler string)
}

type Log struct {
	Service    string `json:"service,omitempty"`
	lambdaName string
	handler    string
	bw         *bufio.Writer
	connection net.Conn
	connTimer  *time.Timer
	wg         sync.WaitGroup
}

type LogData struct {
	Service   string                 `json:"service,omitempty"`
	Level     string                 `json:"level,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Event     string                 `json:"event,omitempty"`
	LogObject map[string]interface{} `json:"properties,omitempty"`
	Timestamp string                 `json:"@timestamp,omitempty"`
}

func NewLogger() LogAPI {
	log := &Log{
		lambdaName: env.GetString("AWS_LAMBDA_FUNCTION_NAME", "unknown"),
		bw:         bufio.NewWriter(os.Stdout),
	}

	log.establishConnection()

	return log
}

func (l *Log) SetHandler(handler string) {
	if handler != "" && l.lambdaName != "unknown" {
		l.handler = handler
	}
}

func (l *Log) Info(eventName string, fields ...models.LoggerField) {
	l.wg.Add(1)
	go l.sendLog(infoLevel, eventName, nil, fields)
}

func (l *Log) Warning(eventName string, err error, fields ...models.LoggerField) {
	l.wg.Add(1)
	go l.sendLog(warningLevel, eventName, err, fields)
}

func (l *Log) Error(eventName string, err error, fields ...models.LoggerField) {
	l.wg.Add(1)
	go l.sendLog(errLevel, eventName, err, fields)
}

func (l *Log) LogLambdaTime(startingTime time.Time, err error, panicErr interface{}) {
	duration := time.Since(startingTime).Seconds()
	durationData := map[string]interface{}{
		"f_duration": duration,
	}

	if panicErr != nil {
		panicObject := getPanicObject(panicErr)

		l.Critical("lambda_panicked", models.Any("duration_data", durationData), panicObject)
		return
	}

	if err != nil {
		l.Error("lambda_execution_finished", err, models.Any("duration_data", durationData))
	}

	l.Info("lambda_execution_finished", models.Any("duration_data", durationData))
}

func (l *Log) Critical(eventName string, fields ...models.LoggerField) {
	l.wg.Add(1)
	go l.sendLog(panicLevel, eventName, nil, fields)
}

func (l *Log) sendLog(level logLevel, eventName string, errToLog error, fields []models.LoggerField) {
	defer l.wg.Done()

	data := l.getLogDataAsBytes(level, eventName, errToLog, fields)

	err := l.write(data)
	if err != nil {
		//The lambda function shouldn't terminate because logs weren't sent. A future way of handling this
		//could be setting Cloudwatch alarms to monitor this kind of failures.
		errorLogger.Println(fmt.Errorf("logger: error writing data to logstash: %w", err))
	}
}

func (l *Log) getLogDataAsBytes(level logLevel, eventName string, errToLog error, fields []models.LoggerField) []byte {
	logData := &LogData{
		Service:   l.getService(),
		Event:     eventName,
		Level:     string(level),
		LogObject: getLogObjects(fields),
		Timestamp: time.Now().Format(timestampLayout),
	}

	if errToLog != nil {
		logData.Error = errToLog.Error()
	}

	dataBuffer := new(bytes.Buffer)

	err := json.NewEncoder(dataBuffer).Encode(logData)
	if err != nil {
		panic(fmt.Errorf("logger: error encoding log data: %w", err))
	}

	return dataBuffer.Bytes()
}

func (l *Log) getService() string {
	if l.handler != "" {
		return fmt.Sprintf("%s-%s", l.lambdaName, l.handler)
	}

	return l.lambdaName
}

func (l *Log) establishConnection() {
	logstashHost := env.GetString("LOGSTASH_HOST", "")
	logstashPort := env.GetString("LOGSTASH_PORT", "")

	once.Do(func() {
		conn, err := net.DialTimeout("tcp", logstashHost+":"+logstashPort, connectionTimeout)
		if err != nil {
			errorLogger.Println(fmt.Errorf("connection to Logstash failed: %w", err))

			return
		}

		l.bw = bufio.NewWriter(conn)
		l.connection = conn

		l.connTimer = time.NewTimer(connDeadlineIncr)
		go l.closeConnection()

		return
	})
}

// closeConnection closes the connection when the deadline is reached.
// Because the deadline is updated on each write, the connection will only be closed when no writes happen after a
// certain amount of time.
func (l *Log) closeConnection() {
	<-l.connTimer.C

	// this will be nil if a connection can never be made, in which case, there is no connection to close.
	if l.connection == nil {
		return
	}

	err := l.bw.Flush()
	if err != nil {
		errorLogger.Println(fmt.Errorf("error flushing buffer while closing to Logstash server connection: %w", err))
	}

	err = l.connection.Close()
	if err != nil {
		errorLogger.Println(fmt.Errorf("error closing connection to Logstash server: %w", err))
		return
	}

	l.connection = nil
	return
}

func (l *Log) write(data []byte) error {
	_, err := l.bw.Write(data)
	if err != nil {
		return fmt.Errorf("error writing log: %w", err)
	}

	if l.connTimer == nil {
		return nil
	}

	if !l.connTimer.Stop() {
		<-l.connTimer.C
	}
	// Reset connection deadline on each successful write. This way the connection will only be closed when there aren't
	// any writes for a certain amount of time.
	l.connTimer.Reset(connDeadlineIncr)

	return err
}

// Finish sends the remaining buffer's contents to Logstash. When the buffer is full it automatically flushes itself,
// sending the data it contains to Logstash. Therefore, Finish has to wait for all "data-sending" goroutines to be
// completed because there may or may not be a Logstash request underway.
func (l *Log) Finish() error {
	l.wg.Wait()

	err := l.bw.Flush()
	if err != nil {
		err = fmt.Errorf("error flushing logger buffer: %w", err)
		errorLogger.Println(err)
		return err
	}

	return nil
}

func (l *Log) MapToLoggerObject(name string, m map[string]interface{}) models.LoggerField {
	return models.Any(name, m)
}

// getLogObjects transforms the logger objects to a serializable representation.
func getLogObjects(objects []models.LoggerField) map[string]interface{} {
	lObjects := make(map[string]interface{})

	var value interface{}
	var err error

	for _, object := range objects {
		//If a caller pases "nil" as a logger field to one of the logging methods(Info, Error, etc), because is a var arg,
		//it will be converted to an LoggerField array with one element, which is nil.
		if object == nil {
			continue
		}

		value, err = object.GetValue()
		if err != nil {
			errorLogger.Println(fmt.Errorf("logger: error getting value for object %s: %w", object.GetKey(), err))
			continue
		}

		lObjects[object.GetKey()] = value
	}

	return lObjects
}

func getPanicObject(panic interface{}) models.LoggerField {
	clean := stackCleaner.FindAll(debug.Stack(), -1)

	return models.Any("panic", map[string]interface{}{
		"s_message": fmt.Sprintf("%v", panic),
		"s_trace":   string(bytes.Join(clean, []byte("\n\n"))),
	})
}
