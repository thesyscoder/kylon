package routes

import "github.com/gin-gonic/gin"

func InitializeHealthRoutes(rg *gin.RouterGroup) {
	rg.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server is running",
		})
	})
}
