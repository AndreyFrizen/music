package handlers

import (
	"context"
	"log/slog"
	"playlist-service/internal/domain/model"
	"playlist-service/proto/playlist"

	"google.golang.org/grpc"
)

type serverAPI struct {
	playlist.UnimplementedPlaylistServiceServer
	log     *slog.Logger
	service UserService
}

type UserService interface {
	CreatePlaylist(ctx context.Context, p *model.NewPlaylist) (int64, error)
	PlaylistByID(ctx context.Context, id int64) (*model.Playlist, error)
	DeletePlaylist(ctx context.Context, id int64) error
	UpdatePlaylist(ctx context.Context, p *model.Playlist) (int64, error)
	AddTrackToPlaylist(ctx context.Context, p *model.PlaylistTrack) (int64, error)
	RemoveTrackFromPlaylist(ctx context.Context, trackId int64) (int64, error)
}

func NewServerAPI(log *slog.Logger, service UserService) *serverAPI {
	return &serverAPI{
		log:     log,
		service: service,
	}
}

func Register(gRPC *grpc.Server, log *slog.Logger, service UserService) {
	playlist.RegisterPlaylistServiceServer(gRPC, NewServerAPI(log, service))
}
