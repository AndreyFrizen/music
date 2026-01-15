package services

import (
	"context"
	"mess/internal/model"
)

type playlistService interface {
	CreatePlaylist(playlist *model.Playlist, ctx context.Context) error
	PlaylistByID(id int, ctx context.Context) (*model.Playlist, error)
	DeletePlaylist(id int, ctx context.Context) error
}

// PlaylistService creates a new playlist
func (m *Service) CreatePlaylist(playlist *model.Playlist, ctx context.Context) error {
	return m.Repo.CreatePlaylist(playlist, ctx)
}

// PlaylistService retrieves a playlist by ID
func (m *Service) PlaylistByID(id int, ctx context.Context) (*model.Playlist, error) {
	return m.Repo.PlaylistByID(id, ctx)
}

// DeletePlaylist deletes a playlist by ID
func (m *Service) DeletePlaylist(id int, ctx context.Context) error {
	return m.Repo.DeletePlaylist(id, ctx)
}
