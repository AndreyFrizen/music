package handlers

import (
	"net/http"
	"playlist-service/internal/domain/model"
	"playlist-service/proto/playlist"
	"strconv"

	"github.com/gin-gonic/gin"
)

type handler struct {
	playlist.UnimplementedPlaylistServiceServer
	service PlaylistService
}

type PlaylistService interface {
	CreatePlaylist(*model.Playlist, *gin.Context) error
	DeletePlaylist(int64, *gin.Context) error
	PlaylistByID(int64, *gin.Context) (*model.Playlist, error)
	UpdatePlaylist(int64, *model.Playlist, *gin.Context) error
	AddTrack(int64, int64, *gin.Context) error
	RemoveTrack(int64, int64, *gin.Context) error
}

// Create playlist
func (h *handler) CreatePlaylist(c *gin.Context) {
	var playlist model.Playlist
	if err := c.BindJSON(&playlist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreatePlaylist(&playlist, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, playlist)
}

// Delete playlist
func (h *handler) DeletePlaylist(c *gin.Context) {
	id := c.Param("id")
	ids, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.service.DeletePlaylist(ids, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Playlist deleted"})
}

// Get playlist
func (h *Handler) PlaylistByID(c *gin.Context) {
	id := c.Param("id")
	ids, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	playlist, err := h.service.PlaylistByID(ids, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, playlist)
}
