package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thesyscoder/kylon/internal/app/handlers"
	"github.com/thesyscoder/kylon/internal/infrastructure/config"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
)

func InitializeRoutes(cfg *config.Config, db *gorm.DB, kubeClient *kubernetes.Clientset) *gin.Engine {
	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("healthz", handlers.HealthCheckHandler(cfg, db, kubeClient))
	}

	return router
}
