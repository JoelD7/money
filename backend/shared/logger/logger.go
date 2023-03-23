package logger

import (
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/utils"
	"net"
)

type logLevel string

const (
	infoLevel    logLevel = "info"
	errLevel     logLevel = "error"
	warningLevel logLevel = "warning"
)

var (
	logstashServerType = env.GetString("LOGSTASH_TYPE", "tcp")
	logstashHost       = env.GetString("LOGSTASH_HOST", "ec2-54-209-181-211.compute-1.amazonaws.com")
	logstashPort       = env.GetString("LOGSTASH_PORT", "5044")
)

type Logger struct {
	Service string `json:"service,omitempty"`
}

type eventLog struct {
	Source string `json:"source,omitempty"`
}

func (l *Logger) Info() {
	l.logData(infoLevel)
}

func (l *Logger) logData(level logLevel) {
	connection, err := net.Dial(logstashServerType, logstashHost+":"+logstashPort)
	if err != nil {
		panic(fmt.Errorf("error connecting to Logstash server: %w", err))
	}

	tk := &tokenEvent{}
	err = json.Unmarshal([]byte(cleanedMessage), tk)
	if err != nil {
		fmt.Println(err)
	}

	lEvent := &logstashEvent{
		Message: *tk,
	}

	str, err := utils.GetJsonString(lEvent)
	if err != nil {
		fmt.Println(err)
		return err
	}

	///send some data
	_, err = connection.Write([]byte(str + "\n"))
	if err != nil {
		fmt.Println("error writing logs")
	}
}
