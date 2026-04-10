package services

import (
	"context"
	"log/slog"

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
