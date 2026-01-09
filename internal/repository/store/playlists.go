package store

import (
	"fmt"
	"mess/internal/model"

	"github.com/google/uuid"
)

type playlistRepository interface {
	CreatePlaylist(p *model.Playlist) error
	PlaylistByID(id string) (*model.Playlist, error)
	DeletePlaylist(id string) error
}

// Post

// Add Playlist to database
func (s *Store) CreatePlaylist(p *model.Playlist) error {
	query := fmt.Sprintf("INSERT INTO playlists VALUES ('%s', '%s', '%s')",
		uuid.New().String(), p.Title, p.UserID)

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

// Get

// PlaylistByID retrieves a playlist by its ID from the database.
func (s *Store) PlaylistByID(id string) (*model.Playlist, error) {
	query := fmt.Sprintf("SELECT * FROM playlists WHERE id = '%s'", id)

	row := s.db.QueryRow(query)

	var playlist model.Playlist

	err := row.Scan(&playlist.ID, &playlist.Title, &playlist.UserID)

	if err != nil {
		return nil, err
	}

	return &playlist, nil
}

// Delete

// Delete Playlist from database
func (s *Store) DeletePlaylist(id string) error {
	query := fmt.Sprintf("DELETE FROM playlists WHERE id = '%s'", id)

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}
