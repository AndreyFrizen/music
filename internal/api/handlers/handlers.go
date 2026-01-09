package handlers

import (
	"mess/internal/model"
	services "mess/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handl struct {
	service services.MusicInterfaceService
}

type Handler interface {
	RegisterUser(c *gin.Context) error
	AddArtist(c *gin.Context) error
	AddTrack(c *gin.Context) error
	AddAlbum(c *gin.Context) error
	AddPlaylist(c *gin.Context) error
	AddTrackToPlaylist(c *gin.Context) error
}

func NewHandler(repo services.MusicInterfaceService) *Handl {
	return &Handl{
		service: repo,
	}
}

// Register User in app.
func (h *Handl) RegisterUser(c *gin.Context) error {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	err := h.service.RegisterService(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
	return nil
}

// Add Artist ...
func (h *Handl) AddArtist(c *gin.Context) error {
	var artist model.Artist
	if err := c.BindJSON(&artist); err != nil {
		return err
	}

	if err := h.service.ArtistService(&artist); err != nil {
		return err
	}

	return nil
}

// Add Track ...
func (h *Handl) AddTrack(c *gin.Context) error {
	var track model.Track
	if err := c.BindJSON(&track); err != nil {
		return err
	}

	if err := h.service.TrackService(&track); err != nil {
		return err
	}

	return nil
}

// Add Album ...
func (h *Handl) AddAlbum(c *gin.Context) error {
	var album model.Album
	if err := c.BindJSON(&album); err != nil {
		return err
	}

	if err := h.service.AlbumService(&album); err != nil {
		return err
	}

	return nil
}

// Add Playlist ...
func (h *Handl) AddPlaylist(c *gin.Context) error {
	var playlist model.Playlist
	if err := c.BindJSON(&playlist); err != nil {
		return err
	}

	if err := h.service.PlaylistService(&playlist); err != nil {
		return err
	}

	return nil
}

// Add track to playlist ...
func (h *Handl) AddTrackToPlaylist(c *gin.Context) error {
	var trackToPlaylist model.PlaylistTrack
	if err := c.BindJSON(&trackToPlaylist); err != nil {
		return err
	}

	if err := h.service.PlaylistTrackService(&trackToPlaylist); err != nil {
		return err
	}

	return nil
}

// Get track stream ...
// func (h *Handl) GetTrackStream(c *gin.Context) error {

// 	if err := auth.GetTrackStream(); err != nil {
// 		return err
// 	}

// 	return nil
// }
