/**
 * @File: response.go
 * @Title: API Response Utilities
 * @Description: Provides utility functions for consistent API response formatting,
 * @Description: including structured success and error responses for Gin contexts.
 * @Author: thesyscoder (github.com/thesyscoder)
 */

package utils

import (
	"errors"   // For standard error wrapping functions
	"net/http" // Import net/http for default status codes

	"github.com/gin-gonic/gin"
	customerrors "github.com/thesyscoder/kylon/pkg/customErrors"
	"github.com/thesyscoder/kylon/pkg/logger" // Import the custom logger package
)

// log is the logger instance for this package.
var log = logger.GetLogger().WithField("component", "api_response_utils")

// APIResponse represents the standardized API response structure.
type APIResponse struct {
	Success bool            `json:"success"`         // Indicates if the API call was successful
	Message string          `json:"message"`         // A human-readable message about the operation's result
	Data    any             `json:"data,omitempty"`  // Optional payload for successful responses
	Error   *APIErrorDetail `json:"error,omitempty"` // Optional detailed error information for failed responses
}

// APIErrorDetail provides detailed error information for API responses.
type APIErrorDetail struct {
	Code    string `json:"code"`           // Application-specific error code
	Message string `json:"message"`        // Specific error message
	Data    any    `json:"data,omitempty"` // Optional additional error-specific data
}

// SuccessResponse sends a standardized JSON success response to the client.
// `statusCode` is the HTTP status code (e.g., http.StatusOK).
// `message` is a descriptive success message.
// `data` is the optional payload to be returned.
func SuccessResponse(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
	log.Debugf("Success response sent with status %d: %s", statusCode, message)
}

// ErrorResponse sends a standardized JSON error response to the client.
// It automatically detects `*customerrors.CustomError` types and formats the response accordingly.
// For any other error type, it defaults to a generic internal server error.
func ErrorResponse(c *gin.Context, err error) {
	var customErr *customerrors.CustomError
	statusCode := http.StatusInternalServerError // Default HTTP status code for errors

	// Attempt to unwrap the error and check if it's a *customerrors.CustomError.
	if errors.As(err, &customErr) {
		// If it's a CustomError, use its defined HTTPStatus, Code, and Message.
		statusCode = customErr.HTTPStatus
		// Ensure code and message are present, fallbacks for safety
		if customErr.Code == "" {
			customErr.Code = customerrors.ErrCodeInternal
		}
		if customErr.Message == "" {
			customErr.Message = "An unspecified custom error occurred."
		}
		log.WithError(err).WithField("error_code", customErr.Code).
			Warnf("Sending custom error response: %s", customErr.Message)
	} else {
		// If the error is not a *customerrors.CustomError, it's an unhandled internal error.
		// Log it with full details for debugging purposes.
		log.WithError(err).Errorf("Unhandled error type in ErrorResponse: %T, Value: %v", err, err)

		// Create a generic internal server error CustomError for the response.
		customErr = customerrors.NewCustomError(
			customerrors.ErrCodeInternal,             // Application-specific internal error code
			"An unexpected internal error occurred.", // Generic message for the client
			err,                                      // Wrap the original unhandled error
			http.StatusInternalServerError,           // HTTP status code
			nil,                                      // No specific data for generic error
		)
	}

	// Send the JSON response with the determined HTTP status and error details.
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: customErr.Message, // Use the message from the custom error or the default generic message
		Error: &APIErrorDetail{
			Code:    customErr.Code,    // Use the error code from the custom error or default
			Message: customErr.Message, // Use the error message from the custom error or default
			Data:    customErr.Data,    // Include any additional data from the custom error
		},
	})
	log.Debugf("Error response sent with status %d: %s", statusCode, customErr.Message)
}
