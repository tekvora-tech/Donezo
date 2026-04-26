// Package utils provides PostgreSQL error parsing utilities.
package parser

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

// ============================================================================
// ParsePostgresError - Parse PostgreSQL-specific errors
// ============================================================================

// ParsePostgresError parses a PostgreSQL error (pgconn.PgError) and returns
// a standardized HTTPError. If the error is not a PostgreSQL error, returns nil.
//
// Logic:
//   1. Try to cast error to *pgconn.PgError using errors.As.
//   2. If not a PostgreSQL error, return nil (caller should use ParseError instead).
//   3. Map PostgreSQL error code (SQLSTATE) to HTTP status and message.
//   4. For unique violations, inspect ConstraintName for specific messages.
//   5. Return standardized HTTPError.
//
// PostgreSQL Error Codes Reference:
//   https://www.postgresql.org/docs/current/errcodes-appendix.html
//
// Parameters:
//   - err: the error to parse
//
// Returns:
//   - *HTTPError if it's a PostgreSQL error, nil otherwise
//
// Usage:
//
//	if httpErr := utils.ParsePostgresError(err); httpErr != nil {
//	    response.Error(c, httpErr.Code, httpErr.Message, httpErr.Error)
//	    return
//	}
//	// fallback to generic error parsing
//	message, code, detail := utils.ParseError(err)
//	response.Error(c, code, message, detail)
//
// Simpen di mana:
//   - Error code mapping -> hardcoded di function ini.
//   - Constraint name mapping -> hardcoded di switch case.
//   - Tidak ada state yang disimpan.
func ParsePostgresError(err error) *HTTPError {
	var pgErr *pgconn.PgError

	// Coba cast ke pgconn.PgError
	if !errors.As(err, &pgErr) {
		return nil
	}

	switch pgErr.Code {

	// ===================================================================
	// VALIDATION ERRORS (4xx Client Error)
	// ===================================================================

	// 23502 - not_null_violation
	// Terjadi ketika INSERT/UPDATE tanpa nilai untuk kolom NOT NULL
	// Contoh: INSERT INTO users (email) VALUES (NULL)
	case "23502":
		column := extractColumnName(pgErr.ColumnName, pgErr.Detail)
		return &HTTPError{
			Status:  http.StatusBadRequest,
			Code:    http.StatusBadRequest,
			Message: "missing required field",
			Err:   "field '" + column + "' cannot be empty",
		}

	// 23514 - check_violation
	// Terjadi ketika CHECK constraint tidak terpenuhi
	// Contoh: CHECK (age >= 0) tapi diisi -1
	case "23514":
		return &HTTPError{
			Status:  http.StatusBadRequest,
			Code:    http.StatusBadRequest,
			Message: "validation failed",
			Err:   "check constraint violated: " + pgErr.ConstraintName,
		}

	// 22001 - string_data_right_truncation
	// Terjadi ketika string terlalu panjang untuk kolom
	// Contoh: VARCHAR(50) tapi diisi 100 karakter
	case "22001":
		return &HTTPError{
			Status:  http.StatusBadRequest,
			Code:    http.StatusBadRequest,
			Message: "validation failed",
			Err:   "value too long for field",
		}

	// 22003 - numeric_value_out_of_range
	// Terjadi ketika numeric value di luar range
	case "22003":
		return &HTTPError{
			Status:  http.StatusBadRequest,
			Code:    http.StatusBadRequest,
			Message: "validation failed",
			Err:   "numeric value out of range",
		}

	// 22P02 - invalid_text_representation
	// Terjadi ketika format data tidak valid (misal: UUID salah format)
	case "22P02":
		return &HTTPError{
			Status:  http.StatusBadRequest,
			Code:    http.StatusBadRequest,
			Message: "invalid request",
			Err:   "invalid data format",
		}

	// ===================================================================
	// CONFLICT ERRORS (409 Conflict)
	// ===================================================================

	// 23505 - unique_violation
	// Terjadi ketika UNIQUE constraint atau PRIMARY KEY violated
	// Contoh: INSERT email yang sudah ada
	case "23505":
		msg := "resource already exists"

		// Constraint-based message untuk Donezo
		switch pgErr.ConstraintName {
		// Users table
		case "users_email_key", "users_email_unique":
			msg = "email already registered"
		case "users_username_key", "users_username_unique":
			msg = "username already taken"

		// Tags table
		case "tags_user_id_name_key":
			msg = "tag name already exists for this user"

		// Todo-Tags junction
		case "todo_tags_pkey":
			msg = "tag already assigned to this todo"
		}

		return &HTTPError{
			Status:  http.StatusConflict,
			Code:    http.StatusConflict,
			Message: msg,
			Err:   "unique constraint violated",
		}

	// 23503 - foreign_key_violation
	// Terjadi ketika FOREIGN KEY reference tidak ditemukan
	// Contoh: INSERT todo dengan user_id yang tidak ada
	case "23503":
		return &HTTPError{
			Status:  http.StatusConflict,
			Code:    http.StatusConflict,
			Message: "invalid reference",
			Err:   "foreign key constraint violated: " + extractFKDetail(pgErr.Detail),
		}

	// ===================================================================
	// AUTH & PERMISSION ERRORS (403 Forbidden)
	// ===================================================================

	// 42501 - insufficient_privilege
	// Terjadi ketika user tidak punya privilege untuk operasi
	case "42501":
		return &HTTPError{
			Status:  http.StatusForbidden,
			Code:    http.StatusForbidden,
			Message: "access denied",
			Err:   "insufficient privilege",
		}

	// 28P01 - invalid_password
	// Terjadi ketika password authentication gagal (DB level)
	case "28P01":
		return &HTTPError{
			Status:  http.StatusUnauthorized,
			Code:    http.StatusUnauthorized,
			Message: "authentication failed",
			Err:   "invalid database credentials",
		}

	// ===================================================================
	// TRANSACTION ERRORS (503 Service Unavailable)
	// ===================================================================

	// 40001 - serialization_failure
	// Terjadi ketika concurrent transaction conflict (serialization isolation)
	case "40001":
		return &HTTPError{
			Status:  http.StatusServiceUnavailable,
			Code:    http.StatusServiceUnavailable,
			Message: "please retry request",
			Err:   "transaction serialization failure",
		}

	// 40P01 - deadlock_detected
	// Terjadi ketika deadlock terdeteksi dan transaction di-rollback
	case "40P01":
		return &HTTPError{
			Status:  http.StatusServiceUnavailable,
			Code:    http.StatusServiceUnavailable,
			Message: "please retry request",
			Err:   "deadlock detected",
		}

	// ===================================================================
	// TIMEOUT & CANCEL ERRORS (408/504)
	// ===================================================================

	// 57014 - query_canceled
	// Terjadi ketika query di-cancel (timeout atau explicit cancel)
	case "57014":
		return &HTTPError{
			Status:  http.StatusRequestTimeout,
			Code:    http.StatusRequestTimeout,
			Message: "request timeout",
			Err:   "query canceled by timeout",
		}

	// 57015 - query_canceled (alternative)
	case "57015":
		return &HTTPError{
			Status:  http.StatusRequestTimeout,
			Code:    http.StatusRequestTimeout,
			Message: "request timeout",
			Err:   "query canceled",
		}

	// ===================================================================
	// SCHEMA / QUERY ERRORS (500 Internal Server Error)
	// ===================================================================

	// 42P01 - undefined_table
	// Terjadi ketika query ke table yang tidak ada
	case "42P01":
		return &HTTPError{
			Status:  http.StatusInternalServerError,
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Err:   "database table not found",
		}

	// 42703 - undefined_column
	// Terjadi ketika query ke kolom yang tidak ada
	case "42703":
		return &HTTPError{
			Status:  http.StatusInternalServerError,
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Err:   "database column not found",
		}

	// 42883 - undefined_function
	// Terjadi ketika function tidak ditemukan
	case "42883":
		return &HTTPError{
			Status:  http.StatusInternalServerError,
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Err:   "database function not found",
		}

	// 42601 - syntax_error
	// Terjadi ketika SQL syntax salah
	case "42601":
		return &HTTPError{
			Status:  http.StatusInternalServerError,
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Err:   "database syntax error",
		}

	// ===================================================================
	// CONNECTION ERRORS (503 Service Unavailable)
	// ===================================================================

	// 08006 - connection_failure
	// Terjadi ketika koneksi ke database gagal
	case "08006":
		return &HTTPError{
			Status:  http.StatusServiceUnavailable,
			Code:    http.StatusServiceUnavailable,
			Message: "database unavailable",
			Err:   "connection failure",
		}

	// 08001 - sqlclient_unable_to_establish_sqlconnection
	case "08001":
		return &HTTPError{
			Status:  http.StatusServiceUnavailable,
			Code:    http.StatusServiceUnavailable,
			Message: "database unavailable",
			Err:   "unable to establish connection",
		}

	// 08004 - sqlserver_rejected_establishment_of_sqlconnection
	case "08004":
		return &HTTPError{
			Status:  http.StatusServiceUnavailable,
			Code:    http.StatusServiceUnavailable,
			Message: "database unavailable",
			Err:   "connection rejected by server",
		}

	// 53300 - too_many_connections
	// Terjadi ketika connection pool penuh
	case "53300":
		return &HTTPError{
			Status:  http.StatusServiceUnavailable,
			Code:    http.StatusServiceUnavailable,
			Message: "service overloaded",
			Err:   "too many database connections",
		}

	// ===================================================================
	// DATA CORRUPTION (500 Internal Server Error)
	// ===================================================================

	// XX000 - internal_error
	// Generic internal PostgreSQL error
	case "XX000":
		return &HTTPError{
			Status:  http.StatusInternalServerError,
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Err:   "database internal error",
		}

	// XX001 - data_corrupted
	case "XX001":
		return &HTTPError{
			Status:  http.StatusInternalServerError,
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Err:   "data corrupted",
		}

	// XX002 - index_corrupted
	case "XX002":
		return &HTTPError{
			Status:  http.StatusInternalServerError,
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Err:   "index corrupted",
		}
	}

	// ===================================================================
	// Fallback: Unknown PostgreSQL error
	// ===================================================================
	return &HTTPError{
		Status:  http.StatusInternalServerError,
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
		Err:   "postgres error [" + pgErr.Code + "]: " + pgErr.Message,
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

// extractColumnName extracts column name from pgconn error detail.
func extractColumnName(columnName, detail string) string {
	if columnName != "" {
		return columnName
	}
	// Try to extract from detail message
	// Detail format: "Failing row contains (..., null, ...)"
	if strings.Contains(detail, "column") {
		parts := strings.Split(detail, "\"")
		if len(parts) >= 2 {
			return parts[1]
		}
	}
	return "unknown"
}

// extractFKDetail extracts foreign key detail from error message.
func extractFKDetail(detail string) string {
	// Detail format: "Key (user_id)=(xxx) is not present in table \"users\""
	if detail == "" {
		return "foreign key constraint violated"
	}
	return detail
}

// IsPostgresError checks if an error is a PostgreSQL error.
func IsPostgresError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr)
}

// GetPostgresCode returns the PostgreSQL error code if it's a pgconn error.
// Returns empty string if not a PostgreSQL error.
func GetPostgresErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

// IsUniqueViolation checks if the error is a unique constraint violation.
func IsUniqueViolation(err error) bool {
	return GetPostgresErrorCode(err) == "23505"
}

// IsForeignKeyViolation checks if the error is a foreign key violation.
func IsForeignKeyViolation(err error) bool {
	return GetPostgresErrorCode(err) == "23503"
}

// IsNotNullViolation checks if the error is a not-null constraint violation.
func IsNotNullViolation(err error) bool {
	return GetPostgresErrorCode(err) == "23502"
}