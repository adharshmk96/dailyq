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
	ErrNoteNotFound   = newError(http.StatusNotFound, "note_not_found", "note not found")
	ErrTagNotFound    = newError(http.StatusNotFound, "tag_not_found", "tag not found")
	ErrTagNameTaken   = newError(http.StatusConflict, "tag_name_taken", "a tag with that name already exists")
	ErrUnknownTag     = newError(http.StatusBadRequest, "unknown_tag", "one or more tags do not exist")
	ErrInvalidDate    = newError(http.StatusBadRequest, "invalid_date", "date must be in YYYY-MM-DD format")
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
