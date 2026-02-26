package utils

import (
	"fmt"
	"net/http"
)

// AppError represents a custom application error with HTTP status code
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code=%d, message=%s, error=%v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code=%d, message=%s", e.Code, e.Message)
}

// NewAppError creates a new AppError
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewBadRequestError creates a new 400 Bad Request error
func NewBadRequestError(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, message, err)
}

// NewUnauthorizedError creates a new 401 Unauthorized error
func NewUnauthorizedError(message string, err error) *AppError {
	return NewAppError(http.StatusUnauthorized, message, err)
}

// NewForbiddenError creates a new 403 Forbidden error
func NewForbiddenError(message string, err error) *AppError {
	return NewAppError(http.StatusForbidden, message, err)
}

// NewNotFoundError creates a new 404 Not Found error
func NewNotFoundError(message string, err error) *AppError {
	return NewAppError(http.StatusNotFound, message, err)
}

// NewInternalServerError creates a new 500 Internal Server Error
func NewInternalServerError(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, message, err)
}

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Details string `json:"details,omitempty"`
}
