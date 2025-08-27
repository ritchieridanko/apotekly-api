package ce

import "fmt"

type Error struct {
	Code    int
	Message string
	Err     error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func NewError(code int, message, tracer string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     fmt.Errorf("%s > %s: %w", tracer, err.Error(), err),
	}
}
