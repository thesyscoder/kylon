// services/cluster_service.go
package services

import (
	"context"
	"net/http"
	"strings"
	"time" // Keep time imported as it's used for formatting

	"github.com/sirupsen/logrus"
	"github.com/thesyscoder/kylon/internal/domain/models"
	"github.com/thesyscoder/kylon/internal/domain/types" // IMPORTANT: Keep this import
	"github.com/thesyscoder/kylon/internal/infrastructure/repositories"
	customerrors "github.com/thesyscoder/kylon/pkg/customErrors"
)

// ClusterService defines the interface for cluster-related business operations.
type ClusterService interface {
	// Kubeconfig parameter now represents the *content* of the kubeconfig file, not its path.
	RegisterCluster(ctx context.Context, name, kubeconfigContent string) (*models.Cluster, error)
	// <--- UPDATED: Interface method now returns []types.ClusterSummary
	ListClusters(ctx context.Context) ([]types.ClusterSummary, error)
}

// ClusterServiceImpl provides the implementation of ClusterService.
type ClusterServiceImpl struct {
	clusterRepo repositories.ClusterRepository
	log         *logrus.Logger
}

// NewClusterService creates a new ClusterServiceImpl.
func NewClusterService(clusterRepo repositories.ClusterRepository, log *logrus.Logger) ClusterService {
	if clusterRepo == nil {
		log.Fatal("ClusterRepository is nil when creating ClusterService. This indicates a critical setup error.")
	}
	return &ClusterServiceImpl{
		clusterRepo: clusterRepo,
		log:         log,
	}
}

// RegisterCluster validates input and creates a new cluster record.
// kubeconfigContent is the actual content read from the uploaded file.
func (s *ClusterServiceImpl) RegisterCluster(ctx context.Context, name, kubeconfigContent string) (*models.Cluster, error) {
	s.log.WithContext(ctx).WithField("cluster_name", name).Info("Attempting to register new cluster with uploaded kubeconfig.")

	// Input validation for name and kubeconfig content
	if strings.TrimSpace(name) == "" || strings.TrimSpace(kubeconfigContent) == "" {
		s.log.WithContext(ctx).Warn("Invalid cluster registration input: name or kubeconfig content is empty.")
		return nil, customerrors.NewCustomError(
			customerrors.ErrCodeInvalidInput,

			"Cluster name and kubeconfig content are required.",

			nil, // No underlying error

			http.StatusBadRequest,

			nil,
		)
	}

	// Create the domain model
	cluster := &models.Cluster{
		Name: name,
		// Assign the content to the Kubeconfig field in the model
		Kubeconfig: kubeconfigContent, // Assuming models.Cluster has Kubeconfig string
	}

	if err := s.clusterRepo.Create(ctx, cluster); err != nil {
		s.log.WithContext(ctx).WithError(err).Error("Failed to persist cluster via repository.")
		// The error from repo should ideally already be a custom error.
		// If not, wrap it into a relevant custom error.
		return nil, err
	}

	s.log.WithContext(ctx).WithField("cluster_id", cluster.ID).Info("Cluster registered successfully.")
	return cluster, nil
}

// ListClusters retrieves all cluster records from the repository and converts them to summary DTOs.
// <--- UPDATED: Implementation return type now matches the interface.
func (s *ClusterServiceImpl) ListClusters(ctx context.Context) ([]types.ClusterSummary, error) {
	s.log.WithContext(ctx).Info("Attempting to list all clusters.")

	clusters, err := s.clusterRepo.List(ctx)
	if err != nil {
		s.log.WithContext(ctx).WithError(err).Error("Failed to retrieve clusters from repository.")
		return nil, err // Assuming repository already returns custom errors, or wrap generic errors here.
	}

	summaries := make([]types.ClusterSummary, len(clusters))
	for i, cluster := range clusters {
		summaries[i] = types.ClusterSummary{

			ID:   cluster.ID.String(),
			Name: cluster.Name,

			CreatedAt: cluster.CreatedAt.Format(time.RFC3339),
			UpdatedAt: cluster.UpdatedAt.Format(time.RFC3339),
		}
	}

	s.log.WithContext(ctx).Infof("Successfully retrieved %d clusters and converted to summaries.", len(summaries))
	return summaries, nil
}
