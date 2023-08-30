package logger

import (
	"bytes"
	"github.com/JoelD7/money/backend/models"
	"time"
)

type LogMock struct {
	Output *bytes.Buffer
}

func (l *LogMock) MapToLoggerObject(name string, m map[string]interface{}) models.LoggerObject {
	return &ObjectWrapper{
		name:       name,
		properties: m,
	}
}

// NewLoggerMock mocks the logger client. This is important to prevent unit tests from sending logs.
func NewLoggerMock(buf *bytes.Buffer) *LogMock {
	logMock := &LogMock{
		Output: buf,
	}

	if buf == nil {
		logMock.Output = bytes.NewBuffer(nil)
	}

	LogClient = logMock

	return logMock
}

func (l *LogMock) Info(eventName string, objects ...models.LoggerObject) {
	_, _ = l.Output.Write([]byte(eventName))
}
func (l *LogMock) Warning(eventName string, err error, objects ...models.LoggerObject) {
	_, _ = l.Output.Write([]byte(eventName))
	_, _ = l.Output.Write([]byte(err.Error()))
}
func (l *LogMock) Error(eventName string, err error, objects ...models.LoggerObject) {
	_, _ = l.Output.Write([]byte(eventName))
	_, _ = l.Output.Write([]byte(err.Error()))
}

func (l *LogMock) Critical(eventName string, objects ...models.LoggerObject) {
	_, _ = l.Output.Write([]byte(eventName))
}
func (*LogMock) LogLambdaTime(startingTime time.Time, err error, panic interface{}) {}

func (l *LogMock) Close() error { return nil }
