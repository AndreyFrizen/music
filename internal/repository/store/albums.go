package store

import (
	"fmt"
	"mess/internal/model"

	"github.com/google/uuid"
)

type albumRepository interface {
	AddAlbum(a *model.Album) error
	AlbumByID(id string) (*model.Album, error)
}

// Post

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

// Get

// AlbumByID retrieves an album by its ID from the database.
func (s *Store) AlbumByID(id string) (*model.Album, error) {
	query := fmt.Sprintf("SELECT * FROM albums WHERE id = '%s'", id)

	row := s.db.QueryRow(query)

	var album model.Album

	err := row.Scan(&album.Title, &album.ArtistID, &album.ReleaseDate, &album.ID)

	if err != nil {
		return nil, err
	}

	return &album, nil
}
