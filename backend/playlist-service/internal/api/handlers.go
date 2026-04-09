package handlers

import (
	"context"
	"log/slog"
	"playlist-service/internal/domain/errors"
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

func (s *serverAPI) CreatePlaylist(ctx context.Context, req *playlist.CreatePlaylistRequest) (*playlist.CreatePlaylistResponse, error) {
	const op = "handlers.CreatePlaylist"
	s.log.With(slog.String("op", op))

	_, err := s.service.CreatePlaylist(ctx, &model.NewPlaylist{
		Title:  req.Title,
		UserID: req.UserId,
	})
	if err != nil {
		s.log.ErrorContext(ctx, "failed to create playlist", slog.String("op", op), slog.String("error", err.Error()))
		return nil, err
	}
	return &playlist.CreatePlaylistResponse{}, nil
}

func (s *serverAPI) PlaylistByID(ctx context.Context, req *playlist.GetPlaylistRequest) (*playlist.GetPlaylistResponse, error) {
	const op = "handlers.PlaylistByID"
	s.log.With(slog.String("op", op))

	p, err := s.service.PlaylistByID(ctx, req.Id)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to ", slog.String("op", op), slog.String("error", err.Error()))
		return nil, err
	}

	play := &playlist.Playlist{
		Id:     p.ID,
		UserId: p.UserID,
		Title:  p.Title,
	}

	return &playlist.GetPlaylistResponse{Playlist: play}, nil
}

func (s *serverAPI) DeletePlaylist(ctx context.Context, req *playlist.DeletePlaylistRequest) (*playlist.DeletePlaylistResponse, error) {
	const op = "handlers.DeletePlaylist"
	s.log.With(slog.String("op", op))

	err := s.service.DeletePlaylist(ctx, req.Id)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to ", slog.String("op", op), slog.String("error", err.Error()))
		return nil, err
	}
	return &playlist.DeletePlaylistResponse{Success: true}, nil
}

func (s *serverAPI) UpdatePlaylist(ctx context.Context, req *playlist.UpdatePlaylistRequest) (*playlist.UpdatePlaylistResponse, error) {
	const op = "handlers.UpdatePlaylist"
	s.log.With(slog.String("op", op))

	userId, ok := ctx.Value("user_id").(int64)
	if !ok {
		s.log.ErrorContext(ctx, "failed to ", slog.String("op", op), slog.String("error", "user_id not found"))
		return nil, errors.UnauthorizedError(op, "user unauthorized")
	}

	p, err := s.service.UpdatePlaylist(ctx, &model.Playlist{ID: req.Id, Title: req.Title, UserID: userId})
	if err != nil {
		s.log.ErrorContext(ctx, "failed to ", slog.String("op", op), slog.String("error", err.Error()))
		return nil, err
	}

	play := &playlist.Playlist{
		Id:     p,
		UserId: userId,
		Title:  req.Title,
	}

	return &playlist.UpdatePlaylistResponse{Playlist: play}, nil
}

func (s *serverAPI) AddTrackToPlaylist(ctx context.Context, req *playlist.AddTrackRequest) (*playlist.AddTrackResponse, error) {
	const op = "handlers.AddTrackToPlaylist"
	s.log.With(slog.String("op", op))

	p, err := s.service.AddTrackToPlaylist(ctx, &model.PlaylistTrack{PlaylistID: req.Id, TrackID: req.TrackId, Position: req.Position})
	if err != nil {
		s.log.ErrorContext(ctx, "failed to ", slog.String("op", op), slog.String("error", err.Error()))
		return nil, err
	}

	return &playlist.AddTrackResponse{Id: p}, nil
}

func (s *serverAPI) RemoveTrackFromPlaylist(ctx context.Context, req *playlist.RemoveTrackToPlaylistRequest) (*playlist.RemoveTrackToPlaylistResponse, error) {
	const op = "handlers.RemoveTrackFromPlaylist"
	s.log.With(slog.String("op", op))

	p, err := s.service.RemoveTrackFromPlaylist(ctx, req.Id)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to ", slog.String("op", op), slog.String("error", err.Error()))
		return nil, err
	}

	return &playlist.RemoveTrackToPlaylistResponse{Id: p}, nil
}
