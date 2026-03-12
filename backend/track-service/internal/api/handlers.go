package handlers

import (
	"context"
	"log/slog"
	services "track-service/internal/service"
	"track-service/proto/track"

	"google.golang.org/grpc"
)

type serverAPI struct {
	track.UnimplementedUserServiceServer
	log     *slog.Logger
	service UserAPI
}

type UserAPI interface {
	TrackByID(ctx context.Context, req *services.GetTrackRequest) (*services.GetTrackResponse, error)
	CreateTrack(ctx context.Context, req *services.CreateTrackRequest) (*services.CreateTrackResponse, error)
	UpdateTrack(ctx context.Context, req *services.UpdateTrackRequest) (*services.UpdateTrackResponse, error)
	DeleteTrack(ctx context.Context, req *services.DeleteTrackRequest) (*services.DeleteTrackResponse, error)
}

func NewServerAPI(log *slog.Logger, service UserAPI) *serverAPI {
	return &serverAPI{
		log:     log,
		service: service,
	}
}

func Register(gRPC *grpc.Server, log *slog.Logger, service UserAPI) {
	track.RegisterUserServiceServer(gRPC, NewServerAPI(log, service))
}
