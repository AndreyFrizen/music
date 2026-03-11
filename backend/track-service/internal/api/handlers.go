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
	Register(ctx context.Context, req *services.RegisterRequest) (*services.UserResponse, error)
	UserByID(ctx context.Context, id int64) (*services.UserResponse, error)
	UpdateUser(ctx context.Context, req *services.UpdateUserRequest) (*services.UserResponse, error)
	Login(ctx context.Context, req *services.LoginRequest) (*services.LoginResponse, error)
	Logout(ctx context.Context, token string) error
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
