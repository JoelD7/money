package logger

type ObjectWrapper struct {
	name       string
	properties map[string]interface{}
}

func (o *ObjectWrapper) Key() string {
	return o.name
}

func (o *ObjectWrapper) Value() map[string]interface{} {
	return o.properties
}
