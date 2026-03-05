package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	// 4xx Client Errors
	CodeBadRequest       ErrorCode = "BAD_REQUEST"
	CodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	CodeForbidden        ErrorCode = "FORBIDDEN"
	CodeNotFound         ErrorCode = "NOT_FOUND"
	CodeConflict         ErrorCode = "CONFLICT"
	CodeValidationFailed ErrorCode = "VALIDATION_FAILED"

	// 5xx Server Errors
	CodeInternalServer      ErrorCode = "INTERNAL_SERVER_ERROR"
	CodeDatabaseError       ErrorCode = "DATABASE_ERROR"
	CodeExternalServiceDown ErrorCode = "EXTERNAL_SERVICE_DOWN"
)

// AppError - структурированная ошибка приложения
type AppError struct {
	Code    ErrorCode         `json:"code"`
	Message string            `json:"message"`
	Op      string            `json:"op,omitempty"`
	Err     error             `json:"-"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case CodeBadRequest, CodeValidationFailed:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func NotFoundError(op string, msg string) *AppError {
	return &AppError{
		Code:    CodeNotFound,
		Message: msg,
		Op:      op,
	}
}

func ValidationError(op string, fields map[string]string) *AppError {
	return &AppError{
		Code:    CodeValidationFailed,
		Message: "validation failed",
		Op:      op,
		Fields:  fields,
	}
}

func DatabaseError(op string, err error) *AppError {
	return &AppError{
		Code:    CodeDatabaseError,
		Message: "database error",
		Op:      op,
		Err:     err,
	}
}

func InternalError(op string, err error) *AppError {
	return &AppError{
		Code:    CodeInternalServer,
		Message: "internal error",
		Op:      op,
		Err:     err,
	}
}

func IsNotFound(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == CodeNotFound
	}
	return false
}

func UnauthorizedError(op string, msg string) *AppError {
	return &AppError{
		Code:    CodeUnauthorized,
		Message: msg,
		Op:      op,
	}
}

func ConflictError(op string, msg string) *AppError {
	return &AppError{
		Code:    CodeConflict,
		Message: msg,
		Op:      op,
	}
}
