/**
 * @File: health_handler.go
 * @Title: Health Check Handler
 * @Description: Provides an HTTP handler for performing application health checks,
 * @Description: including detailed status of core dependencies like database and Kubernetes.
 * @Author: thesyscoder (github.com/thesyscoder)
 */

package handlers

import (
	"context" // For context with timeout
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thesyscoder/kylon/internal/app/utils"             // For SuccessResponse and ErrorResponse
	"github.com/thesyscoder/kylon/internal/infrastructure/config" // For application configuration
	customerrors "github.com/thesyscoder/kylon/pkg/customErrors"  // For structured custom errors
	"github.com/thesyscoder/kylon/pkg/logger"                     // For centralized logging
	"gorm.io/gorm"                                                // For database interaction
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"                 // For Kubernetes API options
	"k8s.io/client-go/kubernetes"                                 // For Kubernetes client interaction
)

// log is the logger instance for this package, providing contextual logging for the health handler.
var log = logger.GetLogger().WithField("component", "health_handler")

// HealthCheckResponse defines the structure for a detailed health check API response.
type HealthCheckResponse struct {
	Status       string `json:"status"`                // Overall service status (e.g., "UP", "DOWN")
	Message      string `json:"message"`               // A human-readable message about the service status
	Application  string `json:"application"`           // Name of the application
	Version      string `json:"version,omitempty"`     // Application version from config
	Environment  string `json:"environment,omitempty"` // Application environment from config
	Timestamp    string `json:"timestamp"`             // Timestamp of when the health check was performed
	Dependencies struct {
		Database   string `json:"database"`   // Status of the database dependency
		Kubernetes string `json:"kubernetes"` // Status of the Kubernetes client dependency
		// Add more dependencies here as the application grows, e.g., "minio": "UP"
	} `json:"dependencies"`
}

// HealthCheckHandler returns a gin.HandlerFunc that performs a comprehensive health check.
// It assesses the overall application status and the connectivity/status of key dependencies,
// such as the configured database and Kubernetes cluster.
// Parameters:
//   - cfg: Application configuration, used for fetching app details (name, version, environment).
//   - db: GORM database instance; used to check database connectivity. Can be nil to skip check.
//   - kubeClient: Kubernetes Clientset instance; used to check Kubernetes API connectivity. Can be nil to skip check.
func HealthCheckHandler(cfg *config.Config, db *gorm.DB, kubeClient *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize the response with default "UP" status.
		response := HealthCheckResponse{
			Status:      http.StatusText(http.StatusOK), // Start with "OK" status text
			Message:     "Service is healthy",
			Application: cfg.App.Name,
			Version:     cfg.App.Version,
			Environment: cfg.App.Env,
			Timestamp:   time.Now().Format(time.RFC3339),
			Dependencies: struct {
				Database   string `json:"database"`
				Kubernetes string `json:"kubernetes"`
			}{
				Database:   "N/A", // Default to "N/A" if dependency is not checked or configured
				Kubernetes: "N/A",
			},
		}

		// Initialize overall status to OK. This will be degraded if any dependency fails.
		overallHTTPStatus := http.StatusOK
		overallMessage := "Service is healthy"

		// --- Database Health Check ---
		if db != nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second) // 2-second timeout for DB ping
			defer cancel()

			sqlDB, err := db.DB() // Get the underlying *sql.DB for pinging
			if err != nil {
				log.WithError(err).Error("Health Check: Failed to retrieve underlying DB connection pool.")
				response.Dependencies.Database = "DOWN - Connection Pool Error"
				overallHTTPStatus = http.StatusInternalServerError
				overallMessage = "Service unhealthy: Database connection pool issue"
			} else if err := sqlDB.PingContext(ctx); err != nil {
				log.WithError(err).Error("Health Check: Database connectivity check failed.")
				response.Dependencies.Database = fmt.Sprintf("DOWN - %s", err.Error())
				overallHTTPStatus = http.StatusInternalServerError
				overallMessage = "Service unhealthy: Database connectivity failed"
			} else {
				response.Dependencies.Database = "UP"
			}
		} else {
			log.Debug("Health Check: Database instance not provided (nil), skipping DB health check.")
		}

		// --- Kubernetes Client Health Check ---
		if kubeClient != nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second) // 3-second timeout for K8s API call
			defer cancel()

			// Perform a lightweight Kubernetes API call, like listing namespaces, to verify connectivity.
			_, err := kubeClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
			if err != nil {
				log.WithError(err).Error("Health Check: Kubernetes API connectivity check failed.")
				response.Dependencies.Kubernetes = fmt.Sprintf("DOWN - %s", err.Error())
				// Only set overall status to error if it's not already an error from previous checks
				if overallHTTPStatus == http.StatusOK {
					overallHTTPStatus = http.StatusInternalServerError
					overallMessage = "Service unhealthy: Kubernetes API communication failed"
				}
			} else {
				response.Dependencies.Kubernetes = "UP"
			}
		} else {
			log.Debug("Health Check: Kubernetes client not provided (nil), skipping K8s health check.")
		}

		// Update the overall status and message in the response based on dependency checks.
		response.Status = http.StatusText(overallHTTPStatus)
		response.Message = overallMessage

		// Send the final API response. If overall status is not OK, use ErrorResponse.
		if overallHTTPStatus != http.StatusOK {
			// For an unhealthy status, use the common error response utility,
			// mapping to an internal error code and message.
			utils.ErrorResponse(c, customerrors.NewCustomError(
				customerrors.ErrCodeInternal, // Application-specific internal error code
				overallMessage,               // The summarized message for the client
				nil,                          // No underlying specific error to wrap for this summary response
				overallHTTPStatus,            // The HTTP status code determined by health checks
				response,                     // Pass the full health check response as data in the error payload
			))
		} else {
			// For a healthy status, use the success response utility.
			utils.SuccessResponse(c, http.StatusOK, response.Message, response)
		}
	}
}
