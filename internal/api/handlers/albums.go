package handlers

import (
	"mess/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type albumHandler interface {
	AlbumByID(c *gin.Context) error
	AddAlbum(c *gin.Context) error
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
