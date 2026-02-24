package store

import (
	"context"
	"encoding/json"
	"fmt"
	"mess/internal/model"
	"strconv"
	"time"
)

type playlistRepository interface {
	CreatePlaylist(p *model.Playlist, ctx context.Context) error
	PlaylistByID(id int, ctx context.Context) (*model.Playlist, error)
	DeletePlaylist(id int, ctx context.Context) error
}

// Post

// Add Playlist to database
func (s *Store) CreatePlaylist(p *model.Playlist, ctx context.Context) error {
	query := fmt.Sprintf("INSERT INTO playlists VALUES ('%s', '%v')",
		p.Title, p.UserID)

	_, err := s.db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	jsonData, _ := json.MarshalIndent(p, "", "  ")
	s.cash.Set(ctx, strconv.Itoa(p.ID), string(jsonData), time.Minute*10)
	return nil
}

// Get

// PlaylistByID retrieves a playlist by its ID from the database.
func (s *Store) PlaylistByID(id int, ctx context.Context) (*model.Playlist, error) {
	var playlist model.Playlist

	play, err := s.cash.Get(ctx, strconv.Itoa(id)).Bytes()
	if err == nil {
		err = json.Unmarshal(play, &playlist)
	} else {
		query := fmt.Sprintf("SELECT * FROM playlists WHERE id = %d", id)

		row := s.db.QueryRowContext(ctx, query)

		err := row.Scan(&playlist.ID, &playlist.Title, &playlist.UserID)

		if err != nil {
			return nil, err
		}

		jsonData, _ := json.MarshalIndent(playlist, "", "  ")
		s.cash.Set(ctx, strconv.Itoa(playlist.ID), string(jsonData), time.Minute*10)
	}

	return &playlist, nil
}

// Delete

// Delete Playlist from database
func (s *Store) DeletePlaylist(id int, ctx context.Context) error {
	query := fmt.Sprintf("DELETE FROM playlists WHERE id = %d", id)

	_, err := s.db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	return nil
}
