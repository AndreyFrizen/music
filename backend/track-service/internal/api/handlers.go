package handlers

import (
	"context"
	"io"
	"log/slog"
	"os"
	services "track-service/internal/service"
	"track-service/proto/track"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *serverAPI) CreateTrack(stream track.UserService_CreateTrackServer) error {
	const op = "serverAPI.CreateTrack"
	fileHandle, err := os.Create("g")
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
			// Удаляем файл при ошибке
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

	return stream.SendAndClose(&track.CreateTrackResponse{
		Id: 1,
	})
}
