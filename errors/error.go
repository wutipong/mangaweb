package errors

import (
	"fmt"
)

var ErrUnknown = New(0, "unknown error.")
var ErrNotImplemented = New(1, "not implemented.")

type Error struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"cause"`
}

func New(code uint, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}

func (m Error) Is(target error) bool {
	if err, ok := target.(Error); ok {
		return err.Code == m.Code
	}

	return m.Error() == target.Error()
}

func (m Error) Error() string {
	return m.Message
}

func (m Error) Unwrap() error {
	return m.Cause
}

func (m Error) Wrap(cause error) Error {
	return Error{
		Code:    m.Code,
		Message: m.Message,
		Cause:   cause,
	}
}

func (m Error) Format(param ...any) Error {
	return Error{
		Code:    m.Code,
		Message: fmt.Sprintf(m.Message, param...),
		Cause:   m.Cause,
	}
}
