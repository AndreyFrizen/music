package middl

import (
	"errors"
	services "mess/internal/service"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey = []byte("123")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		id, err := parseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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
