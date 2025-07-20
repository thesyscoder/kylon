package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thesyscoder/kylon/internal/app/utils"
)

type HealthCheckResposne struct {
	Status       string   `json:"status"`
	Message      string   `json:"message"`
	Application  string   `json:"application"`
	Version      string   `json:"version"`
	Environment  string   `json:"environment,omitempty"`
	Timestamp    string   `json:"timestamp"`
	Dependencies struct{} `json:"dependencies"`
}

func SetupHealthCheckHandler() gin.HandlerFunc {
	response := HealthCheckResposne{}
	return func(ctx *gin.Context) {
		utils.SuccessResponse(ctx, http.StatusOK, "Service is health.", response)
	}
}
