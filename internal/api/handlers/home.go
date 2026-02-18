package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type homeHandler interface {
	Home(c *gin.Context)
}

func (h *Handler) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{"message": "home"})
}
