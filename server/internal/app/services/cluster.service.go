// services/cluster_service.go
package services

import (
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/thesyscoder/kylon/internal/domain/models"
	"github.com/thesyscoder/kylon/internal/infrastructure/repositories"
	customerrors "github.com/thesyscoder/kylon/pkg/customErrors"
)

// ClusterService defines the interface for cluster-related business operations.
type ClusterService interface {
	// Kubeconfig parameter now represents the *content* of the kubeconfig file, not its path.
	RegisterCluster(ctx context.Context, name, kubeconfigContent string) (*models.Cluster, error)
	ListClusters(ctx context.Context) ([]models.Cluster, error)
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
		return nil, err
	}

	s.log.WithContext(ctx).WithField("cluster_id", cluster.ID).Info("Cluster registered successfully.")
	return cluster, nil
}

// ListClusters remains unchanged as its logic does not depend on file upload.
func (s *ClusterServiceImpl) ListClusters(ctx context.Context) ([]models.Cluster, error) {
	s.log.WithContext(ctx).Info("Attempting to list all clusters.")

	clusters, err := s.clusterRepo.List(ctx)
	if err != nil {
		s.log.WithContext(ctx).WithError(err).Error("Failed to retrieve clusters from repository.")
		return nil, err
	}

	s.log.WithContext(ctx).Infof("Successfully retrieved %d clusters.", len(clusters))
	return clusters, nil
}
