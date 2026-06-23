package utils

import "net/http"

type AppError struct {
	Code    string
	Status  int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

const (
	ErrValidation      = "VALIDATION_ERROR"
	ErrUnauthenticated = "UNAUTHENTICATED"
	ErrUnauthorized    = "UNAUTHORIZED"
	ErrNotFound        = "NOT_FOUND"
	ErrConflict        = "CONFLICT"
	ErrInternal        = "INTERNAL_ERROR"
)

func NewAppError(code string, status int, message string) *AppError {
	return &AppError{
		Code:    code,
		Status:  status,
		Message: message,
	}
}

func NewValidationError(message string) *AppError {
	return NewAppError(ErrValidation, http.StatusUnprocessableEntity, message)
}

func NewInternalError(message string) *AppError {
	return NewAppError(ErrInternal, http.StatusInternalServerError, message)
}
