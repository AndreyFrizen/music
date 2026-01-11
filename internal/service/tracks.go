package services

import "mess/internal/model"

type trackService interface {
	AddTrack(track *model.Track) error
	AddTrackToPlaylist(track *model.PlaylistTrack) error
	TrackByID(id string) (*model.Track, error)
	TrackFromPlaylist(id string) (*model.Track, error)
	DeleteTrackFromPlaylist(id string) error
}

// AddTrack adds a track to the database
func (m *Service) AddTrack(track *model.Track) error {
	return m.Repo.AddTrack(track)
}

// AddTrackToPlaylist adds a track to a playlist
func (m *Service) AddTrackToPlaylist(track *model.PlaylistTrack) error {
	return m.Repo.AddTrackToPlaylist(track)
}

// TrackByID retrieves a track by ID
func (m *Service) TrackByID(id string) (*model.Track, error) {
	return m.Repo.TrackByID(id)
}

// TrackFromPlaylist retrieves a track from a playlist
func (m *Service) TrackFromPlaylist(id string) (*model.Track, error) {
	return m.Repo.TrackFromPlaylist(id)
}

// DeleteTrackFromPlaylist deletes a track from a playlist
func (m *Service) DeleteTrackFromPlaylist(id string) error {
	return m.Repo.DeleteTrackFromPlaylist(id)
}
