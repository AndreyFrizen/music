package repository

import (
	"context"
	"playlist-service/internal/app/database"
	"playlist-service/internal/domain/errors"
	"playlist-service/internal/domain/model"
	"strconv"
)

type store struct {
	db database.DB
}

// Add Playlist to database.
func (s *store) CreatePlaylist(ctx context.Context, p model.NewPlaylist) (int64, error) {
	const op = "repository.PlaylistRepository.CreatePlaylist"

	query := "INSERT INTO playlists(title, user_id) VALUES ($1, $2) RETURNING id"

	var id int64
	err := s.db.QueryRowContext(ctx, query, p.Title, p.UserID).Scan(&id)

	if err != nil {
		return 0, errors.DatabaseError(op, err)
	}

	s.setPlaylistToCache(ctx, strconv.Itoa(int(id)), &model.Playlist{ID: id, Title: p.Title, UserID: p.UserID})
	return id, nil
}

// PlaylistByID retrieves a playlist by its ID from the database.
func (s *store) PlaylistByID(ctx context.Context, id int64) (*model.Playlist, error) {
	const op = "repository.PlaylistRepository.PlaylistByID"

	key := strconv.Itoa(int(id))
	playlist, err := s.getPlaylistFromCache(ctx, key)
	if err == nil {
		return playlist, nil
	}

	var p model.Playlist
	query := "SELECT * FROM playlists WHERE id = $1"
	row := s.db.QueryRowContext(ctx, query, id)
	err = row.Scan(&p.ID, &p.Title, &p.UserID)
	if err != nil {
		return nil, errors.DatabaseError(op, err)
	}

	s.setPlaylistToCache(ctx, key, &p)
	return &p, nil
}

// Delete Playlist from database.
func (s *store) DeletePlaylist(ctx context.Context, id int64) error {
	const op = "repository.PlaylistRepository.DeletePlaylist"

	query := "DELETE FROM playlists WHERE id = $1"

	_, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return errors.DatabaseError(op, err)
	}

	return nil
}

// UpdatePlaylist updates a playlist in the database.
func (s *store) UpdatePlaylist(ctx context.Context, p model.Playlist) (int64, error) {
	const op = "repository.PlaylistRepository.UpdatePlaylist"

	query := "UPDATE playlists SET title = $1 WHERE id = $2"
	_, err := s.db.ExecContext(ctx, query, p.Title, p.ID)

	if err != nil {
		return 0, errors.DatabaseError(op, err)
	}

	s.setPlaylistToCache(ctx, strconv.Itoa(int(p.ID)), &p)

	return p.ID, nil
}

// AddTrackToPlaylist adds a track to a playlist in the database
func (s *store) AddTrackToPlaylist(ctx context.Context, playlistID int64, trackID int64) (int64, error) {
	const op = "repository.PlaylistRepository.AddTrackToPlaylist"

	query := "INSERT INTO playlist_tracks (playlist_id, track_id) VALUES ($1, $2) RETURNING id"

	var id int64
	err := s.db.QueryRowContext(ctx, query, playlistID, trackID).Scan(&id)
	if err != nil {
		return 0, errors.DatabaseError(op, err)
	}

	return id, nil
}

// RemoveTrackFromPlaylist removes a track from a playlist in the database.
func (s *store) RemoveTrackFromPlaylist(ctx context.Context, trackID int64) (int64, error) {
	const op = "repository.PlaylistRepository.RemoveTrackFromPlaylist"

	query := "DELETE FROM playlist_tracks WHERE track_id = $1 RETURNING playlist_id"

	var playlistID int64
	err := s.db.QueryRowContext(ctx, query, trackID).Scan(&playlistID)
	if err != nil {
		return 0, errors.DatabaseError(op, err)
	}

	return playlistID, nil
}
