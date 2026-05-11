package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhasnanr/ewallet-ums/helpers"
)

type HealthCheck struct{}

func NewHealthCheck() *HealthCheck {
	return &HealthCheck{}
}

func (h *HealthCheck) RegisterRoute(r *gin.Engine) {
	healthAPI := r.Group("/health-check")
	healthAPI.GET("/", h.checkHealth)
}

func (h *HealthCheck) checkHealth(c *gin.Context) {
	helpers.SendResponseHTTP(c, http.StatusOK, "Server ums is healthy", nil)
}
