package apigateway

var (
	// ErrInternalError is returned when there's an internal error that must be retried
	ErrInternalError = &Error{
		Code:     10500,
		HTTPCode: 500,
		Message:  "Internal server error, try again later",
	}
)

// Error represents an API error
type Error struct {
	Code     int    `json:"code"`
	HTTPCode int    `json:"http_code"`
	Message  string `json:"message"`
}

// Error returns the error message
func (e *Error) Error() string {
	return e.Message
}

// NewError method to initialize custom error
func NewError(message string, code int) error {
	return &Error{
		Code:     10000 + code,
		HTTPCode: code,
		Message:  message,
	}
}
