package handlers

import (
	"fmt"
	"log"
	"mess/internal/model"
	templs "mess/static/templates"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//go:generate mockgen -source=artists.go -destination=mocks/mock_artist_handler.go

type artistHandler interface {
	CreateArtist(c *gin.Context)
	ArtistByID(c *gin.Context)
	Artists(c *gin.Context)
	ArtistsByName(c *gin.Context)
	FindArtists(c *gin.Context)
	ArtistWebSocket(c *gin.Context)
}

// ArtistByID retrieves an artist by ID.
func (h *Handler) ArtistByID(c *gin.Context) {
	id := c.Param("id")

	ids, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	artist, err := h.service.ArtistByID(ids, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, artist)
}

// CreateArtist creates a new artist.
func (h *Handler) CreateArtist(c *gin.Context) {
	var artist model.Artist

	if err := c.ShouldBindJSON(&artist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateArtist(&artist, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "artist created"})
}

func (h *Handler) Artists(c *gin.Context) {
	artists, err := h.service.Artists(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := templs.ArtistsPage(artists).Render(c, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// ArtistsByName retrieves artists by name.
func (h *Handler) ArtistsByName(c *gin.Context) {
	name := c.Param("name")

	artists, err := h.service.ArtistsByName(name, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "artists retrieved", "artists": artists})
}

// FindArtists retrieves artists by input string.
func (h *Handler) FindArtists(c *gin.Context) {
	c.HTML(http.StatusOK, "artists.html", gin.H{"message": "artists retrieved"})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) ArtistWebSocket(c *gin.Context) {
	log.Printf("WebSocket connection opened")
	// Upgrade the HTTP connection to a WebSocket connection.
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Handle WebSocket messages here
	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Received message: %s\n", data)

		// Echo the message back to the client
		if err := conn.WriteMessage(messageType, data); err != nil {
			fmt.Println(err)
			return
		}
	}
}
