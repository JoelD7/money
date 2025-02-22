package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/JoelD7/money/backend/models"
)

// ConsoleLogger is a logger that writes to stdErr
type ConsoleLogger struct {
	Service string `json:"service,omitempty"`
}

func newConsoleLogger(service string) LogAPI {
	return &ConsoleLogger{Service: service}
}

func (c *ConsoleLogger) Info(eventName string, fields ...models.LoggerField) {
	c.write(infoLevel, eventName, nil, fields)
}

func (c *ConsoleLogger) Warning(eventName string, err error, fields ...models.LoggerField) {
	c.write(warningLevel, eventName, err, fields)
}

func (c *ConsoleLogger) Error(eventName string, err error, fields ...models.LoggerField) {
	c.write(errLevel, eventName, err, fields)
}

func (c *ConsoleLogger) Critical(eventName string, fields ...models.LoggerField) {
	c.write(panicLevel, eventName, nil, fields)
}

func (c *ConsoleLogger) LogLambdaTime(startingTime time.Time, err error, panic interface{}) {
}

func (c *ConsoleLogger) Finish() error {
	return nil
}

func (c *ConsoleLogger) AddToContext(key string, value interface{}) {
	//TODO implement me
	panic("implement me")
}

func (c *ConsoleLogger) SetHandler(handler string) {}

// MapToLoggerObject not implemented
func (c *ConsoleLogger) MapToLoggerObject(name string, m map[string]interface{}) models.LoggerField {
	return models.Any(name, m)
}

func (c *ConsoleLogger) write(level logLevel, eventName string, errToLog error, objects []models.LoggerField) {
	logData := &LogData{
		Service:   c.Service,
		Event:     eventName,
		Level:     string(level),
		LogObject: getLogObjects(objects, nil),
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
