package handlers

import (
	"mess/internal/model"
	templs "mess/static/templates"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=artists.go -destination=mocks/mock_artist_handler.go

type artistHandler interface {
	CreateArtist(c *gin.Context)
	ArtistByID(c *gin.Context)
	Artists(c *gin.Context)
	ArtistsByName(c *gin.Context)
}

// ArtistByID retrieves an artist by ID.
func (h *Handler) ArtistByID(c *gin.Context) {
	id := c.Param("id")

	ids, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	artist, err := h.service.ArtistByID(ids, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, artist)
}

// CreateArtist creates a new artist.
func (h *Handler) CreateArtist(c *gin.Context) {
	var artist model.Artist

	if err := c.ShouldBindJSON(&artist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateArtist(&artist, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "artist created"})
}

func (h *Handler) Artists(c *gin.Context) {
	artists, err := h.service.Artists(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := templs.ArtistsPage(artists).Render(c, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// ArtistsByName retrieves artists by name.
func (h *Handler) ArtistsByName(c *gin.Context) {
	name := c.Param("name")

	artists, err := h.service.ArtistsByName(name, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "artists retrieved", "artists": artists})
}
