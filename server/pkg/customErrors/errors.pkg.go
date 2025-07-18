/**
 * @File: customerrors.go
 * @Title: Custom Application Error Type
 * @Description: Defines a flexible custom error type (`CustomError`) for application-specific
 * @Description: errors, supporting error wrapping for detailed cause tracking, and including
 * @Description: fields suitable for direct API response mapping.
 * @Author: thesyscoder (github.com/thesyscoder)
 */

package customerrors

import (
	"errors" // Standard library for error handling utilities
	"fmt"
)

// CustomError represents a structured application error.
// It includes a unique code, a human-readable message, an optional
// wrapped underlying error for cause analysis (Go 1.13+ error wrapping),
// an HTTP status code for API responses, and optional additional data.
type CustomError struct {
	Code       string `json:"code"`           // Unique application-specific error code (e.g., "CONFIG_LOAD_FAILED")
	Message    string `json:"message"`        // Human-readable message describing the error
	Err        error  `json:"-"`              // The underlying error, omitted from JSON marshaling
	HTTPStatus int    `json:"-"`              // HTTP status code appropriate for this error, omitted from JSON
	Data       any    `json:"data,omitempty"` // Additional error-specific data for the API response
}

// Error implements the `error` interface for CustomError.
// It returns a formatted string representation of the error.
func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code: %s, message: %s, underlying_error: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

// Unwrap implements the `errors.Unwrap` interface for CustomError.
// This allows `errors.Is` and `errors.As` to be used for inspecting the error chain.
func (e *CustomError) Unwrap() error {
	return e.Err
}

// NewCustomError creates and returns a new CustomError instance.
// `code` is the application-specific error code (e.g., "INVALID_INPUT").
// `message` is a human-readable description.
// `err` is an optional underlying error to wrap.
// `httpStatus` is the HTTP status code to be returned in API responses.
// `data` is optional additional data relevant to the error.
func NewCustomError(code, message string, err error, httpStatus int, data any) *CustomError {
	return &CustomError{
		Code:       code,
		Message:    message,
		Err:        err,
		HTTPStatus: httpStatus,
		Data:       data,
	}
}

// IsCustomError checks if the given error is of type *CustomError.
// This is useful for type assertion and specific error handling, even if wrapped.
func IsCustomError(err error) bool {
	var ce *CustomError
	return errors.As(err, &ce)
}

const (
	// ErrCodeConfigLoadFailed indicates an error occurred while loading application configuration.
	ErrCodeConfigLoadFailed = "CONFIG_LOAD_FAILED"
	// ErrCodeK8sClientInitFailed indicates a failure to initialize the Kubernetes client.
	ErrCodeK8sClientInitFailed = "K8S_CLIENT_INIT_FAILED"
	// ErrCodeK8sClientNotInitialized indicates that the Kubernetes client was requested before being initialized.
	ErrCodeK8sClientNotInitialized = "K8S_CLIENT_NOT_INITIALIZED"
	// ErrCodeStorageInitFailed indicates a failure to initialize storage (e.g., database, file system).
	ErrCodeStorageInitFailed = "STORAGE_INIT_FAILED"
	// ErrCodeInvalidInput indicates that input provided to a function or method was invalid.
	ErrCodeInvalidInput = "INVALID_INPUT"
	// ErrCodeResourceNotFound indicates that a requested resource could not be found.
	ErrCodeResourceNotFound = "RESOURCE_NOT_FOUND"
	// ErrCodeAlreadyExists indicates an attempt to create a resource that already exists.
	ErrCodeAlreadyExists = "RESOURCE_ALREADY_EXISTS"
	// ErrCodeDatabaseOperationFailed indicates a generic failure during a database operation.
	ErrCodeDatabaseOperationFailed = "DATABASE_OPERATION_FAILED"
	// ErrCodeUnauthorized indicates an authentication failure (e.g., invalid credentials, missing token).
	ErrCodeUnauthorized = "UNAUTHORIZED"
	// ErrCodePermissionDenied indicates an authorization failure (user lacks required permissions).
	ErrCodePermissionDenied = "PERMISSION_DENIED"
	// ErrCodeExternalServiceError indicates an error from an external dependency.
	ErrCodeExternalServiceError = "EXTERNAL_SERVICE_ERROR"
	// ErrCodeInternal represents a generic, unexpected internal error.
	ErrCodeInternal = "INTERNAL_ERROR"
)
