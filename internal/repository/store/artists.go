package store

import (
	"fmt"
	"mess/internal/model"

	"github.com/google/uuid"
)

type artistRepository interface {
	CreateArtist(a *model.Artist) error
	ArtistByID(id string) (*model.Artist, error)
	Artists() ([]model.Artist, error)
}

// Post

// CreateArtist creates a new artist in the database
func (s *Store) CreateArtist(a *model.Artist) error {
	query := fmt.Sprintf("INSERT INTO artists VALUES ('%s', '%s')",
		uuid.New().String(), a.Name)

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

// GET

// ArtistByID retrieves an artist by their ID from the database
func (s *Store) ArtistByID(id string) (*model.Artist, error) {
	query := fmt.Sprintf("SELECT * FROM artists WHERE id = '%s'", id)

	row := s.db.QueryRow(query)

	var artist model.Artist

	err := row.Scan(&artist.Name, &artist.ID)

	if err != nil {
		return nil, err
	}

	return &artist, nil
}

func (s *Store) Artists() ([]model.Artist, error) {
	query := "SELECT name FROM artists"

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var artists []model.Artist

	for rows.Next() {
		var artist model.Artist

		err := rows.Scan(&artist.Name)

		if err != nil {
			return nil, err
		}

		artists = append(artists, artist)
	}

	return artists, nil
}
