package api

import (
	"contree/api/auth"
	"contree/api/ping"

	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.Engine) {
	api := r.Group("/api")
	ping.ApplyRoutes(api)
	auth.ApplyRoutes(api)
}
