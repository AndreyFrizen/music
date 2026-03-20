package services

import (
	"catalog-service/internal/domain/errors"
	"catalog-service/internal/domain/model"
	"context"
	"log/slog"

	"github.com/go-playground/validator/v10"
)

type service struct {
	repo     CatalogRepository
	log      *slog.Logger
	validate *validator.Validate
}

func NewService(repo CatalogRepository, log *slog.Logger) *service {
	return &service{
		repo:     repo,
		log:      log,
		validate: validator.New(),
	}
}

type CatalogRepository interface {
	CreateAlbum(ctx context.Context, a *model.Album) (int64, error)
	CreateArtist(ctx context.Context, a *model.Artist) (int64, error)
	ArtistByID(ctx context.Context, id int64) (*model.Artist, error)
	AlbumByID(ctx context.Context, id int64) (*model.Album, error)
	DeleteArtist(ctx context.Context, id int64) (int64, error)
	DeleteAlbum(ctx context.Context, id int64) (int64, error)
}

// CreateAlbum creates an album in the database
func (s *service) CreateAlbum(ctx context.Context, req *model.CreateAlbumRequest) (*model.CreateAlbumResponse, error) {
	const op = "service.CatalogService.CreateAlbum"

	if err := s.validate.Struct(req); err != nil {
		return nil, errors.ValidationError(op, map[string]string{
			"track": "invalid track",
		})
	}

	a := &model.Album{
		Title:       req.Title,
		ArtistID:    req.ArtistID,
		ReleaseDate: req.ReleaseDate,
	}

	albumId, err := s.repo.CreateAlbum(ctx, a)
	if err != nil {
		s.log.Error(op, "failed to create album", err)
		return nil, errors.InternalError(op, err)
	}

	s.log.Info(op, "album created successfully", req.Title)

	return &model.CreateAlbumResponse{ID: albumId}, nil
}

// CreateArtist creates an artist in the database
func (s *service) CreateArtist(ctx context.Context, req *model.CreateArtistRequest) (*model.CreateArtistResponse, error) {
	const op = "service.CatalogService.CreateArtist"

	if err := s.validate.Struct(req); err != nil {
		return nil, errors.ValidationError(op, map[string]string{
			"track": "invalid track",
		})
	}

	a := &model.Artist{
		Name: req.Name,
	}

	artistId, err := s.repo.CreateArtist(ctx, a)
	if err != nil {
		s.log.Error(op, "failed to create album", err)
		return nil, errors.InternalError(op, err)
	}

	s.log.Info(op, "artist created successfully", req.Name)

	return &model.CreateArtistResponse{ID: artistId}, nil
}

// ArtistByID returns an artist by its ID.
func (s *service) ArtistByID(ctx context.Context, req *model.GetArtistRequest) (*model.GetArtistResponse, error) {
	const op = "service.CatalogService.ArtistByID"

	artist, err := s.repo.ArtistByID(ctx, req.ID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundError(op, "artist not found")
		}
		s.log.Error(op, "failed to get artist", err)
		return nil, errors.InternalError(op, err)
	}

	s.log.Info(op, "artist retrieved successfully", artist.ID)

	return &model.GetArtistResponse{
		Artist: artist,
	}, nil
}

// AlbumByID returns an album by its ID.
func (s *service) AlbumByID(ctx context.Context, req *model.GetAlbumRequest) (*model.GetAlbumResponse, error) {
	const op = "service.CatalogService.AlbumByID"

	album, err := s.repo.AlbumByID(ctx, req.ID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundError(op, "album not found")
		}
		s.log.Error(op, "failed to get album", err)
		return nil, errors.InternalError(op, err)
	}

	s.log.Info(op, "album retrieved successfully", album.ID)

	return &model.GetAlbumResponse{
		Album: album,
	}, nil
}

func (s *service) DeleteArtist(ctx context.Context, req *model.DeleteArtistRequest) (*model.DeleteArtistResponse, error) {
	const op = "service.CatalogService.DeleteArtist"

	id, err := s.repo.DeleteArtist(ctx, req.ID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundError(op, "not found")
		}
		s.log.Error(op, "failed to delete", err)
		return nil, errors.InternalError(op, err)
	}

	s.log.Info(op, "deleted successfully", req.ID)

	return &model.DeleteArtistResponse{
		ID: id,
	}, nil
}
func (s *service) DeleteAlbum(ctx context.Context, req *model.DeleteAlbumRequest) (*model.DeleteAlbumResponse, error) {
	const op = "service.CatalogService.DeleteAlbum"

	id, err := s.repo.DeleteAlbum(ctx, req.ID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundError(op, "not found")
		}
		s.log.Error(op, "failed to delete", err)
		return nil, errors.InternalError(op, err)
	}

	s.log.Info(op, "deleted successfully", req.ID)

	return &model.DeleteAlbumResponse{
		ID: id,
	}, nil
}
