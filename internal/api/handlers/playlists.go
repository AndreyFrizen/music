package handlers

import (
	"mess/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type playlistHandler interface {
	PlaylistByID(c *gin.Context) error
	CreatePlaylist(c *gin.Context) error
	DeletePlaylist(c *gin.Context) error
}

// Create playlist
func (h *Handler) CreatePlaylist(c *gin.Context) error {
	var playlist model.Playlist
	if err := c.BindJSON(&playlist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	if err := h.service.CreatePlaylist(&playlist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}
	c.JSON(http.StatusCreated, playlist)
	return nil
}

// Delete playlist
func (h *Handler) DeletePlaylist(c *gin.Context) error {
	id := c.Param("id")
	if err := h.service.DeletePlaylist(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}
	c.JSON(http.StatusOK, gin.H{"message": "Playlist deleted"})
	return nil
}

// Get playlist
func (h *Handler) PlaylistByID(c *gin.Context) error {
	id := c.Param("id")
	playlist, err := h.service.PlaylistByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}
	c.JSON(http.StatusOK, playlist)
	return nil
}
