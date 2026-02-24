package store

import (
	"context"
	"encoding/json"
	"fmt"
	"mess/internal/model"
	"strconv"
	"time"
)

type artistRepository interface {
	CreateArtist(a *model.Artist, ctx context.Context) error
	ArtistByID(id int, ctx context.Context) (*model.Artist, error)
	Artists(ctx context.Context) ([]model.Artist, error)
	ArtistsByName(name string, ctx context.Context) ([]model.Artist, error)
	FindArtists(input string) (error, []model.Artist)
}

// Post

// CreateArtist creates a new artist in the database
func (s *Store) CreateArtist(a *model.Artist, ctx context.Context) error {
	query := fmt.Sprintf("INSERT INTO artists (name) VALUES ('%s')",
		a.Name)

	_, err := s.db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	jsonData, _ := json.MarshalIndent(a, "", "  ")
	s.cash.Set(ctx, strconv.Itoa(a.ID), string(jsonData), time.Minute*10)

	return nil
}

// GET

// ArtistByID retrieves an artist by their ID from the database
func (s *Store) ArtistByID(id int, ctx context.Context) (*model.Artist, error) {
	var artist model.Artist

	arts, err := s.cash.Get(ctx, strconv.Itoa(id)).Bytes()
	if err == nil {
		err = json.Unmarshal(arts, &artist)
	} else {
		query := fmt.Sprintf("SELECT * FROM artists WHERE id = '%v'", id)

		row := s.db.QueryRow(query)

		err := row.Scan(&artist.Name, &artist.ID)

		if err != nil {
			return nil, err
		}

		jsonData, _ := json.MarshalIndent(artist, "", "  ")
		s.cash.Set(ctx, strconv.Itoa(artist.ID), string(jsonData), time.Minute*10)
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

// FindArtists retrieves all artists by their name from the database
func (s *Store) FindArtists(input string) (error, []model.Artist) {

	var artists []model.Artist

	a := fmt.Sprintf("SELECT * FROM artists WHERE name = '%v'", input)

	rowsArtists, err := s.db.Query(a)
	if err != nil {
		return err, nil
	}
	defer rowsArtists.Close()

	for rowsArtists.Next() {
		var artist model.Artist

		err := rowsArtists.Scan(&artist.ID, &artist.Name)

		if err != nil {
			return err, nil
		}

		artists = append(artists, artist)
	}

	return nil, artists
}
