package errors

import (
	"errors"
	"fmt"
)

// Error standar untuk semua service
var (
	// 4xx Client Errors
	ErrBadRequest       = errors.New("bad request")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrNotFound         = errors.New("resource not found")
	ErrConflict         = errors.New("resource conflict")
	ErrUnprocessable    = errors.New("unprocessable entity")
	ErrTooManyRequests  = errors.New("too many requests")

	// 5xx Server Errors
	ErrInternal         = errors.New("internal server error")
	ErrNotImplemented   = errors.New("not implemented")
	ErrServiceUnavailable = errors.New("service unavailable")
)

// AppError error dengan context tambahan
type AppError struct {
	Err     error  // Error asli
	Message string // Pesan untuk user
	Code    string // Error code: "AUTH_001", "USER_404"
	Status  int    // HTTP status code
}

// Error implementasi error interface
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

// Unwrap untuk errors.Is
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError buat AppError baru
func NewAppError(err error, message, code string, status int) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    code,
		Status:  status,
	}
}

// Wrap tambah context ke error
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// Is cek error sama atau nggak (support wrapping)
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As cast error ke type tertentu
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}