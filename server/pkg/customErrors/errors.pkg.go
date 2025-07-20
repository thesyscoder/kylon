/**
 * @File: customerrors.go
 * @Title: Custom Application Error Type
 * @Description: Defines a flexible custom error type (`CustomError`) for application-specific
 * errors, supporting error wrapping for detailed cause tracking, and including
 * fields suitable for direct API response mapping.
 * @Author: thesyscoder (github.com/thesyscoder)
 */

package customerrors

import (
	"errors"
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
	HTTPStatus int    `json:"-"`              // HTTP status code for this error, omitted from JSON
	Data       any    `json:"data,omitempty"` // Additional error-specific data for the API response
}

// Error implements the error interface for CustomError.
func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code: %s, message: %s, underlying_error: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

// Unwrap allows error unwrapping for errors.Is/errors.As.
func (e *CustomError) Unwrap() error {
	return e.Err
}

// NewCustomError returns a new instance of CustomError.
func NewCustomError(code, message string, err error, httpStatus int, data interface{}) *CustomError {
	return &CustomError{
		Code:       code,
		Message:    message,
		Err:        err,
		HTTPStatus: httpStatus,
		Data:       data,
	}
}

// IsCustomError checks if the given error is a *CustomError, even if wrapped.
func IsCustomError(err error) bool {
	var ce *CustomError
	return errors.As(err, &ce)
}

// APPLICATION ERROR CODES (expand as needed for your domain)
const (
	// Configuration & Startup
	ErrCodeConfigLoadFailed       = "CONFIG_LOAD_FAILED"
	ErrCodeConfigValidationFailed = "CONFIG_VALIDATION_FAILED"
	ErrCodeEnvMissingRequired     = "ENV_MISSING_REQUIRED"
	ErrCodeSecretLoadFailed       = "SECRET_LOAD_FAILED"

	// Database/Storage
	ErrCodeStorageInitFailed        = "STORAGE_INIT_FAILED"
	ErrCodeDatabaseConnectionFailed = "DATABASE_CONNECTION_FAILED"
	ErrCodeDatabaseOperationFailed  = "DATABASE_OPERATION_FAILED"
	ErrCodeMigrationFailed          = "MIGRATION_FAILED"
	ErrCodeTransactionFailed        = "TRANSACTION_FAILED"

	// Kubernetes/Cloud APIs
	ErrCodeK8sClientInitFailed      = "K8S_CLIENT_INIT_FAILED"
	ErrCodeK8sClientNotInitialized  = "K8S_CLIENT_NOT_INITIALIZED"
	ErrCodeK8sResourceError         = "K8S_RESOURCE_ERROR"
	ErrCodeK8sAPIError              = "K8S_API_ERROR"
	ErrCodeCloudProviderError       = "CLOUD_PROVIDER_ERROR"
	ErrCodeObjectStorageUnavailable = "OBJECT_STORAGE_UNAVAILABLE"

	// Application/Business Logic
	ErrCodeInvalidInput        = "INVALID_INPUT"
	ErrCodeResourceNotFound    = "RESOURCE_NOT_FOUND"
	ErrCodeAlreadyExists       = "RESOURCE_ALREADY_EXISTS"
	ErrCodePermissionDenied    = "PERMISSION_DENIED"
	ErrCodeUnauthorized        = "UNAUTHORIZED"
	ErrCodeRateLimitExceeded   = "RATE_LIMIT_EXCEEDED"
	ErrCodeOperationTimeout    = "OPERATION_TIMEOUT"
	ErrCodeOperationInProgress = "OPERATION_IN_PROGRESS"
	ErrCodeDependencyFailed    = "DEPENDENCY_FAILED"
	ErrCodeBackupFailed        = "BACKUP_FAILED"
	ErrCodeRestoreFailed       = "RESTORE_FAILED"
	ErrCodeSnapshotFailed      = "SNAPSHOT_FAILED"
	ErrCodeValidationFailed    = "VALIDATION_FAILED"

	// External Services / Dependencies
	ErrCodeExternalServiceError = "EXTERNAL_SERVICE_ERROR"
	ErrCodeIntegrationFailed    = "INTEGRATION_FAILED"
	ErrCodeWebhookFailed        = "WEBHOOK_FAILED"
	ErrCodeNotificationFailed   = "NOTIFICATION_FAILED"
	ErrCodeThirdPartyAPIError   = "THIRD_PARTY_API_ERROR"
	ErrCodeAuthProviderError    = "AUTH_PROVIDER_ERROR"

	// System/Infrastructure
	ErrCodeInternal            = "INTERNAL_ERROR"
	ErrCodeNotImplemented      = "NOT_IMPLEMENTED"
	ErrCodeUnavailable         = "UNAVAILABLE"
	ErrCodeMaintenanceMode     = "MAINTENANCE_MODE"
	ErrCodeServiceShuttingDown = "SERVICE_SHUTTING_DOWN"
)
