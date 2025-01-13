package domain

import (
	"github.com/pkg/errors"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

// WrapError to create a new custom error
func WrapError(code int, message string) error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func ExtractCustomError(err error) (*Error, bool) {
	customErr, ok := errors.Cause(err).(*Error)
	return customErr, ok
}
