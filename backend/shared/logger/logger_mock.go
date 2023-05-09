package logger

import (
	"time"
)

type LogMock struct{}

// InitLoggerMock mocks the logger client. This is important to prevent unit tests from sending logs.
func InitLoggerMock() {
	LogClient = &LogMock{}
}

func (*LogMock) Info(eventName string, objects []Object)                 {}
func (*LogMock) Warning(eventName string, err error, objects []Object)   {}
func (*LogMock) Error(eventName string, err error, objects []Object)     {}
func (*LogMock) Critical(eventName string, objects []Object)             {}
func (*LogMock) LogLambdaTime(startingTime time.Time, panic interface{}) {}
