package cerrors

import "errors"

type (
	Error struct {
		Message      string `json:"message"`
		Code         string `json:"code,omitempty"`
		matcherError error
	}

	JSONError struct {
		Message string `json:"message"`
		Code    string `json:"code,omitempty"`
	}
)

func New(message string, code string) *Error {
	matcherError := errors.New("")

	return &Error{
		Message:      message,
		Code:         code,
		matcherError: matcherError,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Is(target error) bool {
	var targetError *Error

	success := errors.As(target, &targetError)
	if !success {
		return false
	}

	return errors.Is(targetError.matcherError, e.matcherError)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}
