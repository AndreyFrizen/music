package authmiddleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Check if user is authenticated
		if !isAuthenticated(c) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

func isAuthenticated(c *gin.Context) bool {

	return true
}
