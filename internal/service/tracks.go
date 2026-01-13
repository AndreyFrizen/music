package services

import (
	"context"
	"mess/internal/model"
)

type trackService interface {
	AddTrack(track *model.Track, ctx context.Context) error
	AddTrackToPlaylist(track *model.PlaylistTrack, ctx context.Context) error
	TrackByID(id string, ctx context.Context) (*model.Track, error)
	TrackFromPlaylist(id string, ctx context.Context) (*model.Track, error)
	TracksByArtist(artistID string, ctx context.Context) ([]model.Track, error)
	TracksByTitle(title string, ctx context.Context) ([]model.Track, error)
	DeleteTrackFromPlaylist(id string, ctx context.Context) error
}

// AddTrack adds a track to the database
func (m *Service) AddTrack(track *model.Track, ctx context.Context) error {
	return m.Repo.AddTrack(track, ctx)
}

// AddTrackToPlaylist adds a track to a playlist
func (m *Service) AddTrackToPlaylist(track *model.PlaylistTrack, ctx context.Context) error {
	return m.Repo.AddTrackToPlaylist(track, ctx)
}

// TrackByID retrieves a track by ID
func (m *Service) TrackByID(id string, ctx context.Context) (*model.Track, error) {
	return m.Repo.TrackByID(id, ctx)
}

// TrackFromPlaylist retrieves a track from a playlist
func (m *Service) TrackFromPlaylist(id string, ctx context.Context) (*model.Track, error) {
	return m.Repo.TrackFromPlaylist(id, ctx)
}

// DeleteTrackFromPlaylist deletes a track from a playlist
func (m *Service) DeleteTrackFromPlaylist(id string, ctx context.Context) error {
	return m.Repo.DeleteTrackFromPlaylist(id, ctx)
}

// TracksByArtist retrieves tracks by artist ID.
func (m *Service) TracksByArtist(artistID string, ctx context.Context) ([]model.Track, error) {
	return m.Repo.TracksByArtist(artistID, ctx)
}

// TracksByTitle retrieves tracks by title.
func (m *Service) TracksByTitle(title string, ctx context.Context) ([]model.Track, error) {
	return m.Repo.TracksByTitle(title, ctx)
}
