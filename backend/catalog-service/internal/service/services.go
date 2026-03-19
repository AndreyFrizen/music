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
