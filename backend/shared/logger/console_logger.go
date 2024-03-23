package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/JoelD7/money/backend/models"
)

var (
	errorLogger  = log.New(os.Stderr, "ERROR ", log.Llongfile)
	backupLogger = log.New(os.Stdout, "", log.Llongfile)
)

// ConsoleLogger is a logger that writes to stdErr
type ConsoleLogger struct {
	Service string `json:"service,omitempty"`
}

func NewConsoleLogger(service string) LogAPI {
	return &ConsoleLogger{Service: service}
}

func (c *ConsoleLogger) Info(eventName string, objects []models.LoggerObject) {
	c.write(infoLevel, eventName, nil, objects)
}

func (c *ConsoleLogger) Warning(eventName string, err error, objects []models.LoggerObject) {
	c.write(warningLevel, eventName, err, objects)
}

func (c *ConsoleLogger) Error(eventName string, err error, objects []models.LoggerObject) {
	c.write(errLevel, eventName, err, objects)
}

func (c *ConsoleLogger) Critical(eventName string, objects []models.LoggerObject) {
	c.write(panicLevel, eventName, nil, objects)
}

func (c *ConsoleLogger) LogLambdaTime(startingTime time.Time, err error, panic interface{}) {}

func (c *ConsoleLogger) Close() error {
	return nil
}

func (c *ConsoleLogger) SetHandler(handler string) {}

// MapToLoggerObject not implemented
func (c *ConsoleLogger) MapToLoggerObject(name string, m map[string]interface{}) models.LoggerObject {
	return &ObjectWrapper{
		name:       name,
		properties: m,
	}
}

func (c *ConsoleLogger) write(level logLevel, eventName string, errToLog error, objects []models.LoggerObject) {
	logData := &LogData{
		Service:   c.Service,
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
		errorLogger.Println(fmt.Errorf("console logger: error encoding log data: %w", err))
		return
	}

	errorLogger.Println(dataAsBytes.String())
}
