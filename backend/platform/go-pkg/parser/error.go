// Package utils provides utility functions for the Donezo application.
package parser

import (
	"net/http"
	"strings"
)

// HTTPError represents a standardized HTTP error response.
type HTTPError struct {
	Status  int    `json:"-"`              // HTTP status code (for Gin)
	Code    int    `json:"code"`           // HTTP status code (in response body)
	Message string `json:"message"`        // User-friendly message
	Err   string `json:"error"`          // Technical detail
}

// Error implements the error interface for HTTPError.
func (e *HTTPError) Error() string {
	return e.Err
}

// ToResponse returns the HTTPError as a map for JSON response.
func (e *HTTPError) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"code":    e.Code,
		"message": e.Message,
		"error":   e.Error,
	}
}

// ============================================================================
// ParseError - Parse generic application errors
// ============================================================================

// ParseError parses a generic error and returns standardized HTTP error components.
//
// Logic:
//   1. If err is nil, return success (200, "", "").
//   2. Extract the error message string.
//   3. Split by ":" to get key and detail.
//   4. Map the key to appropriate HTTP status code and message.
//   5. If key doesn't match any known pattern, fallback to string matching.
//
// Parameters:
//   - err: the error to parse
//
// Returns:
//   - message: user-friendly message
//   - code: HTTP status code
//   - detail: technical error detail
//
// Usage:
//
//	message, code, detail := utils.ParseError(err)
//	response.Error(c, code, message, detail)
//
// Simpen di mana:
//   - Error messages → hardcoded di function ini.
//   - Tidak ada state yang disimpan.
func ParseError(err error) (message string, code int, detail string) {
	if err == nil {
		return "success", http.StatusOK, ""
	}

	raw := err.Error()

	// Default values
	message = "internal server error"
	code = http.StatusInternalServerError
	detail = raw

	// Split by ":" to extract key and detail
	parts := strings.SplitN(raw, ":", 2)

	var key string
	if len(parts) > 0 {
		key = strings.TrimSpace(strings.ToLower(parts[0]))
	}

	if len(parts) == 2 {
		detail = strings.TrimSpace(parts[1])
	}

	// Mapping berdasarkan key (prefix error)
	switch key {
	case "invalid request":
		message = "invalid request"
		code = http.StatusBadRequest

	case "not found":
		message = "data not found"
		code = http.StatusNotFound

	case "unauthorized":
		message = "unauthorized"
		code = http.StatusUnauthorized

	case "forbidden":
		message = "forbidden"
		code = http.StatusForbidden

	case "conflict":
		message = "data already exists"
		code = http.StatusConflict

	case "validation failed":
		message = "validation failed"
		code = http.StatusBadRequest

	case "timeout":
		message = "request timeout"
		code = http.StatusRequestTimeout

	case "unprocessable":
		message = "unprocessable request"
		code = http.StatusUnprocessableEntity

	case "too many requests":
		message = "too many requests"
		code = http.StatusTooManyRequests

	case "service unavailable":
		message = "service unavailable"
		code = http.StatusServiceUnavailable

	default:
		// Fallback detection dari full error string (case-insensitive)
		lower := strings.ToLower(raw)
		switch {
		case strings.Contains(lower, "duplicate"),
			strings.Contains(lower, "unique"),
			strings.Contains(lower, "already exists"):
			message = "data already exists"
			code = http.StatusConflict

		case strings.Contains(lower, "not found"),
			strings.Contains(lower, "does not exist"),
			strings.Contains(lower, "no rows"):
			message = "data not found"
			code = http.StatusNotFound

		case strings.Contains(lower, "invalid"),
			strings.Contains(lower, "bad request"),
			strings.Contains(lower, "malformed"):
			message = "invalid request"
			code = http.StatusBadRequest

		case strings.Contains(lower, "unauthorized"),
			strings.Contains(lower, "unauthenticated"):
			message = "unauthorized"
			code = http.StatusUnauthorized

		case strings.Contains(lower, "forbidden"),
			strings.Contains(lower, "permission denied"):
			message = "forbidden"
			code = http.StatusForbidden

		case strings.Contains(lower, "timeout"),
			strings.Contains(lower, "deadline exceeded"):
			message = "request timeout"
			code = http.StatusRequestTimeout

		case strings.Contains(lower, "connection"),
			strings.Contains(lower, "refused"),
			strings.Contains(lower, "unavailable"):
			message = "service unavailable"
			code = http.StatusServiceUnavailable

		case strings.Contains(lower, "validation"),
			strings.Contains(lower, "required"),
			strings.Contains(lower, "cannot be empty"):
			message = "validation failed"
			code = http.StatusBadRequest
		}
	}

	return message, code, detail
}

// ParseErrorToHTTPError wraps ParseError and returns an HTTPError struct.
//
// Usage:
//
//	httpErr := utils.ParseErrorToHTTPError(err)
//	response.Error(c, httpErr.Code, httpErr.Message, httpErr.Error)
func ParseErrorToHTTPError(err error) *HTTPError {
	message, code, detail := ParseError(err)
	return &HTTPError{
		Status:  code,
		Code:    code,
		Message: message,
		Err:   detail,
	}
}

// NewHTTPError creates a new HTTPError with the given parameters.
//
// Usage:
//
//	return utils.NewHTTPError(http.StatusNotFound, "Todo not found", "todo does not exist")
func NewHTTPError(status int, message, err string) *HTTPError {
	return &HTTPError{
		Status:  status,
		Code:    status,
		Message: message,
		Err:   err,
	}
}

// IsNotFound checks if the error represents a "not found" response.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, code, _ := ParseError(err)
	return code == http.StatusNotFound
}

// IsConflict checks if the error represents a "conflict" response.
func IsConflict(err error) bool {
	if err == nil {
		return false
	}
	_, code, _ := ParseError(err)
	return code == http.StatusConflict
}

// IsUnauthorized checks if the error represents an "unauthorized" response.
func IsUnauthorized(err error) bool {
	if err == nil {
		return false
	}
	_, code, _ := ParseError(err)
	return code == http.StatusUnauthorized
}