package auth

import (
	"github.com/gin-gonic/gin"
)

//ApplyRoutes Apply pings routes
func ApplyRoutes(r *gin.RouterGroup) {

	auth := r.Group("/login")
	{
		auth.GET("/", authMiddleware, login)
	}
}
