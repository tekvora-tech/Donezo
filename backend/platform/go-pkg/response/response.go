package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTTPResponse format standar API response
type HTTPResponse struct {
	Code    int         `json:"code"`              // HTTP status code
	Message string      `json:"message"` // Pesan untuk user
	Error string `json:"error,omitempty"`
	Result    interface{} `json:"result,omitempty"`    // Hasil
}

// ─── Success Responses ────────────────────────

// OK response sukses (200)
func OK(c *gin.Context, message string, result interface{}) {
	c.JSON(http.StatusOK, HTTPResponse{
		Code:    http.StatusOK,
		Message: message,
		Result:    result,
	})
}

// Created resource berhasil dibuat (201)
func Created(c *gin.Context, message string, result interface{}) {
	c.JSON(http.StatusCreated, HTTPResponse{
		Code:    http.StatusCreated,
		Message: message,
		Result:    result,
	})
}

// Accepted request diterima, proses async (202)
func Accepted(c *gin.Context, message string, result interface{}) {
	c.JSON(http.StatusAccepted, HTTPResponse{
		Code:    http.StatusAccepted,
		Message: message,
		Result:    result,
	})
}

// NoContent sukses tanpa data (204)
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// ─── Error Responses ──────────────────────────

// BadRequest input user salah (400)
func BadRequest(c *gin.Context, message, err string) {
	c.JSON(http.StatusBadRequest, HTTPResponse{
		Code:    http.StatusBadRequest,
		Message: message,
		Error: err,
	})
}

// Unauthorized belum login / token invalid (401)
func Unauthorized(c *gin.Context, message, err string) {
	if message == "" {
		message = "unauthorized"
	}
	c.JSON(http.StatusUnauthorized, HTTPResponse{
		Code:    http.StatusUnauthorized,
		Message: message,
		Error: err,
	})
}

// Forbidden nggak punya izin (403)
func Forbidden(c *gin.Context, message, err string) {
	if message == "" {
		message = "forbidden"
	}
	c.JSON(http.StatusForbidden, HTTPResponse{
		Code:    http.StatusForbidden,
		Message: message,
		Error: err,
	})
}

// NotFound resource tidak ada (404)
func NotFound(c *gin.Context, message, err string) {
	if message == "" {
		message = "resource not found"
	}
	c.JSON(http.StatusNotFound, HTTPResponse{
		Code:    http.StatusNotFound,
		Message: message,
		Error: err,
	})
}

// Conflict data bentrok (409)
func Conflict(c *gin.Context, message, err string) {
	c.JSON(http.StatusConflict, HTTPResponse{
		Code:    http.StatusConflict,
		Message: message,
		Error: err,
	})
}

// UnprocessableEntity validasi gagal (422)
func UnprocessableEntity(c *gin.Context, message, err string) {
	c.JSON(http.StatusUnprocessableEntity, HTTPResponse{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
		Error: err,
	})
}

// TooManyRequests rate limit (429)
func TooManyRequests(c *gin.Context, message, err string) {
	c.JSON(http.StatusTooManyRequests, HTTPResponse{
		Code:    http.StatusTooManyRequests,
		Message: message,
		Error: err,
	})
}

// Error error umum dengan custom status code
func Error(c *gin.Context, code int, message, err string) {
	c.JSON(code, HTTPResponse{
		Code:    code,
		Message: message,
		Error: err,
	})
}

// InternalServerError error server (500)
func InternalServerError(c *gin.Context, message, err string) {
	if message == "" {
		message = "internal server error"
	}
	c.JSON(http.StatusInternalServerError, HTTPResponse{
		Code:    http.StatusInternalServerError,
		Message: message,
		Error: err,
	})
}

// ServiceUnavailable server sedang maintenance (503)
func ServiceUnavailable(c *gin.Context, message, err string) {
	c.JSON(http.StatusServiceUnavailable, HTTPResponse{
		Code:    http.StatusServiceUnavailable,
		Message: message,
		Error: err,
	})
}

// ─── Pagination Response ──────────────────────

// Paginated response dengan metadata pagination
func Paginated(c *gin.Context, message string, result interface{}) {
	c.JSON(http.StatusOK, HTTPResponse{
		Code:    http.StatusOK,
		Message: message,
		Result:    result,
	})
}