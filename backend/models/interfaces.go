package models

type LoggerObject interface {
	LogName() string
	LogProperties() map[string]interface{}
}
