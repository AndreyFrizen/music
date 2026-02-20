package middl

import (
	"errors"
	"mess/internal/model"
	services "mess/internal/service"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		id, err := parse(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("userID", id)
		c.Next()
	}
}

func parse(tokenString string) (string, error) {
	token, err := parseToken(tokenString)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*services.TokenClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	return claims.UserID, nil
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &services.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return model.SecretKey, nil
	})
}
