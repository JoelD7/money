package logger

type ObjectWrapper struct {
	name       string
	properties map[string]interface{}
}

func (o *ObjectWrapper) LogName() string {
	return o.name
}

func (o *ObjectWrapper) LogProperties() map[string]interface{} {
	return o.properties
}
