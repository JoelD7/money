package logger

import (
	"github.com/JoelD7/money/backend/models"
	"time"
)

// LogImplementation is the type of logger to use
type LogImplementation string

type LogAPI interface {
	// Info logs an event with info level
	Info(eventName string, fields ...models.LoggerField)

	// Warning logs an event with warning level
	Warning(eventName string, err error, fields ...models.LoggerField)

	// Error logs an event with error level
	Error(eventName string, err error, fields ...models.LoggerField)

	// Critical logs an event with critical level
	Critical(eventName string, fields ...models.LoggerField)

	// LogLambdaTime logs the time it took for a lambda to execute
	LogLambdaTime(startingTime time.Time, err error, panic interface{})

	// Finish closes the logger's resources
	Finish() error

	// AddToContext adds a key-value pair to the logger's context so that it can be logged with every event
	AddToContext(key string, value interface{})

	MapToLoggerObject(name string, m map[string]interface{}) models.LoggerField
	SetHandler(handler string)
}
