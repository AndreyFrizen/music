package handlers

import (
	"mess/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userHandler interface {
	RegisterUser(c *gin.Context)
	UserByID(c *gin.Context)
	LoginUser(c *gin.Context)
}

// Register User in app
func (h *Handler) RegisterUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Register(&user, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// UserByID retrieves a user by ID
func (h *Handler) UserByID(c *gin.Context) {
	id := c.Param("id")
	ids, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	user, err := h.service.UserByID(ids, c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Login User in app
func (h *Handler) LoginUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Auth(&user, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "not token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "token": token})
}
