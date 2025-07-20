package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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
}

// NewClusterHandler creates a new ClusterHandler.
func NewClusterHandler(clusterService services.ClusterService, log *logrus.Logger) *ClusterHandler {
	if clusterService == nil {
		log.Fatal("ClusterService is nil when creating ClusterHandler. Critical setup error.")
	}
	return &ClusterHandler{
		clusterService: clusterService,
		log:            log,
	}
}

// RegisterCluster handles registration of a new cluster via file upload (multipart/form-data).
func (h *ClusterHandler) RegisterCluster(c *gin.Context) {
	ctx := c.Request.Context()
	h.log.WithContext(ctx).Info("Received request to register a new cluster via file upload.")

	name := c.PostForm("name")
	if strings.TrimSpace(name) == "" {
		h.log.WithContext(ctx).Warn("Missing or empty 'name' in cluster registration request.")
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

	fileHeader, err := c.FormFile("kubeconfig_file")
	if err != nil {
		h.log.WithContext(ctx).WithError(err).Warn("Failed to get uploaded kubeconfig file.")
		var errToReturn error
		if err == http.ErrMissingFile {
			errToReturn = customerrors.NewCustomError(
				customerrors.ErrCodeInvalidInput,
				"Kubeconfig file is required (field name 'kubeconfig_file').",
				nil,
				http.StatusBadRequest,
				nil,
			)
		} else {
			errToReturn = customerrors.NewCustomError(
				customerrors.ErrCodeInvalidInput,
				"Kubeconfig file 'kubeconfig_file' is required.",
				err,
				http.StatusBadRequest,
				nil,
			)
		}
		utils.ErrorResponse(c, errToReturn)
		return
	}

	src, err := fileHeader.Open()
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
	defer src.Close()

	// Prepare directory for kubeconfigs
	saveDir := "./data/kubeconfigs"
	if err := os.MkdirAll(saveDir, 0700); err != nil {
		utils.ErrorResponse(
			c, customerrors.NewCustomError(
				customerrors.ErrCodeInternal,
				"Failed to prepare kubeconfig storage.",
				err,
				http.StatusInternalServerError,
				nil,
			),
		)
		return
	}

	// Build a unique filename: <timestamp>_<cluster_name>.yaml (sanitized name for path safety)
	safeName := strings.ReplaceAll(strings.TrimSpace(name), " ", "_")
	timestamp := time.Now().Unix()
	filePath := filepath.Join(saveDir,
		filepath.Base(
			fmt.Sprintf("%d_%s.yaml", timestamp, safeName),
		),
	)

	outFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		utils.ErrorResponse(
			c, customerrors.NewCustomError(
				customerrors.ErrCodeInternal,
				"Failed to create kubeconfig file.",
				err,
				http.StatusInternalServerError,
				nil,
			),
		)
		return
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, src); err != nil {
		utils.ErrorResponse(
			c, customerrors.NewCustomError(
				customerrors.ErrCodeInternal,
				"Failed to save kubeconfig file.",
				err,
				http.StatusInternalServerError,
				nil,
			),
		)
		return
	}

	// Pass file path only to service (never the YAML content)
	cluster, err := h.clusterService.RegisterCluster(ctx, name, filePath)
	if err != nil {
		h.log.WithContext(ctx).WithError(err).Error("Failed to register cluster after file upload.")
		utils.ErrorResponse(c, err)
		return
	}

	h.log.WithContext(ctx).
		WithField("cluster_id", cluster.ID).
		WithField("kubeconfig_path", filePath).
		Info("Cluster registered successfully via API with uploaded file.")
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

	h.log.WithContext(c.Request.Context()).
		Infof("Successfully retrieved %d clusters for list request.", len(clusters))
	utils.SuccessResponse(c, http.StatusOK, "Clusters retrieved successfully.", clusters)
}
