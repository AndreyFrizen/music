package handlers

import (
	"collection-service/internal/domain/models"
	"context"
	"log/slog"
	"music/collection-service/proto/collection"
	modelsColl "music/collection-service/proto/models"
	"music/collection-service/proto/request"
	"music/collection-service/proto/response"

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

func (h *handler) GetArtists(ctx context.Context, req *request.GetArtistsRequest) (*response.GetArtistsResponse, error) {
	const op = "handlers.CollectionHandlers.GetArtists"

	resp, err := h.service.GetArtists(ctx, &models.GetArtistsRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}

	artists := make([]*modelsColl.Artist, 0, len(resp.Artists))
	for _, artist := range resp.Artists {
		artists = append(artists, &modelsColl.Artist{
			UserId:   artist.UserId,
			ArtistId: artist.ArtistId,
		})
	}

	return &response.GetArtistsResponse{
		Artists: artists,
	}, nil
}

func (h *handler) AddArtist(ctx context.Context, req *request.AddArtistRequest) (*response.AddArtistResponse, error) {
	const op = "handlers.CollectionHandlers.AddArtist"

	resp, err := h.service.AddArtist(ctx, &models.AddArtistRequest{
		UserId:   req.UserId,
		ArtistId: req.ArtistId,
	})
	if err != nil {
		return nil, err
	}

	return &response.AddArtistResponse{
		ArtistId: resp.ArtistId,
	}, nil
}

func (h *handler) DeleteArtist(ctx context.Context, req *request.DeleteArtistRequest) error {
	const op = "handlers.CollectionHandlers.DeleteArtist"

	err := h.service.DeleteArtist(ctx, &models.RemoveArtistRequest{
		UserId:   req.UserId,
		ArtistId: req.ArtistId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (h *handler) GetTracks(ctx context.Context, req *request.GetTracksRequest) (*response.GetTracksResponse, error) {
	const op = "handlers.CollectionHandlers.GetTracks"

	resp, err := h.service.GetTracks(ctx, &models.GetTracksRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}

	tracks := make([]*modelsColl.Track, 0, len(resp.Tracks))
	for _, track := range resp.Tracks {
		tracks = append(tracks, &modelsColl.Track{
			UserId:  track.UserId,
			TrackId: track.TrackId,
		})
	}

	return &response.GetTracksResponse{
		Tracks: tracks,
	}, nil
}

func (h *handler) AddTrack(ctx context.Context, req *request.AddTrackRequest) (*response.AddTrackResponse, error) {
	const op = "handlers.CollectionHandlers.AddTrack"

	resp, err := h.service.AddTrack(ctx, &models.AddTrackRequest{
		UserId:  req.UserId,
		TrackId: req.TrackId,
	})
	if err != nil {
		return nil, err
	}

	return &response.AddTrackResponse{
		TrackId: resp.TrackId,
	}, nil
}

func (h *handler) DeleteTrack(ctx context.Context, req *request.DeleteTrackRequest) error {
	const op = "handlers.CollectionHandlers.DeleteTrack"

	err := h.service.DeleteTrack(ctx, &models.RemoveTrackRequest{
		UserId:  req.UserId,
		TrackId: req.TrackId,
	})
	if err != nil {
		return err
	}

	return nil
}
