package handlers

import (
	"context"
	"io"
	"log/slog"
	"os"
	"track-service/internal/domain/model"
	"music/track-service/proto/track"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	track.UnimplementedTrackServiceServer
	log     *slog.Logger
	service Service
}

type Service interface {
	TrackByID(ctx context.Context, req *model.GetTrackRequest) (*model.GetTrackResponse, error)
	CreateTrack(ctx context.Context, req *model.CreateTrackRequest) (*model.CreateTrackResponse, error)
	UpdateTrack(ctx context.Context, req *model.UpdateTrackRequest) (*model.UpdateTrackResponse, error)
	DeleteTrack(ctx context.Context, req *model.DeleteTrackRequest) (*model.DeleteTrackResponse, error)
}

func NewServerAPI(log *slog.Logger, service Service) *serverAPI {
	return &serverAPI{
		log:     log,
		service: service,
	}
}

func Register(gRPC *grpc.Server, log *slog.Logger, service Service) {
	track.RegisterTrackServiceServer(gRPC, NewServerAPI(log, service))
}

func (s *serverAPI) CreateTrack(stream track.TrackService_CreateTrackServer) error {

	err := s.uploadTrack(stream)
	if err != nil {
		return err
	}

	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot receive chunk: %v", err)
	}

	resp, err := s.service.CreateTrack(context.Background(), &model.CreateTrackRequest{
		Title:    req.GetTrack().Title,
		Duration: req.GetTrack().Duration,
		ArtistID: req.GetTrack().ArtistId,
		AlbumID:  req.GetTrack().AlbumId,
	})

	return stream.SendAndClose(&track.CreateTrackResponse{
		Id: resp.ID,
	})
}

func (s *serverAPI) uploadTrack(stream track.TrackService_CreateTrackServer) error {
	fileHandle, err := os.Create("")
	if err != nil {
		return status.Errorf(codes.Internal, "cannot create audio file: %v", err)
	}
	defer fileHandle.Close()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			os.Remove("")
			return status.Errorf(codes.Unknown, "cannot receive chunk: %v", err)
		}

		chunk := req.GetChunk()
		if chunk == nil {
			continue
		}

		_, err = fileHandle.Write(chunk)
		if err != nil {
			os.Remove("")
			return status.Errorf(codes.Internal, "cannot write chunk: %v", err)
		}
	}

	return nil
}

func (s *serverAPI) GetTrack(ctx context.Context, req *track.GetTrackRequest) (*track.GetTrackResponse, error) {
	resp, err := s.service.TrackByID(ctx, &model.GetTrackRequest{ID: req.Id})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get track: %v", err)
	}

	t := &track.Track{
		Id:       resp.Track.ID,
		Title:    resp.Track.Title,
		Duration: resp.Track.Duration,
		ArtistId: resp.Track.ArtistID,
		AlbumId:  resp.Track.AlbumID,
		AudioUrl: resp.Track.AudioURL,
	}

	return &track.GetTrackResponse{Track: t}, nil
}
