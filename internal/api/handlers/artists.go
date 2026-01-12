package handlers

import (
	"mess/internal/model"
	templs "mess/static/templates"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=artists.go -destination=mocks/mock_artist_handler.go

type artistHandler interface {
	CreateArtist(c *gin.Context) error
	ArtistByID(c *gin.Context) error
	Artists(c *gin.Context) error
	ArtistsByName(c *gin.Context) error
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

func (h *Handler) Artists(c *gin.Context) error {
	artists, err := h.service.Artists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	return templs.ArtistsPage(artists).Render(c, c.Writer)
}

// ArtistsByName retrieves artists by name.
func (h *Handler) ArtistsByName(c *gin.Context) error {
	name := c.Param("name")

	_, err := h.service.ArtistsByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	return nil
}
