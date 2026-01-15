package services

import (
	"context"
	"mess/internal/model"
)

type trackService interface {
	AddTrack(track *model.Track, ctx context.Context) error
	AddTrackToPlaylist(track *model.PlaylistTrack, ctx context.Context) error
	TrackByID(id int, ctx context.Context) (*model.Track, error)
	TrackFromPlaylist(id int, ctx context.Context) (*model.Track, error)
	TracksByArtist(artistID int, ctx context.Context) ([]model.Track, error)
	TracksByTitle(title string, ctx context.Context) ([]model.Track, error)
	DeleteTrackFromPlaylist(id int, ctx context.Context) error
	FindTracks(input string) (error, []model.Track)
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
func (m *Service) TrackByID(id int, ctx context.Context) (*model.Track, error) {
	return m.Repo.TrackByID(id, ctx)
}

// TrackFromPlaylist retrieves a track from a playlist
func (m *Service) TrackFromPlaylist(id int, ctx context.Context) (*model.Track, error) {
	return m.Repo.TrackFromPlaylist(id, ctx)
}

// DeleteTrackFromPlaylist deletes a track from a playlist
func (m *Service) DeleteTrackFromPlaylist(id int, ctx context.Context) error {
	return m.Repo.DeleteTrackFromPlaylist(id, ctx)
}

// TracksByArtist retrieves tracks by artist ID.
func (m *Service) TracksByArtist(artistID int, ctx context.Context) ([]model.Track, error) {
	return m.Repo.TracksByArtist(artistID, ctx)
}

// TracksByTitle retrieves tracks by title.
func (m *Service) TracksByTitle(title string, ctx context.Context) ([]model.Track, error) {
	return m.Repo.TracksByTitle(title, ctx)
}

// FindTracks retrieves tracks by input string.
func (m *Service) FindTracks(input string) (error, []model.Track) {
	return m.Repo.FindTracks(input)
}
