package logger

type AnyField struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func Any(key string, value interface{}) *AnyField {
	return &AnyField{
		Key:   key,
		Value: value,
	}
}

func (a *AnyField) GetKey() string {
	return a.Key
}

func (a *AnyField) GetValue() (interface{}, error) {
	return a.Value, nil
}
