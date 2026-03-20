package handlers

import (
	"catalog-service/internal/domain/model"
	"catalog-service/proto/catalog"
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	catalog.UnimplementedCatalogServiceServer
	log     *slog.Logger
	service Service
}

type Service interface {
	DeleteAlbum(ctx context.Context, req *model.DeleteAlbumRequest) (*model.DeleteAlbumResponse, error)
	DeleteArtist(ctx context.Context, req *model.DeleteArtistRequest) (*model.DeleteArtistResponse, error)
	CreateAlbum(ctx context.Context, req *model.CreateAlbumRequest) (*model.CreateAlbumResponse, error)
	CreateArtist(ctx context.Context, req *model.CreateArtistRequest) (*model.CreateArtistResponse, error)
	AlbumByID(ctx context.Context, req *model.GetAlbumRequest) (*model.GetAlbumResponse, error)
	ArtistByID(ctx context.Context, req *model.GetArtistRequest) (*model.GetArtistResponse, error)
}

func NewServerAPI(log *slog.Logger, service Service) *serverAPI {
	return &serverAPI{
		log:     log,
		service: service,
	}
}

func Register(gRPC *grpc.Server, log *slog.Logger, service Service) {
	catalog.RegisterCatalogServiceServer(gRPC, NewServerAPI(log, service))
}

func (s *serverAPI) CreateAlbum(ctx context.Context, req *catalog.AddAlbumRequest) (*catalog.AddAlbumResponse, error) {
	const op = "handler.CreateAlbum"

	resp, err := s.service.CreateAlbum(ctx, &model.CreateAlbumRequest{
		Title:       req.Album.Title,
		ArtistID:    req.Album.ArtistId,
		ReleaseDate: req.Album.ReleaseDate,
	})
	if err != nil {
		s.log.ErrorContext(ctx, "failed to create album",
			"op", op,
			"error", err,
		)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &catalog.AddAlbumResponse{
		Id: resp.ID,
	}, nil
}

func (s *serverAPI) CreateArtist(ctx context.Context, req *catalog.AddArtistRequest) (*catalog.AddArtistResponse, error) {
	const op = "handler.CreateArtist"

	resp, err := s.service.CreateArtist(ctx, &model.CreateArtistRequest{
		Name: req.Artist.Name,
	})
	if err != nil {
		s.log.ErrorContext(ctx, "failed to create artist",
			"op", op,
			"error", err,
		)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &catalog.AddArtistResponse{
		Id: resp.ID,
	}, nil
}

func (s *serverAPI) AlbumByID(ctx context.Context, req *catalog.GetAlbumRequest) (*catalog.GetAlbumResponse, error) {
	const op = "handler.GetAlbum"

	resp, err := s.service.AlbumByID(ctx, &model.GetAlbumRequest{
		ID: req.Id,
	})
	if err != nil {
		s.log.ErrorContext(ctx, "failed to get album",
			"op", op,
			"error", err,
		)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &catalog.GetAlbumResponse{
		Album: &catalog.Album{
			Id:          resp.Album.ID,
			Title:       resp.Album.Title,
			ArtistId:    resp.Album.ArtistID,
			ReleaseDate: resp.Album.ReleaseDate,
		},
	}, nil
}

func (s *serverAPI) ArtistByID(ctx context.Context, req *catalog.GetArtistRequest) (*catalog.GetArtistResponse, error) {
	const op = "handler.GetArtist"

	resp, err := s.service.ArtistByID(ctx, &model.GetArtistRequest{
		ID: req.Id,
	})
	if err != nil {
		s.log.ErrorContext(ctx, "failed to get artist",
			"op", op,
			"error", err,
		)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &catalog.GetArtistResponse{
		Artist: &catalog.Artist{
			Id:   resp.Artist.ID,
			Name: resp.Artist.Name,
		},
	}, nil
}

func (s *serverAPI) DeleteAlbum(ctx context.Context, req *catalog.DeleteAlbumRequest) (*catalog.DeleteAlbumResponse, error) {
	const op = "handler.DeleteAlbum"

	id, err := s.service.DeleteAlbum(ctx, &model.DeleteAlbumRequest{
		ID: req.Id,
	})
	if err != nil {
		s.log.ErrorContext(ctx, "failed to delete album",
			"op", op,
			"error", err,
		)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &catalog.DeleteAlbumResponse{
		Id: id.ID,
	}, nil
}

func (s *serverAPI) DeleteArtist(ctx context.Context, req *catalog.DeleteArtistRequest) (*catalog.DeleteArtistResponse, error) {
	const op = "handler.DeleteArtist"

	id, err := s.service.DeleteArtist(ctx, &model.DeleteArtistRequest{
		ID: req.Id,
	})
	if err != nil {
		s.log.ErrorContext(ctx, "failed to delete artist",
			"op", op,
			"error", err,
		)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &catalog.DeleteArtistResponse{
		Id: id.ID,
	}, nil
}
