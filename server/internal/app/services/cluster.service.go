/**
 * @File: cluster_service.go
 * @Title: Cluster Service
 * @Description: Provides business logic for managing Cluster entities,
 * @Description: including registration and validation.
 */

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
	RegisterCluster(ctx context.Context, name, kubeconfig string) (*models.Cluster, error)
}

// ClusterServiceImpl provides the implementation of ClusterService.
type ClusterServiceImpl struct {
	clusterRepo repositories.ClusterRepository
	log         *logrus.Logger
}

// NewClusterService creates a new ClusterServiceImpl.
func NewClusterService(clusterRepo repositories.ClusterRepository, log *logrus.Logger) ClusterService {
	return &ClusterServiceImpl{
		clusterRepo: clusterRepo,
		log:         log,
	}
}

// RegisterCluster validates input and creates a new cluster record.
func (s *ClusterServiceImpl) RegisterCluster(ctx context.Context, name, kubeconfig string) (*models.Cluster, error) {
	s.log.WithField("cluster_name", name).Info("Attempting to register new cluster.")

	if strings.TrimSpace(name) == "" || strings.TrimSpace(kubeconfig) == "" {
		s.log.Warn("Invalid cluster registration input: name or kubeconfig is empty.")
		return nil, customerrors.NewCustomError(
			customerrors.ErrCodeInvalidInput,
			"Cluster name and kubeconfig are required.",
			nil,
			http.StatusBadRequest,
			nil,
		)
	}

	cluster := &models.Cluster{
		Name:       name,
		Kubeconfig: kubeconfig,
	}

	if err := s.clusterRepo.Create(ctx, cluster); err != nil {
		s.log.WithError(err).Error("Failed to persist cluster via repository.")
		// The error returned from the repository is already a custom error.
		return nil, err
	}

	s.log.WithField("cluster_id", cluster.ID).Info("Cluster registered successfully.")
	return cluster, nil
}
