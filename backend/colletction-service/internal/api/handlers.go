package handlers

import (
	"collection-service/internal/domain/models"
	"collection-service/proto/collection"
	modelsColl "collection-service/proto/models"
	"collection-service/proto/request"
	"collection-service/proto/response"
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

type handler struct {
	collection.UnimplementedCollectionServiceServer
	service Service
}

type Service interface {
	GetAlbums(ctx context.Context, req *models.GetAlbumsRequest) (*models.GetAlbumsResponse, error)
	AddAlbum(ctx context.Context, req *models.AddAlbumRequest) (*models.AddAlbumResponse, error)
	DeleteAlbum(ctx context.Context, req *models.RemoveAlbumRequest) error
	GetTracks(ctx context.Context, req *models.GetTracksRequest) (*models.GetTracksResponse, error)
	AddTrack(ctx context.Context, req *models.AddTrackRequest) (*models.AddTrackResponse, error)
	DeleteTrack(ctx context.Context, req *models.RemoveTrackRequest) error
	GetArtists(ctx context.Context, req *models.GetArtistsRequest) (*models.GetArtistsResponse, error)
	AddArtist(ctx context.Context, req *models.AddArtistRequest) (*models.AddArtistResponse, error)
	DeleteArtist(ctx context.Context, req *models.RemoveArtistRequest) error
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func Register(gRPC *grpc.Server, log *slog.Logger, service Service) {
	collection.RegisterCollectionServiceServer(gRPC, NewHandler(service))
}

func (h *handler) GetAlbums(ctx context.Context, req *request.GetAlbumsRequest) (*response.GetAlbumsResponse, error) {
	const op = "handlers.CollectionHandlers.GetAlbums"

	resp, err := h.service.GetAlbums(ctx, &models.GetAlbumsRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}

	albums := make([]*modelsColl.Album, 0, len(resp.Albums))
	for _, album := range resp.Albums {
		albums = append(albums, &modelsColl.Album{
			UserId:  album.UserId,
			AlbumId: album.AlbumId,
		})
	}

	return &response.GetAlbumsResponse{
		Albums: albums,
	}, nil
}

func (h *handler) AddAlbum(ctx context.Context, req *request.AddAlbumRequest) (*response.AddAlbumResponse, error) {
	const op = "handlers.CollectionHandlers.AddAlbum"

	resp, err := h.service.AddAlbum(ctx, &models.AddAlbumRequest{
		UserId:  req.UserId,
		AlbumId: req.AlbumId,
	})
	if err != nil {
		return nil, err
	}

	return &response.AddAlbumResponse{
		AlbumId: resp.AlbumId,
	}, nil
}

func (h *handler) DeleteAlbum(ctx context.Context, req *request.DeleteAlbumRequest) error {
	const op = "handlers.CollectionHandlers.DeleteAlbum"

	err := h.service.DeleteAlbum(ctx, &models.RemoveAlbumRequest{
		UserId:  req.UserId,
		AlbumId: req.AlbumId,
	})
	if err != nil {
		return err
	}

	return nil
}
