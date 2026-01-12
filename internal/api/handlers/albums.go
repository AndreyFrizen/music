package handlers

import (
	"mess/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=albums.go -destination=mocks/mock_album_handler.go

type albumHandler interface {
	AlbumByID(c *gin.Context) error
	AddAlbum(c *gin.Context) error
	AlbumsByTitle(c *gin.Context) error
	AlbumsByArtist(c *gin.Context) error
}

// AddAlbum adds a new album to the platform
func (h *Handler) AddAlbum(c *gin.Context) error {
	var album model.Album
	if err := c.BindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	if err := h.service.AddAlbum(&album); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Album added successfully"})
	return nil
}

// GetAlbumByID retrieves an album by its ID
func (h *Handler) AlbumByID(c *gin.Context) error {
	id := c.Param("id")

	album, err := h.service.AlbumByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusOK, album)
	return nil
}

// AlbumsByTitle retrieves albums by title
func (h *Handler) AlbumsByTitle(c *gin.Context) error {
	title := c.Param("title")

	albums, err := h.service.AlbumsByTitle(title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusOK, albums)
	return nil
}

// AlbumsByArtist retrieves albums by artist
func (h *Handler) AlbumsByArtist(c *gin.Context) error {
	artist := c.Param("artist")

	albums, err := h.service.AlbumsByArtist(artist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusOK, albums)
	return nil
}
