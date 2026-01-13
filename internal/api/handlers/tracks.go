package handlers

import (
	"mess/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type trackHandler interface {
	AddTrack(c *gin.Context)
	AddTrackToPlaylist(c *gin.Context)
	TrackByID(c *gin.Context)
	TrackFromPlaylist(c *gin.Context)
	DeleteTrackFromPlaylist(c *gin.Context)
	TracksByTitle(c *gin.Context)
	TracksByArtist(c *gin.Context)
}

// Add Track to Playlist
func (h *Handler) AddTrackToPlaylist(c *gin.Context) {
	var track model.PlaylistTrack
	if err := c.BindJSON(&track); err != nil {
		return
	}

	if err := h.service.AddTrackToPlaylist(&track, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Track added to playlist successfully"})
}

// Add Track to platform
func (h *Handler) AddTrack(c *gin.Context) {
	var track model.Track
	if err := c.BindJSON(&track); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddTrack(&track, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Track added successfully"})
}

// Delete Track from playlist
func (h *Handler) DeleteTrackFromPlaylist(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteTrackFromPlaylist(id, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Track deleted successfully"})
}

// Get Track by ID
func (h *Handler) TrackByID(c *gin.Context) {
	id := c.Param("id")
	track, err := h.service.TrackByID(id, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, track)

}

// Get Track from playlist
func (h *Handler) TrackFromPlaylist(c *gin.Context) {
	id := c.Param("id")
	track, err := h.service.TrackFromPlaylist(id, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, track)
}

// TracksByTitle retrieves tracks by title
func (h *Handler) TracksByTitle(c *gin.Context) {
	title := c.Param("title")

	tracks, err := h.service.TracksByTitle(title, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

// TracksByArtist retrieves tracks by artist
func (h *Handler) TracksByArtist(c *gin.Context) {
	artist := c.Param("artist")

	tracks, err := h.service.TracksByArtist(artist, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracks)
}
