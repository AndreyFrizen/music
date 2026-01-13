package handlers

import (
	"mess/internal/model"
	"net/http"

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
	user, err := h.service.UserByID(id, c)
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
<<<<<<< HEAD
=======

// AuthMiddleware middleware for authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		id, err := parseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("userID", id)
		c.Next()
	}
}

func parseToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &services.TokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid method")
		}
		return secretKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*services.TokenClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	return claims.UserID, nil
}
>>>>>>> ed7a845 (v1)
