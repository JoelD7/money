package logger

import (
	"fmt"
	"time"

	"github.com/JoelD7/money/backend/models"
)

var (
	loggerSingleton LogAPI
)

const (
	LogstashImplementation LogImplementation = "logstash"
	ConsoleImplementation  LogImplementation = "console"
)

// LogImplementation is the type of logger to use
type LogImplementation string

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

// InitLogger initializes the logger. Don't forget to call Finish() when the application is shutting down to ensure a
// proper closing of the logger's resources.
func InitLogger(impl LogImplementation) LogAPI {
	if loggerSingleton != nil {
		return loggerSingleton
	}

	switch impl {
	case LogstashImplementation:
		loggerSingleton = initLogstash()
	case ConsoleImplementation:
		loggerSingleton = initConsole("unknown")
	default:
		fmt.Println("Unknown logger implementation, using console logger")
		loggerSingleton = initConsole("unknown")
	}

	return loggerSingleton
}

func Info(eventName string, fields ...models.LoggerField) {
	loggerSingleton.Info(eventName, fields...)
}

func Warning(eventName string, err error, fields ...models.LoggerField) {
	loggerSingleton.Warning(eventName, err, fields...)
}

func Error(eventName string, err error, fields ...models.LoggerField) {
	loggerSingleton.Error(eventName, err, fields...)
}

func Critical(eventName string, fields ...models.LoggerField) {
	loggerSingleton.Critical(eventName, fields...)
}

func LogLambdaTime(startingTime time.Time, err error, panic interface{}) {
	loggerSingleton.LogLambdaTime(startingTime, err, panic)
}

func Finish() error {
	return loggerSingleton.Finish()
}

func MapToLoggerObject(name string, m map[string]interface{}) models.LoggerField {
	return loggerSingleton.MapToLoggerObject(name, m)
}

func SetHandler(handler string) {
	loggerSingleton.SetHandler(handler)
}
