package handlers

import (
	"mess/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type trackHandler interface {
	AddTrack(c *gin.Context) error
	AddTrackToPlaylist(c *gin.Context) error
	TrackByID(c *gin.Context) error
	TrackFromPlaylist(c *gin.Context) error
	DeleteTrackFromPlaylist(c *gin.Context) error
}

// Add Track to Playlist
func (h *Handler) AddTrackToPlaylist(c *gin.Context) error {
	var track model.PlaylistTrack
	if err := c.BindJSON(&track); err != nil {
		return err
	}

	if err := h.service.AddTrackToPlaylist(&track); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Track added to playlist successfully"})
	return nil
}

// Add Track to platform
func (h *Handler) AddTrack(c *gin.Context) error {
	var track model.Track
	if err := c.BindJSON(&track); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	if err := h.service.AddTrack(&track); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Track added successfully"})
	return nil
}

// Delete Track from playlist
func (h *Handler) DeleteTrackFromPlaylist(c *gin.Context) error {
	id := c.Param("id")
	if err := h.service.DeleteTrackFromPlaylist(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusOK, gin.H{"message": "Track deleted successfully"})
	return nil
}

// Get Track by ID
func (h *Handler) TrackByID(c *gin.Context) error {
	id := c.Param("id")
	track, err := h.service.TrackByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusOK, track)
	return nil
}

// Get Track from playlist
func (h *Handler) TrackFromPlaylist(c *gin.Context) error {
	id := c.Param("id")
	track, err := h.service.TrackFromPlaylist(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusOK, track)
	return nil
}
