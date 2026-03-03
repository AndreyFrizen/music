package handlers

import (
	"context"
	"log/slog"
	"user-service/proto/user"

	"google.golang.org/grpc"
)

type serverAPI struct {
	user.UnimplementedUserServiceServer
	logger *slog.Logger
	// service *UserService
}

func Register(gRPC *grpc.Server) {
	user.RegisterUserServiceServer(gRPC, &serverAPI{})
}

func (s *serverAPI) RegisterUser(ctx context.Context, req *user.RegisterUserRequest) (*user.RegisterUserResponse, error) {
	return nil, nil
}

func (s *serverAPI) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	return nil, nil
}
func (s *serverAPI) UpdateUser(context.Context, *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	return nil, nil
}
func (s *serverAPI) DeleteUser(context.Context, *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	return nil, nil
}
func (s *serverAPI) LoginUser(context.Context, *user.LoginUserRequest) (*user.LoginUserResponse, error) {
	return nil, nil
}
func (s *serverAPI) LogoutUser(context.Context, *user.LogoutUserRequest) (*user.LogoutUserResponse, error) {
	return nil, nil
}
