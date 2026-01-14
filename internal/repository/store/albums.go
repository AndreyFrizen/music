package store

import (
	"context"
	"encoding/json"
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
	query := fmt.Sprintf("INSERT INTO albums VALUES ('%s', '%s', '%s', '%s')",
		uuid.New().String(), a.Title, a.ArtistID, a.ReleaseDate)

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	jsonData, _ := json.MarshalIndent(a, "", "  ")
	s.cash.Set(ctx, a.ID.String(), string(jsonData), time.Minute*10)

	return nil
}

// Get

// AlbumByID retrieves an album by its ID from the database.
func (s *Store) AlbumByID(id string, ctx context.Context) (*model.Album, error) {

	var album model.Album

	albs, err := s.cash.Get(ctx, id).Bytes()
	if err == nil {
		err = json.Unmarshal(albs, &album)
	} else {
		query := fmt.Sprintf("SELECT * FROM albums WHERE id = '%s'", id)

		row := s.db.QueryRowContext(ctx, query)

		err = row.Scan(&album.Title, &album.ArtistID, &album.ReleaseDate, &album.ID)

		if err != nil {
			return nil, err
		}

		jsonData, _ := json.MarshalIndent(album, "", "  ")
		s.cash.Set(ctx, album.ID.String(), string(jsonData), time.Minute*10)
	}

	return &album, nil
}

// AlbumsByTitle retrieves all albums by their title from the database.
func (s *Store) AlbumsByTitle(title string, ctx context.Context) ([]model.Album, error) {
	var albums []model.Album

	albs, err := s.cash.Get(ctx, title).Bytes()
	if err == nil {
		err = json.Unmarshal(albs, &albums)
	} else {
		query := fmt.Sprintf("SELECT * FROM albums WHERE title = '%s'", title)

		rows, err := s.db.QueryContext(ctx, query)

		if err != nil {
			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			var album model.Album

			err := rows.Scan(&album.Title, &album.ArtistID, &album.ReleaseDate, &album.ID)

			if err != nil {
				return nil, err
			}

			albums = append(albums, album)
		}

		jsonData, _ := json.MarshalIndent(albums, "", "  ")
		s.cash.Set(ctx, title, string(jsonData), time.Minute*10)
	}

	return albums, nil
}

// AlbumsByArtistID retrieves all albums by their artist ID from the database.
func (s *Store) AlbumsByArtistID(artistID string, ctx context.Context) ([]model.Album, error) {
	var albums []model.Album

	albs, err := s.cash.Get(ctx, artistID).Bytes()
	if err == nil {
		err = json.Unmarshal(albs, &albums)
	} else {
		query := fmt.Sprintf("SELECT * FROM albums WHERE artist_id = '%s'", artistID)

		rows, err := s.db.QueryContext(ctx, query)

		if err != nil {
			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			var album model.Album

			err := rows.Scan(&album.Title, &album.ArtistID, &album.ReleaseDate, &album.ID)

			if err != nil {
				return nil, err
			}

			albums = append(albums, album)
		}

		jsonData, _ := json.MarshalIndent(albums, "", "  ")
		s.cash.Set(ctx, artistID, string(jsonData), time.Minute*10)
	}

	return albums, nil
}
