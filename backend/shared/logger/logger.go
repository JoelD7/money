package logger

import (
	"fmt"
	"time"

	"github.com/JoelD7/money/backend/models"
)

var (
	loggerInstance LogAPI
)

const (
	LogstashImplementation LogImplementation = "logstash"
	ConsoleImplementation  LogImplementation = "console"
)

// InitLogger initializes the logger. id is the identifier used to trace logs of the same request, and impl is the
// implementation to use.
//
// Don't forget to call Finish() when the application is shutting down to ensure a
// proper closing of the logger's resources.
func InitLogger(impl LogImplementation) LogAPI {
	switch impl {
	case LogstashImplementation:
		loggerInstance = newLogstashLogger()
	case ConsoleImplementation:
		loggerInstance = newConsoleLogger("unknown")
	default:
		fmt.Println("Unknown logger implementation, using console logger")
		loggerInstance = newConsoleLogger("unknown")
	}

	return loggerInstance
}

func Info(eventName string, fields ...models.LoggerField) {
	loggerInstance.Info(eventName, fields...)
}

func Warning(eventName string, err error, fields ...models.LoggerField) {
	loggerInstance.Warning(eventName, err, fields...)
}

func Error(eventName string, err error, fields ...models.LoggerField) {
	loggerInstance.Error(eventName, err, fields...)
}

func Critical(eventName string, fields ...models.LoggerField) {
	loggerInstance.Critical(eventName, fields...)
}

func LogLambdaTime(startingTime time.Time, err error, panic interface{}) {
	loggerInstance.LogLambdaTime(startingTime, err, panic)
}

func Finish() error {
	return loggerInstance.Finish()
}

func MapToLoggerObject(name string, m map[string]interface{}) models.LoggerField {
	return loggerInstance.MapToLoggerObject(name, m)
}

func SetHandler(handler string) {
	loggerInstance.SetHandler(handler)
}

func AddToContext(key string, value interface{}) {
	loggerInstance.AddToContext(key, value)
}
