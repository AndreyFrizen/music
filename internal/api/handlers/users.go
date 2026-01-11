package handlers

import (
	"mess/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler interface {
	RegisterUser(c *gin.Context) error
	UserByID(c *gin.Context) error
}

// Register User in app
func (h *Handler) RegisterUser(c *gin.Context) error {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	err := h.service.Register(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
	return nil
}

// UserByID retrieves a user by ID
func (h *Handler) UserByID(c *gin.Context) error {
	id := c.Param("id")
	user, err := h.service.UserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusOK, user)
	return nil
}
