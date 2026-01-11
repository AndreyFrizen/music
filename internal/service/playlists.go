package services

import "mess/internal/model"

type playlistService interface {
	CreatePlaylist(playlist *model.Playlist) error
	PlaylistByID(id string) (*model.Playlist, error)
	DeletePlaylist(id string) error
}

// PlaylistService creates a new playlist
func (m *Service) CreatePlaylist(playlist *model.Playlist) error {
	return m.repo.CreatePlaylist(playlist)
}

// PlaylistService retrieves a playlist by ID
func (m *Service) PlaylistByID(id string) (*model.Playlist, error) {
	return m.repo.PlaylistByID(id)
}

// DeletePlaylist deletes a playlist by ID
func (m *Service) DeletePlaylist(id string) error {
	return m.repo.DeletePlaylist(id)
}
