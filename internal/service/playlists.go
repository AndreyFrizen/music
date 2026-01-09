package services

import "mess/internal/model"

type playlistService interface {
	CreatePlaylist(playlist *model.Playlist) error
	PlaylistByID(id string) (*model.Playlist, error)
}

// PlaylistService creates a new playlist
func (m *Service) CreatePlaylist(playlist *model.Playlist) error {
	return m.repo.CreatePlaylist(playlist)
}

// PlaylistService retrieves a playlist by ID
func (m *Service) PlaylistByID(id string) (*model.Playlist, error) {
	return m.repo.PlaylistByID(id)
}
