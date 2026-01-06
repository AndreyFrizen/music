package store

import (
	"database/sql"
	"fmt"
	"mess/internal/model"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(u *model.User) error
	GetUserByID(id string) (*model.User, error)
	Authenticate(u *model.User) (*model.User, error)
	AddArtist(a *model.Artist) error
	AddTrack(t *model.Track) error
	AddAlbum(a *model.Album) error
	CreatePlaylist(p *model.Playlist) error
	AddTrackToPlaylist(p *model.PlaylistTrack) error
}

type Store struct {
	db *sql.DB
}

// Methods : CreateUser, GetUserByID, Authenticate, AddArtist, AddTrack,
// AddAlbum, CreatePlaylist, AddTrackToPlaylist
func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// CreateUser creates a new user in the database.
func (s *Store) CreateUser(u *model.User) error {
	query := fmt.Sprintf("INSERT INTO users VALUES ('%s', '%s', '%s', '%s')",
		uuid.New().String(), u.Username, u.EncryptedPassword, u.Email,
	)

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

// GetUserByID retrieves a user by their ID from the database.
func (s *Store) GetUserByID(id string) (*model.User, error) {
	query := fmt.Sprintf("SELECT * FROM users WHERE id = '%s'", id)

	row := s.db.QueryRow(query)

	var user model.User

	err := row.Scan(&user.ID, &user.Username, &user.EncryptedPassword, &user.Email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Authenticate a user by their email and password.
func (s *Store) Authenticate(u *model.User) (*model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", u.Email)

	row := s.db.QueryRow(query)

	err := row.Scan(&user.ID, &user.Username, &user.EncryptedPassword, &user.Email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Add Artist to database.
func (s *Store) AddArtist(a *model.Artist) error {
	query := fmt.Sprintf("INSERT INTO artists VALUES ('%s', '%s')",
		uuid.New().String(), a.Name)

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

// Add Track to database.
func (s *Store) AddTrack(t *model.Track) error {
	query := fmt.Sprintf("INSERT INTO tracks VALUES ('%s', '%s', '%s', '%s')",
		uuid.New().String(), t.Title, t.Duration, t.AudioURL)

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

// Add Album to database.
func (s *Store) AddAlbum(a *model.Album) error {
	query := fmt.Sprintf("INSERT INTO albums VALUES ('%s', '%s', '%s', '%s')",
		uuid.New().String(), a.Title, a.ArtistID, a.ReleaseDate)

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

// Add Playlist to database.
func (s *Store) CreatePlaylist(p *model.Playlist) error {
	query := fmt.Sprintf("INSERT INTO playlists VALUES ('%s', '%s', '%s')",
		uuid.New().String(), p.Title, p.UserID)

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

// Add Track to Playlist in database.
func (s *Store) AddTrackToPlaylist(t *model.PlaylistTrack) error {
	query := fmt.Sprintf("INSERT INTO track_to_playlist VALUES ('%s', '%s', '%d')",
		t.PlaylistID, t.TrackID, t.Position)

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}
