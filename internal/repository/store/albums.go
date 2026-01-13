package store

import (
	"context"
	"fmt"
	"mess/internal/model"
	"time"

	"github.com/google/uuid"
)

type albumRepository interface {
	AddAlbum(a *model.Album, ctx context.Context) error
	AlbumByID(id string, ctx context.Context) (*model.Album, error)
	AlbumsByTitle(title string, ctx context.Context) ([]model.Album, error)
	AlbumsByArtistID(artistID string, ctx context.Context) ([]model.Album, error)
}

// Post

// Add Album to database.
func (s *Store) AddAlbum(a *model.Album, ctx context.Context) error {

	if err := s.cash.Get(ctx, a.ID.String()).Err(); err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO albums VALUES ('%s', '%s', '%s', '%s')",
		uuid.New().String(), a.Title, a.ArtistID, a.ReleaseDate)

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}

	s.cash.HSet(ctx, a.ID.String(), map[string]any{
		"title":        a.Title,
		"artist_id":    a.ArtistID.String(),
		"release_date": a.ReleaseDate.String(),
	}, time.Minute*20)

	return nil
}

// Get

// AlbumByID retrieves an album by its ID from the database.
func (s *Store) AlbumByID(id string, ctx context.Context) (*model.Album, error) {
	query := fmt.Sprintf("SELECT * FROM albums WHERE id = '%s'", id)

	row := s.db.QueryRowContext(ctx, query)

	var album model.Album

	err := row.Scan(&album.Title, &album.ArtistID, &album.ReleaseDate, &album.ID)

	if err != nil {
		return nil, err
	}

	return &album, nil
}

// AlbumsByTitle retrieves all albums by their title from the database.
func (s *Store) AlbumsByTitle(title string, ctx context.Context) ([]model.Album, error) {
	query := fmt.Sprintf("SELECT * FROM albums WHERE title = '%s'", title)

	rows, err := s.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var albums []model.Album

	for rows.Next() {
		var album model.Album

		err := rows.Scan(&album.Title, &album.ArtistID, &album.ReleaseDate, &album.ID)

		if err != nil {
			return nil, err
		}

		albums = append(albums, album)
	}

	return albums, nil
}

// AlbumsByArtistID retrieves all albums by their artist ID from the database.
func (s *Store) AlbumsByArtistID(artistID string, ctx context.Context) ([]model.Album, error) {
	query := fmt.Sprintf("SELECT * FROM albums WHERE artist_id = '%s'", artistID)

	rows, err := s.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var albums []model.Album

	for rows.Next() {
		var album model.Album

		err := rows.Scan(&album.Title, &album.ArtistID, &album.ReleaseDate, &album.ID)

		if err != nil {
			return nil, err
		}

		albums = append(albums, album)
	}

	return albums, nil
}
