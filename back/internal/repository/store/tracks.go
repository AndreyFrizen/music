package store

import (
	"context"
	"encoding/json"
	"fmt"
	"mess/internal/model"
	"strconv"
	"time"
)

type trackRepository interface {
	AddTrack(t *model.Track, ctx context.Context) error
	AddTrackToPlaylist(t *model.PlaylistTrack, ctx context.Context) error
	TrackByID(id int, ctx context.Context) (*model.Track, error)
	TrackFromPlaylist(id int, ctx context.Context) (*model.Track, error)
	DeleteTrackFromPlaylist(id int, ctx context.Context) error
	TracksByTitle(title string, ctx context.Context) ([]model.Track, error)
	TracksByArtist(artistID int, ctx context.Context) ([]model.Track, error)
	FindTracks(input string) (error, []model.Track)
}

// Post

// Add track to database
func (s *Store) AddTrack(t *model.Track, ctx context.Context) error {
	query := fmt.Sprintf("INSERT INTO tracks(title, duration, audio_url, artist_id) VALUES ('%s', '%d', '%s', '%d')",
		t.Title, t.Duration, t.AudioURL, t.ArtistID)

	_, err := s.db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	return nil
}

// Add track to playlist
func (s *Store) AddTrackToPlaylist(t *model.PlaylistTrack, ctx context.Context) error {
	query := fmt.Sprintf("INSERT INTO track_to_playlist VALUES ('%v', '%v', '%d')",
		t.PlaylistID, t.TrackID, t.Position)

	_, err := s.db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	jsonData, _ := json.MarshalIndent(t, "", "  ")
	s.cash.Set(ctx, strconv.Itoa(t.TrackID), string(jsonData), time.Minute*10)
	return nil
}

// Get

// TrackFromPlaylist retrieves a track from playlist from the database.
func (s *Store) TrackFromPlaylist(id int, ctx context.Context) (*model.Track, error) {
	var track model.Track

	tr, err := s.cash.Get(ctx, strconv.Itoa(id)).Bytes()
	if err == nil {
		err = json.Unmarshal(tr, &track)
	} else {
		query := fmt.Sprintf("SELECT * FROM playlist_tracks WHERE id = '%v'", id)

		row := s.db.QueryRowContext(ctx, query)

		err := row.Scan(&track.ID, &track.Title, &track.Duration, &track.AudioURL, &track.ArtistID)

		if err != nil {
			return nil, err
		}

		jsonData, _ := json.MarshalIndent(track, "", "  ")
		s.cash.Set(ctx, strconv.Itoa(track.ID), string(jsonData), time.Minute*10)
	}

	return &track, nil
}

// TrackByID retrieves a track by its ID from the database
func (s *Store) TrackByID(id int, ctx context.Context) (*model.Track, error) {
	var track model.Track

	tr, err := s.cash.Get(ctx, strconv.Itoa(id)).Bytes()
	if err == nil {
		err = json.Unmarshal(tr, &track)
	} else {
		query := fmt.Sprintf("SELECT * FROM tracks WHERE id = '%v'", id)

		row := s.db.QueryRowContext(ctx, query)

		err := row.Scan(&track.ID, &track.Title, &track.Duration, &track.AudioURL, &track.ArtistID)

		if err != nil {
			return nil, err
		}

		jsonData, _ := json.MarshalIndent(track, "", "  ")
		s.cash.Set(ctx, strconv.Itoa(track.ID), string(jsonData), time.Minute*10)
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
func (s *Store) TracksByArtist(artistID int, ctx context.Context) ([]model.Track, error) {
	query := fmt.Sprintf("SELECT * FROM tracks WHERE artist_id = '%v'", artistID)

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
func (s *Store) DeleteTrackFromPlaylist(id int, ctx context.Context) error {
	query := fmt.Sprintf("DELETE FROM playlist_tracks WHERE id = '%v'", id)

	_, err := s.db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	return nil
}

// FindTracks
func (s *Store) FindTracks(input string) (error, []model.Track) {
	var tracks []model.Track

	t := fmt.Sprintf("SELECT * FROM tracks WHERE title = '%v*'", input)
	rowsTracks, err := s.db.Query(t)
	if err != nil {
		return err, nil
	}
	defer rowsTracks.Close()

	for rowsTracks.Next() {
		var track model.Track

		err := rowsTracks.Scan(&track.ID, &track.Title, &track.AudioURL, &track.Duration, &track.ArtistID)

		if err != nil {
			return err, nil
		}

		tracks = append(tracks, track)
	}

	return nil, tracks
}
