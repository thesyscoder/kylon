package routes

import "github.com/gin-gonic/gin"

func InitializeRoutes() *gin.Engine {
	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		InitializeHealthRoutes(apiV1)
	}

	return router
}
