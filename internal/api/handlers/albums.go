package handlers

import (
	"mess/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=albums.go -destination=mocks/mock_album_handler.go

type albumHandler interface {
	AlbumByID(c *gin.Context)
	AddAlbum(c *gin.Context)
	AlbumsByTitle(c *gin.Context)
	AlbumsByArtist(c *gin.Context)
}

// AddAlbum adds a new album to the platform
func (h *Handler) AddAlbum(c *gin.Context) {
	var album model.Album
	if err := c.BindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddAlbum(&album, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Album added successfully"})

}

// GetAlbumByID retrieves an album by its ID
func (h *Handler) AlbumByID(c *gin.Context) {
	id := c.Param("id")

	album, err := h.service.AlbumByID(id, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, album)
}

// AlbumsByTitle retrieves albums by title
func (h *Handler) AlbumsByTitle(c *gin.Context) {
	title := c.Param("title")

	albums, err := h.service.AlbumsByTitle(title, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, albums)
}

// AlbumsByArtist retrieves albums by artist
func (h *Handler) AlbumsByArtist(c *gin.Context) {
	artist := c.Param("artist")

	albums, err := h.service.AlbumsByArtist(artist, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, albums)
}
