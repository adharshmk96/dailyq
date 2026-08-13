package journal

import (
	"errors"
	"net/http"
)

// Error is a domain error carrying the HTTP status and a stable machine code.
type Error struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string { return e.Message }

func newError(status int, code, message string) *Error {
	return &Error{Status: status, Code: code, Message: message}
}

var (
	ErrTaskNotFound   = newError(http.StatusNotFound, "task_not_found", "task not found")
	ErrInvalidRequest = newError(http.StatusBadRequest, "invalid_request", "request payload is invalid")
	ErrInternal       = newError(http.StatusInternalServerError, "internal_error", "something went wrong")
)

// AsError maps any error to a domain Error, defaulting to a 500.
func AsError(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return ErrInternal
}
