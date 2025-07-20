/**
 * @File: cluster.handler.go
 * @Title: Cluster Handler
 * @Description: Provides HTTP handlers for managing Cluster resources,
 * @Description: including registration.
 */

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/thesyscoder/kylon/internal/app/services"
	"github.com/thesyscoder/kylon/internal/app/utils"
	"github.com/thesyscoder/kylon/internal/domain/types"
	customerrors "github.com/thesyscoder/kylon/pkg/customErrors"
	// Use Kylon's custom logger
)

// ClusterHandler handles HTTP requests related to cluster management.
type ClusterHandler struct {
	clusterService services.ClusterService
	log            *logrus.Logger
}

// NewClusterHandler creates a new ClusterHandler.
func NewClusterHandler(clusterService services.ClusterService, log *logrus.Logger) *ClusterHandler {
	return &ClusterHandler{
		clusterService: clusterService,
		log:            log,
	}
}

// RegisterCluster handles the registration of a new cluster via an HTTP POST request.
func (h *ClusterHandler) RegisterCluster(c *gin.Context) {
	var req types.RegisterClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Warn("Invalid cluster registration request payload.")
		utils.ErrorResponse(
			c,
			customerrors.NewCustomError(
				customerrors.ErrCodeInvalidInput,
				"Invalid request: Name and kubeconfig are required and must be valid JSON fields.",
				err,
				http.StatusBadRequest,
				nil,
			),
		)
		return
	}

	cluster, err := h.clusterService.RegisterCluster(c.Request.Context(), req.Name, req.Kubeconfig)
	if err != nil {
		// The service layer is expected to return a custom error, which ErrorResponse can handle.
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Cluster registered successfully.", cluster)
}
