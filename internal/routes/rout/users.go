package routes

import (
	"mess/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, handler handlers.HandlerInterface) {
	user := r.Group("/user")
	{
		user.POST("/register", handler.RegisterUser)
		user.POST("/login", handler.LoginUser)
	}
}
