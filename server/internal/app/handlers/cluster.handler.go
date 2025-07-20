// internal/app/handlers/cluster.handler.go
package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/thesyscoder/kylon/internal/app/services"
	"github.com/thesyscoder/kylon/internal/app/utils"
	customerrors "github.com/thesyscoder/kylon/pkg/customErrors"
)

// ClusterHandler handles HTTP requests related to cluster management.
type ClusterHandler struct {
	clusterService services.ClusterService
	log            *logrus.Logger
	// If you need direct access to config here, uncomment this:
	// cfg            *config.Config
}

// NewClusterHandler creates a new ClusterHandler.
func NewClusterHandler(clusterService services.ClusterService, log *logrus.Logger /*, cfg *config.Config */) *ClusterHandler {
	if clusterService == nil {
		log.Fatal("ClusterService is nil when creating ClusterHandler. This indicates a critical setup error.")
	}
	return &ClusterHandler{
		clusterService: clusterService,
		log:            log,
		// If you uncommented cfg above, also uncomment here:
		// cfg:            cfg,
	}
}

// RegisterCluster handles the registration of a new cluster via an HTTP POST request.
// It now expects a multipart/form-data request with 'name' and 'kubeconfig_file'.
func (h *ClusterHandler) RegisterCluster(c *gin.Context) {
	ctx := c.Request.Context() // Get context for logging and service calls
	h.log.WithContext(ctx).Info("Received request to register a new cluster via file upload.")

	// 1. Get the cluster name from form data
	name := c.PostForm("name")
	if strings.TrimSpace(name) == "" {
		h.log.WithContext(ctx).Warn("Invalid cluster registration request: 'name' form field is missing or empty.")
		utils.ErrorResponse(
			c,
			customerrors.NewCustomError(
				customerrors.ErrCodeInvalidInput,
				"Cluster name is required.",
				nil,
				http.StatusBadRequest,
				nil,
			),
		)
		return
	}

	// 2. Get the uploaded kubeconfig file
	// The variable is named 'file' to avoid the 'declared and not used' error if it was 'fileHeader'
	file, err := c.FormFile("kubeconfig_file") // 'kubeconfig_file' is the expected field name in the form
	if err != nil {
		h.log.WithContext(ctx).WithError(err).Warn("Failed to get uploaded kubeconfig file.")
		// Improved error message for missing file
		errToReturn := customerrors.NewCustomError(customerrors.ErrCodeInvalidInput, "Kubeconfig file 'kubeconfig_file' is required.", err, http.StatusBadRequest, nil)
		if err == http.ErrMissingFile {
			errToReturn = customerrors.NewCustomError(customerrors.ErrCodeInvalidInput, "Kubeconfig file is required (field name 'kubeconfig_file').", nil, http.StatusBadRequest, nil)
		}
		utils.ErrorResponse(c, errToReturn)
		return
	}

	// 3. Open the uploaded file
	src, err := file.Open() // Correctly uses 'file'
	if err != nil {
		h.log.WithContext(ctx).WithError(err).Error("Failed to open uploaded kubeconfig file.")
		utils.ErrorResponse(
			c,
			customerrors.NewCustomError(
				customerrors.ErrCodeInternal,
				"Failed to process uploaded file.",
				err,
				http.StatusInternalServerError,
				nil,
			),
		)
		return
	}
	defer src.Close() // Ensure the file is closed

	// 4. Read the content of the file
	kubeconfigContentBytes, err := io.ReadAll(src)
	if err != nil {
		h.log.WithContext(ctx).WithError(err).Error("Failed to read kubeconfig file content.")
		utils.ErrorResponse(
			c,
			customerrors.NewCustomError(
				customerrors.ErrCodeInternal,
				"Failed to read kubeconfig file content.",
				err,
				http.StatusInternalServerError,
				nil,
			),
		)
		return
	}
	kubeconfigContent := string(kubeconfigContentBytes)

	// 5. Delegate to the service layer with the content
	cluster, err := h.clusterService.RegisterCluster(ctx, name, kubeconfigContent)
	if err != nil {
		h.log.WithContext(ctx).WithError(err).Error("Failed to register cluster via service after file upload.")
		utils.ErrorResponse(c, err)
		return
	}

	h.log.WithContext(ctx).WithField("cluster_id", cluster.ID).Info("Cluster registered successfully via API with uploaded file.")
	utils.SuccessResponse(c, http.StatusCreated, "Cluster registered successfully.", cluster)
}

// ListClusters remains unchanged.
func (h *ClusterHandler) ListClusters(c *gin.Context) {
	h.log.WithContext(c.Request.Context()).Info("Received request to list all clusters.")

	clusters, err := h.clusterService.ListClusters(c.Request.Context())
	if err != nil {
		h.log.WithContext(c.Request.Context()).WithError(err).Error("Failed to list clusters via service.")
		utils.ErrorResponse(c, err)
		return
	}

	h.log.WithContext(c.Request.Context()).Infof("Successfully retrieved %d clusters for list request.", len(clusters))
	utils.SuccessResponse(c, http.StatusOK, "Clusters retrieved successfully.", clusters)
}
