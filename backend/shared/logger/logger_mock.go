package logger

import (
	"bytes"
	"time"
)

type LogMock struct {
	Output *bytes.Buffer
}

// InitLoggerMock mocks the logger client. This is important to prevent unit tests from sending logs.
func InitLoggerMock(buf *bytes.Buffer) *LogMock {
	logMock := &LogMock{
		Output: buf,
	}

	LogClient = logMock

	return logMock
}

func (l *LogMock) Info(eventName string, objects []Object) {
	l.Output.Reset()
	_, _ = l.Output.Write([]byte(eventName))
}
func (l *LogMock) Warning(eventName string, err error, objects []Object) {
	l.Output.Reset()
	_, _ = l.Output.Write([]byte(eventName))
}
func (l *LogMock) Error(eventName string, err error, objects []Object) {
	l.Output.Reset()
	_, _ = l.Output.Write([]byte(eventName))
}
func (l *LogMock) Critical(eventName string, objects []Object) {
	l.Output.Reset()
	_, _ = l.Output.Write([]byte(eventName))
}
func (*LogMock) LogLambdaTime(startingTime time.Time, err error, panic interface{}) {

}
