package services

import (
	"context"
	"log/slog"
	"playlist-service/internal/domain/model"

	"github.com/go-playground/validator/v10"
)

type service struct {
	repo     PlaylistRepository
	log      *slog.Logger
	validate *validator.Validate
}

func NewService(repo PlaylistRepository, log *slog.Logger) *service {
	return &service{
		repo:     repo,
		log:      log,
		validate: validator.New(),
	}
}

type PlaylistRepository interface {
	CreatePlaylist(ctx context.Context, p *model.NewPlaylist) (int64, error)
	PlaylistByID(ctx context.Context, id int64) (*model.Playlist, error)
	DeletePlaylist(ctx context.Context, id int64) error
	UpdatePlaylist(ctx context.Context, p *model.Playlist) (int64, error)
	AddTrackToPlaylist(ctx context.Context, p *model.PlaylistTrack) (int64, error)
	RemoveTrackFromPlaylist(ctx context.Context, trackId int64) (int64, error)
}

// PlaylistService creates a new playlist
func (m *service) CreatePlaylist(ctx context.Context, p *model.NewPlaylist) (int64, error) {
	const op = "playlist-service.CreatePlaylist"

	id, err := m.repo.CreatePlaylist(ctx, p)
	if err != nil {
		m.log.Error(op, "error creating playlist", err)
		return 0, err
	}

	return id, nil
}

// PlaylistService retrieves a playlist by ID
func (m *service) PlaylistByID(ctx context.Context, id int64) (*model.Playlist, error) {
	const op = "playlist-service.PlaylistByID"

	p, err := m.repo.PlaylistByID(ctx, id)
	if err != nil {
		m.log.Error(op, "error getting playlist", err)
		return nil, err
	}

	return p, nil
}

// DeletePlaylist deletes a playlist by ID
func (m *service) DeletePlaylist(ctx context.Context, id int64) error {
	const op = "playlist-service.DeletePlaylist"

	err := m.repo.DeletePlaylist(ctx, id)
	if err != nil {
		m.log.Error(op, "error deleting playlist", err)
		return err
	}

	return nil
}

// UpdatePlaylist updates a playlist by ID
func (m *service) UpdatePlaylist(ctx context.Context, p *model.Playlist) (int64, error) {
	const op = "playlist-service.UpdatePlaylist"

	id, err := m.repo.UpdatePlaylist(ctx, p)
	if err != nil {
		m.log.Error(op, "error updating playlist", err)
		return 0, err
	}

	return id, nil
}

// AddTrackToPlaylist adds a track to a playlist
func (m *service) AddTrackToPlaylist(ctx context.Context, p *model.PlaylistTrack) (int64, error) {
	const op = "playlist-service.AddTrackToPlaylist"

	id, err := m.repo.AddTrackToPlaylist(ctx, p)
	if err != nil {
		m.log.Error(op, "error adding track to playlist", err)
		return 0, err
	}

	return id, nil
}

// RemoveTrackFromPlaylist removes a track from a playlist
func (m *service) RemoveTrackFromPlaylist(ctx context.Context, trackId int64) (int64, error) {
	const op = "playlist-service.RemoveTrackFromPlaylist"

	id, err := m.repo.RemoveTrackFromPlaylist(ctx, trackId)
	if err != nil {
		m.log.Error(op, "error removing track from playlist", err)
		return 0, err
	}

	return id, nil
}
