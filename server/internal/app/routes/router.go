package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/thesyscoder/kylon/internal/app/handlers"
	middleware "github.com/thesyscoder/kylon/internal/app/middlewares"
	"github.com/thesyscoder/kylon/internal/infrastructure/config"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
)

func InitializeRoutes(cfg *config.Config, db *gorm.DB, log *logrus.Logger, kubeClient *kubernetes.Clientset) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CorsMiddleware())
	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("healthz", handlers.SetupHealthCheckHandler())
	}
	// Clusters (call RegisterClusterRoutes)
	RegisterClusterRoutes(apiV1, cfg, db, log)
	return router
}
