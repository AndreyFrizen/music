package handlers

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"track-service/internal/model"

	"github.com/gin-gonic/gin"
	"go.senan.xyz/taglib"
)

type trackHandler interface {
	TrackByID(c *gin.Context)
	TracksByTitle(c *gin.Context)
	TracksByArtist(c *gin.Context)
	UploadFile(c *gin.Context)
	Upload(c *gin.Context)
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

func (h *Handler) TracksByTitle(c *gin.Context) {
	title := c.Param("title")

	tracks, err := h.service.TracksByTitle(title, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

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
