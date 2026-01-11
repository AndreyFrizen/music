package handlers

import (
	"encoding/json"
	"mess/internal/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func Test_ArtistHandler(t *testing.T) {
	router := gin.New()
	router.POST("/addartist", func(c *gin.Context) {
		var artist model.Artist
		if err := c.ShouldBindJSON(&artist); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, artist)
	})

	w := httptest.NewRecorder()

	// Create an example user for testing
	exampleUser := model.Artist{
		Name: "test",
	}
	userJson, _ := json.Marshal(exampleUser)
	req, _ := http.NewRequest("POST", "/addartist", strings.NewReader(string(userJson)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
