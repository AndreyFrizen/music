package handlers

import (
	"log"
	"mess/internal/model"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.senan.xyz/taglib"
)

type trackHandler interface {
	AddTrackToPlaylist(c *gin.Context)
	TrackByID(c *gin.Context)
	TrackFromPlaylist(c *gin.Context)
	DeleteTrackFromPlaylist(c *gin.Context)
	TracksByTitle(c *gin.Context)
	TracksByArtist(c *gin.Context)
	Stream(c *gin.Context)
	UploadFile(c *gin.Context)
	Upload(c *gin.Context)
}

// Add Track to Playlist
func (h *Handler) AddTrackToPlaylist(c *gin.Context) {
	var track model.PlaylistTrack
	if err := c.BindJSON(&track); err != nil {
		return
	}
	log.Printf("Adding track %v", track)
	if err := h.service.AddTrackToPlaylist(&track, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Track added to playlist successfully"})
}

// Delete Track from playlist
func (h *Handler) DeleteTrackFromPlaylist(c *gin.Context) {
	id := c.Param("id")
	ids, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.service.DeleteTrackFromPlaylist(ids, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Track deleted successfully"})
}

// Get Track by ID
func (h *Handler) TrackByID(c *gin.Context) {
	id := c.Param("id")
	ids, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	track, err := h.service.TrackByID(ids, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, track)
}

// Get Track from playlist
func (h *Handler) TrackFromPlaylist(c *gin.Context) {
	id := c.Param("id")
	ids, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	track, err := h.service.TrackFromPlaylist(ids, c)
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
	artist := c.Param("id")
	id, err := strconv.Atoi(artist)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	tracks, err := h.service.TracksByArtist(id, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

// Stream track
func (h *Handler) Stream(c *gin.Context) {
	ids := c.Param("id")
	log.Print(ids)
	id, err := strconv.Atoi(ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	track, err := h.service.TrackByID(id, c)
	log.Print(track)
	log.Print(track.AudioURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err = os.Stat(track.AudioURL); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	file, err := os.Open(track.AudioURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()

	contentType := "audio/mpeg"

	c.Header("Content-Type", contentType)
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	c.Header("Accept-Ranges", "bytes")

	// Потоковая передача всего файла
	http.ServeContent(c.Writer, c.Request, fileInfo.Name(), fileInfo.ModTime(), file)
}

// html
func (h *Handler) Upload(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"message": "biba"})
}

// Upload track
func (h *Handler) UploadFile(c *gin.Context) {

	c.Request.ParseMultipartForm(10 << 20)

	file, header, err := c.Request.FormFile("myFile")
	if err != nil {
		return
	}
	defer file.Close()

	buffer := make([]byte, header.Size)
	file.Read(buffer)
	var b strings.Builder
	b.WriteString("/home/andrey/projects/music/static/music/")
	b.WriteString(header.Filename)
	data := b.String()
	err = os.WriteFile(data, buffer, 0644)
	prop, err := taglib.ReadProperties(data)
	if err != nil {
		return
	}
	timeTrack := int(prop.Length.Seconds())
	track := model.Track{
		Title:    header.Filename,
		Duration: timeTrack,
		AudioURL: data,
	}
	log.Print(track)
	h.service.AddTrack(&track, c)
}
