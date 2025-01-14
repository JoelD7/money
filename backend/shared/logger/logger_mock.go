package logger

import (
	"bytes"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"time"
)

type LogMock struct {
	Output *bytes.Buffer
}

func (l *LogMock) MapToLoggerObject(name string, m map[string]interface{}) models.LoggerField {
	return models.Any(name, m)
}

// NewLoggerMock mocks the logger client. This is important to prevent unit tests from sending logs.
func NewLoggerMock(buf *bytes.Buffer) *LogMock {
	logMock := &LogMock{
		Output: buf,
	}

	if buf == nil {
		logMock.Output = bytes.NewBuffer(nil)
	}

	return logMock
}

func (l *LogMock) Info(eventName string, objects []models.LoggerField) {
	_, _ = l.Output.Write([]byte(eventName))
}
func (l *LogMock) Warning(eventName string, err error, objects []models.LoggerField) {
	_, _ = l.Output.Write([]byte(eventName))

	if err != nil {
		_, _ = l.Output.Write([]byte(err.Error()))
	}
}
func (l *LogMock) Error(eventName string, err error, objects []models.LoggerField) {
	_, _ = l.Output.Write([]byte(eventName))

	if err != nil {
		_, _ = l.Output.Write([]byte(err.Error()))
	}
}

func (l *LogMock) Critical(eventName string, objects []models.LoggerField) {
	_, _ = l.Output.Write([]byte(eventName))
}
func (*LogMock) LogLambdaTime(startingTime time.Time, err error, panic interface{}) {
	fmt.Println("panic: ", panic)
}

func (l *LogMock) SetHandler(handler string) {}

func (l *LogMock) Finish() error { return nil }
