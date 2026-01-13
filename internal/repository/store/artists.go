package store

import (
	"context"
	"fmt"
	"mess/internal/model"

	"github.com/google/uuid"
)

type artistRepository interface {
	CreateArtist(a *model.Artist, ctx context.Context) error
	ArtistByID(id string, ctx context.Context) (*model.Artist, error)
	Artists(ctx context.Context) ([]model.Artist, error)
	ArtistsByName(name string, ctx context.Context) ([]model.Artist, error)
}

// Post

// CreateArtist creates a new artist in the database
func (s *Store) CreateArtist(a *model.Artist, ctx context.Context) error {
	query := fmt.Sprintf("INSERT INTO artists VALUES ('%s', '%s')",
		uuid.New().String(), a.Name)

	_, err := s.db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	return nil
}

// GET

// ArtistByID retrieves an artist by their ID from the database
func (s *Store) ArtistByID(id string, ctx context.Context) (*model.Artist, error) {
	query := fmt.Sprintf("SELECT * FROM artists WHERE id = '%s'", id)

	row := s.db.QueryRow(query)

	var artist model.Artist

	err := row.Scan(&artist.Name, &artist.ID)

	if err != nil {
		return nil, err
	}

	return &artist, nil
}

// Artists retrieves all artists from the database
func (s *Store) Artists(ctx context.Context) ([]model.Artist, error) {
	query := "SELECT name FROM artists"

	rows, err := s.db.QueryContext(ctx, query)

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

// ArtistsByName retrieves all artists by their name from the database
func (s *Store) ArtistsByName(name string, ctx context.Context) ([]model.Artist, error) {
	query := fmt.Sprintf("SELECT * FROM artists WHERE name = '%s'", name)

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var artists []model.Artist

	for rows.Next() {
		var artist model.Artist

		err := rows.Scan(&artist.Name, &artist.ID)

		if err != nil {
			return nil, err
		}

		artists = append(artists, artist)
	}

	return artists, nil
}
