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
	logstashServerType = env.GetString("LOGSTASH_TYPE", "tcp")
	logstashHost       = env.GetString("LOGSTASH_HOST", "ec2-54-196-205-47.compute-1.amazonaws.com")
	logstashPort       = env.GetString("LOGSTASH_PORT", "5044")

	stackCleaner = regexp.MustCompile(`[^\t]*:\d+`)

	once              sync.Once
	connectionTimeout = time.Second * 3
	connDeadlineIncr  = time.Minute * 5
)

type LogAPI interface {
	Info(eventName string, objects []models.LoggerObject)
	Warning(eventName string, err error, objects []models.LoggerObject)
	Error(eventName string, err error, objects []models.LoggerObject)
	Critical(eventName string, objects []models.LoggerObject)
	LogLambdaTime(startingTime time.Time, err error, panic interface{})
	Finish() error
	MapToLoggerObject(name string, m map[string]interface{}) models.LoggerObject
	SetHandler(handler string)
}

type Log struct {
	Service    string `json:"service,omitempty"`
	bw         *bufio.Writer
	connection net.Conn
	connTimer  *time.Timer
	wg         sync.WaitGroup
}

type LogData struct {
	Service   string                            `json:"service,omitempty"`
	Level     string                            `json:"level,omitempty"`
	Error     string                            `json:"error,omitempty"`
	Event     string                            `json:"event,omitempty"`
	LogObject map[string]map[string]interface{} `json:"properties,omitempty"`
	Timestamp string                            `json:"@timestamp,omitempty"`
}

func NewLogger() LogAPI {
	log := &Log{
		Service: env.GetString("AWS_LAMBDA_FUNCTION_NAME", "unknown"),
		bw:      bufio.NewWriter(os.Stdout),
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

func (l *Log) establishConnection() {
	l.connTimer = time.NewTimer(connDeadlineIncr)

	go l.closeConnection()

	once.Do(func() {
		conn, err := net.DialTimeout("tcp", logstashHost+":"+logstashPort, connectionTimeout)
		if err != nil {
			errorLogger.Println(fmt.Errorf("connection to Logstash failed: %w", err))

			return
		}

		l.bw = bufio.NewWriter(conn)
		l.connection = conn

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

	if !l.connTimer.Stop() {
		<-l.connTimer.C
	}
	// Reset connection deadline on each successful write. This way the connection will only be closed when there aren't
	// any writes for a certain amount of time.
	l.connTimer.Reset(connDeadlineIncr)

	return err
}

// Finish sends the buffer's contents to Logstash in a batch
func (l *Log) Finish() error {
	l.wg.Wait()

	if l.connection == nil {
		return nil
	}

	err := l.bw.Flush()
	if err != nil {
		err = fmt.Errorf("error flushing logger buffer: %w", err)
		errorLogger.Println(err)
		return err
	}

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
