package handlers

import (
	"database/sql"
	"fmt"
	"mess/internal/model"
	auth "mess/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register User in app.
func RegisterUser(c *gin.Context, store model.UserRepository) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := auth.UserService(user, store)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Auth ...
func AuthUser(c *gin.Context, db *sql.DB) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := AuthenticateUser(db, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auth user successfully"})
}

func AuthenticateUser(db *sql.DB, u *model.User) error {
	query := fmt.Sprintf("SELECT password FROM users WHERE email = '%s'", u.Email)

	rows, err := db.Query(query)
	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		var Password string

		if err := rows.Scan(&Password); err != nil {
			return err
		}

		if err := bcrypt.CompareHashAndPassword([]byte(Password), []byte(u.Password)); err != nil {
			return err
		}

		return nil
	}

	return nil
}

// Add Artist ...
func AddArtist(c *gin.Context, store model.ArtistRepository) error {
	var artist model.Artist
	if err := c.BindJSON(&artist); err != nil {
		return err
	}

	if err := auth.ArtistService(artist, store); err != nil {
		return err
	}

	return nil
}

// Add Track ...
func AddTrack(c *gin.Context, store model.TrackRepository) error {
	var track model.Track
	if err := c.BindJSON(&track); err != nil {
		return err
	}

	if err := auth.TrackService(track, store); err != nil {
		return err
	}

	return nil
}

// Add Album ...
func AddAlbum(c *gin.Context, store model.AlbumRepository) error {
	var album model.Album
	if err := c.BindJSON(&album); err != nil {
		return err
	}

	if err := auth.AlbumService(album, store); err != nil {
		return err
	}

	return nil
}

// Add Playlist ...
func CreatePlaylist(c *gin.Context, store model.PlaylistRepository) error {
	var playlist model.Playlist
	if err := c.BindJSON(&playlist); err != nil {
		return err
	}

	if err := auth.PlaylistService(playlist, store); err != nil {
		return err
	}

	return nil
}

// Add track to playlist ...
func AddTrackToPlaylist(c *gin.Context, store model.PlaylistRepository) error {
	var trackToPlaylist model.PlaylistTrack
	if err := c.BindJSON(&trackToPlaylist); err != nil {
		return err
	}

	if err := auth.PlaylistTrackService(trackToPlaylist, store); err != nil {
		return err
	}

	return nil
}
