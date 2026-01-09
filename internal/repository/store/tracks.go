package store

import (
	"fmt"
	"mess/internal/model"

	"github.com/google/uuid"
)

type trackRepository interface {
	AddTrack(t *model.Track) error
	AddTrackToPlaylist(t *model.PlaylistTrack) error
	TrackByID(id string) (*model.Track, error)
	TrackFromPlaylist(id string) (*model.Track, error)
	DeleteTrackFromPlaylist(id string) error
}

// Post

// Add track to database
func (s *Store) AddTrack(t *model.Track) error {
	query := fmt.Sprintf("INSERT INTO tracks VALUES ('%s', '%s', '%s', '%s')",
		uuid.New().String(), t.Title, t.Duration, t.AudioURL)

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

// Add track to playlist
func (s *Store) AddTrackToPlaylist(t *model.PlaylistTrack) error {
	query := fmt.Sprintf("INSERT INTO track_to_playlist VALUES ('%s', '%s', '%d')",
		t.PlaylistID, t.TrackID, t.Position)

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

// Get

// TrackFromPlaylist retrieves a track from playlist from the database.
func (s *Store) TrackFromPlaylist(id string) (*model.Track, error) {
	query := fmt.Sprintf("SELECT * FROM playlist_tracks WHERE id = '%s'", id)

	row := s.db.QueryRow(query)

	var track model.Track

	err := row.Scan(&track.ID, &track.Title, &track.Duration, &track.AudioURL, &track.ArtistID)

	if err != nil {
		return nil, err
	}

	return &track, nil
}

// TrackByID retrieves a track by its ID from the database
func (s *Store) TrackByID(id string) (*model.Track, error) {
	query := fmt.Sprintf("SELECT * FROM tracks WHERE id = '%s'", id)

	row := s.db.QueryRow(query)

	var track model.Track

	err := row.Scan(&track.ID, &track.Title, &track.Duration, &track.AudioURL)

	if err != nil {
		return nil, err
	}

	return &track, nil
}

// Delete

// DeleteTrackFromPlaylist deletes a track from playlist from the database
func (s *Store) DeleteTrackFromPlaylist(id string) error {
	query := fmt.Sprintf("DELETE FROM playlist_tracks WHERE id = '%s'", id)

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}
