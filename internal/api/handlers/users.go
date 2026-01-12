package handlers

import (
	"mess/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler interface {
	RegisterUser(c *gin.Context) error
	UserByID(c *gin.Context) error
	AuthMiddleware() gin.HandlerFunc
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

type inputUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var inp inputUser
		if err := c.ShouldBindJSON(&inp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if err := Authenticate(h, inp.Email, inp.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func Authenticate(h *Handler, email, password string) error {
	err := h.service.Auth(email, password)
	if err != nil {
		return err
	}

	return nil
}
