package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/thesyscoder/kylon/internal/app/handlers"
	"github.com/thesyscoder/kylon/internal/app/services"
	"github.com/thesyscoder/kylon/internal/infrastructure/config"
	"github.com/thesyscoder/kylon/internal/infrastructure/repositories"
	"gorm.io/gorm"
)

func RegisterClusterRoutes(rg *gin.RouterGroup, cfg *config.Config, db *gorm.DB, log *logrus.Logger) {
	clusters := rg.Group("/clusters")
	{
		// call the dependecies
		clusterRepo := repositories.NewClusterRepository(db, log)
		clusterService := services.NewClusterService(clusterRepo, log)
		clusterHandler := handlers.NewClusterHandler(clusterService, log)

		clusters.POST("", clusterHandler.RegisterCluster)

	}
}
