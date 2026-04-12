package services

import (
	"collection-service/internal/domain/errors"
	"collection-service/internal/domain/models"
	"context"
	"log/slog"

	"github.com/go-playground/validator/v10"
)

type service struct {
	repo     Repository
	log      *slog.Logger
	validate *validator.Validate
}

func NewService(repo Repository, log *slog.Logger) *service {
	return &service{
		repo:     repo,
		log:      log,
		validate: validator.New(),
	}
}

type Repository interface {
	AddAlbum(ctx context.Context, album *models.Album) (int64, error)
	GetAlbums(ctx context.Context, userId int64) ([]*models.Album, error)
	DeleteAlbum(ctx context.Context, userId int64, albumId int64) error
	GetTracks(ctx context.Context, userId int64) ([]*models.Track, error)
	AddTrack(ctx context.Context, track *models.Track) (int64, error)
	DeleteTrack(ctx context.Context, userId int64, trackId int64) error
	GetArtists(ctx context.Context, userId int64) ([]*models.Artist, error)
	AddArtist(ctx context.Context, artist *models.Artist) (int64, error)
	DeleteArtist(ctx context.Context, userId int64, artistId int64) error
}

func (s *service) GetAlbums(ctx context.Context, req *models.GetAlbumsRequest) (*models.GetAlbumsResponse, error) {
	const op = "service.CollectionService.GetAlbums"

	s.log.With("op", op)

	albums, err := s.repo.GetAlbums(ctx, req.UserId)
	if err != nil {
		s.log.Error("internall database error", "error", err)
		return nil, errors.InternalError(op, err)
	}

	return &models.GetAlbumsResponse{Albums: albums}, nil
}

func (s *service) AddAlbum(ctx context.Context, req *models.AddAlbumRequest) (*models.AddAlbumResponse, error) {
	const op = "service.CollectionService.AddAlbum"

	s.log.With("op", op)
	userId := ctx.Value("userId").(int64)

	album := &models.Album{UserId: userId, AlbumId: req.AlbumId}
	id, err := s.repo.AddAlbum(ctx, album)
	if err != nil {
		s.log.Error("internall database error", "error", err)
		return nil, errors.InternalError(op, err)
	}

	return &models.AddAlbumResponse{AlbumId: id}, nil
}

func (s *service) DeleteAlbum(ctx context.Context, req *models.RemoveAlbumRequest) error {
	const op = "service.CollectionService.DeleteAlbum"
	s.log.With("op", op)

	userId := ctx.Value("userId").(int64)
	err := s.repo.DeleteAlbum(ctx, userId, req.AlbumId)
	if err != nil {
		s.log.Error("internall database error", "error", err)
		return errors.InternalError(op, err)
	}

	return nil
}

func (s *service) GetTracks(ctx context.Context, req *models.GetTracksRequest) (*models.GetTracksResponse, error) {
	const op = "service.CollectionService.GetTracks"
	s.log.With("op", op)

	userId := ctx.Value("userId").(int64)
	tracks, err := s.repo.GetTracks(ctx, userId)
	if err != nil {
		s.log.Error("internall database error", "error", err)
		return nil, errors.InternalError(op, err)
	}

	return &models.GetTracksResponse{Tracks: tracks}, nil
}

func (s *service) AddTrack(ctx context.Context, req *models.AddTrackRequest) (*models.AddTrackResponse, error) {
	const op = "service.CollectionService.AddTrack"
	s.log.With("op", op)

	track := &models.Track{
		UserId:  req.UserId,
		TrackId: req.TrackId,
	}
	id, err := s.repo.AddTrack(ctx, track)
	if err != nil {
		s.log.Error("internall database error", "error", err)
		return nil, errors.InternalError(op, err)
	}
	return &models.AddTrackResponse{TrackId: id}, nil
}

func (s *service) DeleteTrack(ctx context.Context, req *models.RemoveTrackRequest) error {
	const op = "service.CollectionService.DeleteTrack"
	s.log.With("op", op)

	err := s.repo.DeleteTrack(ctx, req.UserId, req.TrackId)
	if err != nil {
		s.log.Error("internall database error", "error", err)
		return errors.InternalError(op, err)
	}

	return nil
}

func (s *service) GetArtists(ctx context.Context, req *models.GetArtistsRequest) (*models.GetArtistsResponse, error) {
	const op = "service.CollectionService.GetArtists"
	s.log.With("op", op)

	artists, err := s.repo.GetArtists(ctx, req.UserId)
	if err != nil {
		s.log.Error("internall database error", "error", err)
		return nil, errors.InternalError(op, err)
	}
	return &models.GetArtistsResponse{Artists: artists}, nil
}

func (s *service) AddArtist(ctx context.Context, req *models.AddArtistRequest) (*models.AddArtistResponse, error) {
	const op = "service.CollectionService.AddArtist"
	s.log.With("op", op)

	artist := &models.Artist{
		UserId:   req.UserId,
		ArtistId: req.ArtistId,
	}
	id, err := s.repo.AddArtist(ctx, artist)
	if err != nil {
		s.log.Error("internall database error", "error", err)
		return nil, errors.InternalError(op, err)
	}
	return &models.AddArtistResponse{ArtistId: id}, nil
}

func (s *service) DeleteArtist(ctx context.Context, req *models.RemoveArtistRequest) error {
	const op = "service.CollectionService.DeleteArtist"
	s.log.With("op", op)

	err := s.repo.DeleteArtist(ctx, req.UserId, req.ArtistId)
	if err != nil {
		s.log.Error("internall database error", "error", err)
		return errors.InternalError(op, err)
	}
	return nil
}
