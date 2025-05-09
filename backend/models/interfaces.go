package models

// LoggerField is an interface that represents a field in a log entry.
type LoggerField interface {
	GetKey() string
	GetValue() (interface{}, error)
}

type AnyField struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// Any creates a new LoggerField with the given key and value.
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

// Resource is an interface that represents any of the types that can be stored in the database. It's purpose is to
// serve as a generics type.
type Resource interface {
	*User | *Category | *Expense | *SavingGoal | *Income | *Period | *Saving
}
