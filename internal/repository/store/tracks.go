package store

import (
	"context"
	"fmt"
	"mess/internal/model"

	"github.com/google/uuid"
)

type trackRepository interface {
	AddTrack(t *model.Track, ctx context.Context) error
	AddTrackToPlaylist(t *model.PlaylistTrack, ctx context.Context) error
	TrackByID(id string, ctx context.Context) (*model.Track, error)
	TrackFromPlaylist(id string, ctx context.Context) (*model.Track, error)
	DeleteTrackFromPlaylist(id string, ctx context.Context) error
	TracksByTitle(title string, ctx context.Context) ([]model.Track, error)
	TracksByArtist(artistID string, ctx context.Context) ([]model.Track, error)
}

// Post

// Add track to database
func (s *Store) AddTrack(t *model.Track, ctx context.Context) error {
	query := fmt.Sprintf("INSERT INTO tracks VALUES ('%s', '%s', '%s', '%s')",
		uuid.New().String(), t.Title, t.Duration, t.AudioURL)

	_, err := s.db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	return nil
}

// Add track to playlist
func (s *Store) AddTrackToPlaylist(t *model.PlaylistTrack, ctx context.Context) error {
	query := fmt.Sprintf("INSERT INTO track_to_playlist VALUES ('%s', '%s', '%d')",
		t.PlaylistID, t.TrackID, t.Position)

	_, err := s.db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	return nil
}

// Get

// TrackFromPlaylist retrieves a track from playlist from the database.
func (s *Store) TrackFromPlaylist(id string, ctx context.Context) (*model.Track, error) {
	query := fmt.Sprintf("SELECT * FROM playlist_tracks WHERE id = '%s'", id)

	row := s.db.QueryRowContext(ctx, query)

	var track model.Track

	err := row.Scan(&track.ID, &track.Title, &track.Duration, &track.AudioURL, &track.ArtistID)

	if err != nil {
		return nil, err
	}

	return &track, nil
}

// TrackByID retrieves a track by its ID from the database
func (s *Store) TrackByID(id string, ctx context.Context) (*model.Track, error) {
	query := fmt.Sprintf("SELECT * FROM tracks WHERE id = '%s'", id)

	row := s.db.QueryRowContext(ctx, query)

	var track model.Track

	err := row.Scan(&track.ID, &track.Title, &track.Duration, &track.AudioURL)

	if err != nil {
		return nil, err
	}

	return &track, nil
}

// TracksByTitle retrieves tracks by title from the database
func (s *Store) TracksByTitle(title string, ctx context.Context) ([]model.Track, error) {
	query := fmt.Sprintf("SELECT * FROM tracks WHERE title = '%s'", title)

	rows, err := s.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tracks []model.Track

	for rows.Next() {
		var track model.Track

		err := rows.Scan(&track.ID, &track.Title, &track.Duration, &track.AudioURL)

		if err != nil {
			return nil, err
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

// TracksByArtist retrieves tracks by artist from the database
func (s *Store) TracksByArtist(artistID string, ctx context.Context) ([]model.Track, error) {
	query := fmt.Sprintf("SELECT * FROM tracks WHERE artist_id = '%s'", artistID)

	rows, err := s.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tracks []model.Track

	for rows.Next() {
		var track model.Track

		err := rows.Scan(&track.ID, &track.Title, &track.Duration, &track.AudioURL)

		if err != nil {
			return nil, err
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

// Delete

// DeleteTrackFromPlaylist deletes a track from playlist from the database
func (s *Store) DeleteTrackFromPlaylist(id string, ctx context.Context) error {
	query := fmt.Sprintf("DELETE FROM playlist_tracks WHERE id = '%s'", id)

	_, err := s.db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	return nil
}
