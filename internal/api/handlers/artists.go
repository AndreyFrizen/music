package handlers

import (
	"mess/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type artistHandler interface {
	CreateArtist(c *gin.Context) error
	ArtistByID(c *gin.Context) error
}

// ArtistByID retrieves an artist by ID.
func (h *Handler) ArtistByID(c *gin.Context) error {
	id := c.Param("id")

	artist, err := h.service.ArtistByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusOK, artist)
	return nil
}

// CreateArtist creates a new artist.
func (h *Handler) CreateArtist(c *gin.Context) error {
	var artist model.Artist

	if err := c.ShouldBindJSON(&artist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	err := h.service.CreateArtist(&artist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusCreated, gin.H{"message": "artist created"})
	return nil
}
