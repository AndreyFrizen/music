package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var SecretKey = []byte("123")

// User represents a user in the system
type User struct {
	ID                uuid.UUID `json:"id" db:"id" redis:"id"`
	Username          string    `json:"username" db:"username" validate:"required,min=2,max=50" redis:"username"`
	Password          string    `json:"password" validate:"required,min=6,max=100"`
	EncryptedPassword string    `db:"password"`
	Email             string    `json:"email" db:"email" validate:"required,email" redis:"email"`
}

// Artist represents an artist in the system
type Artist struct {
	ID   uuid.UUID `json:"id" db:"id" redis:"id"`
	Name string    `json:"name" db:"name" redis:"name"`
}

// Album represents an album in the system
type Album struct {
	ID          uuid.UUID `json:"id" db:"id" redis:"id"`
	Title       string    `djson:"title" db:"title" redis:"title"`
	ArtistID    uuid.UUID `json:"artist_id" db:"artist_id" redis:"artist_id"`
	ReleaseDate time.Time `json:"release_date" db:"release_date" redis:"release_date"`
}

// Track represents a track in the system
type Track struct {
	ID       uuid.UUID `json:"id" db:"id" redis:"id"`
	Title    string    `json:"title" db:"title" redis:"title"`
	Duration time.Time `json:"duration" db:"duration" redis:"duration"`
	AudioURL string    `json:"audio_url" db:"audio_url" redis:"audio_url"`
	ArtistID uuid.UUID `json:"artist_id" db:"artist_id" redis:"artist_id"`
}

// Playlist represents a playlist in the system
type Playlist struct {
	ID     uuid.UUID `json:"id" db:"id" redis:"id"`
	UserID uuid.UUID `json:"user_id" db:"user_id" redis:"user_id"`
	Title  string    `json:"title" db:"title" redis:"title"`
}

// PlaylistTrack represents a track in a playlist in the system
type PlaylistTrack struct {
	PlaylistID uuid.UUID `json:"playlist_id" db:"playlist_id" redis:"playlist_id"`
	TrackID    uuid.UUID `json:"track_id" db:"track_id" redis:"track_id"`
	Position   int       `json:"position" db:"position" redis:"position"`
}

// AlbumTracks represents a collection of tracks in an album
type AlbumTracks struct {
	AlbumID uuid.UUID `json:"album_id" db:"album_id" redis:"album_id"`
	Tracks  []Track   `json:"tracks" db:"tracks" redis:"tracks"`
}

// ArtistAlbums represents a collection of albums by an artist
type ArtistAlbums struct {
	ArtistID uuid.UUID `json:"artist_id" db:"artist_id" redis:"artist_id"`
	Albums   []Album   `json:"albums" db:"albums" redis:"albums"`
}

// ArtistTracks represents a collection of tracks by an artist
type ArtistTracks struct {
	ArtistID uuid.UUID `json:"artist_id" db:"artist_id" redis:"artist_id"`
	Tracks   []Track   `json:"tracks" db:"tracks" redis:"tracks"`
}

// PlaylistTracks represents a collection of tracks in a playlist
type PlaylistTracks struct {
	PlaylistID uuid.UUID       `json:"playlist_id" db:"playlist_id" redis:"playlist_id"`
	Tracks     []PlaylistTrack `json:"tracks" db:"tracks" redis:"tracks"`
}

// UserAlbums represents a collection of albums by a user
type UserAlbums struct {
	UserID uuid.UUID `json:"user_id" db:"user_id" redis:"user_id"`
	Albums []Album   `json:"albums" db:"albums" redis:"albums"`
}

func (u *User) EncryptPassword() error {

	if len(u.Password) > 0 {
		enc, err := encryptedPassword(u.Password)
		if err != nil {
			return err
		}
		u.EncryptedPassword = string(enc)
	}

	return nil
}

func encryptedPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
