/**
 * @File: cluster.repository.go
 * @Title: Cluster Repository
 * @Description: Defines the interface and implements the PostgreSQL repository
 * @Description: for managing Cluster entities.
 */

package repositories

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/thesyscoder/kylon/internal/domain/models"
	customerrors "github.com/thesyscoder/kylon/pkg/customErrors"
	"gorm.io/gorm"
)

// ClusterRepository defines operations for cluster data.
type ClusterRepository interface {
	Create(ctx context.Context, cluster *models.Cluster) error
}

// postgresClusterRepository is a PostgreSQL implementation of ClusterRepository.
type postgresClusterRepository struct {
	DB  *gorm.DB
	log *logrus.Logger // Use Kylon's custom logger
}

// NewClusterRepository creates a new postgresClusterRepository.
func NewClusterRepository(db *gorm.DB, log *logrus.Logger) ClusterRepository {
	return &postgresClusterRepository{
		DB:  db,
		log: log,
	}
}

// Create persists a new Cluster model to the database.
func (r *postgresClusterRepository) Create(ctx context.Context, cluster *models.Cluster) error {
	r.log.WithField("cluster_name", cluster.Name).Info("Creating new cluster record.")

	if err := r.DB.WithContext(ctx).Create(cluster).Error; err != nil {
		r.log.WithError(err).Error("Failed to create cluster record in database.")
		return customerrors.NewCustomError(
			customerrors.ErrCodeDatabaseOperationFailed,
			fmt.Sprintf("Failed to create cluster '%s'.", cluster.Name),
			err,
			http.StatusInternalServerError,
			nil,
		)
	}

	r.log.WithField("cluster_id", cluster.ID).Info("Successfully created cluster record.")
	return nil
}
