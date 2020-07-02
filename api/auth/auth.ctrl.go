package auth

import (
	"contree/models"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func authMiddleware(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if len(auth) < 7 {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	credentials, err := base64.StdEncoding.DecodeString(auth[6:])
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	lists := strings.Split(string(credentials), ":")
	if len(lists) != 2 {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	var u *models.User
	u, err = models.CheckIfUserExistsAndHasGoodPassword(lists[0], lists[1])

	c.Set("User", u)
	c.Next()
}

func login(c *gin.Context) {
	user := c.MustGet("User").(*models.User)
	c.JSON(http.StatusOK, user)
}
