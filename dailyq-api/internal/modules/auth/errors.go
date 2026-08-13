package auth

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
	ErrEmailTaken         = newError(http.StatusConflict, "email_taken", "an account with this email already exists")
	ErrInvalidCredentials = newError(http.StatusUnauthorized, "invalid_credentials", "invalid email or password")
	ErrUserNotFound       = newError(http.StatusNotFound, "user_not_found", "user not found")
	ErrUnauthorized       = newError(http.StatusUnauthorized, "unauthorized", "authentication required")
	ErrInvalidToken       = newError(http.StatusUnauthorized, "invalid_token", "invalid or expired token")
	ErrSessionExpired     = newError(http.StatusUnauthorized, "session_expired", "session expired or revoked")
	ErrInvalidResetToken  = newError(http.StatusBadRequest, "invalid_reset_token", "invalid or expired reset token")
	ErrSamePassword       = newError(http.StatusBadRequest, "same_password", "new password must differ from the current password")
	ErrWeakPassword       = newError(http.StatusBadRequest, "weak_password", "password does not meet the minimum requirements")
	ErrInvalidRequest     = newError(http.StatusBadRequest, "invalid_request", "request payload is invalid")
	ErrInternal           = newError(http.StatusInternalServerError, "internal_error", "something went wrong")
)

// AsError maps any error to a domain Error, defaulting to a 500.
func AsError(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return ErrInternal
}
