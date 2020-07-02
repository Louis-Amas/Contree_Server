package ping

import (
	"github.com/gin-gonic/gin"
)

//ApplyRoutes Apply pings routes
func ApplyRoutes(r *gin.RouterGroup) {
	pings := r.Group("/ping")
	{
		pings.GET("/", ping)
	}
}
