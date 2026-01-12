package handlers

import (
	"mess/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type playlistHandler interface {
	PlaylistByID(c *gin.Context)
	CreatePlaylist(c *gin.Context)
	DeletePlaylist(c *gin.Context)
}

// Create playlist
func (h *Handler) CreatePlaylist(c *gin.Context) {
	var playlist model.Playlist
	if err := c.BindJSON(&playlist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreatePlaylist(&playlist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, playlist)
}

// Delete playlist
func (h *Handler) DeletePlaylist(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeletePlaylist(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Playlist deleted"})
}

// Get playlist
func (h *Handler) PlaylistByID(c *gin.Context) {
	id := c.Param("id")
	playlist, err := h.service.PlaylistByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, playlist)
}
