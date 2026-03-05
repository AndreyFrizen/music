package handlers

import (
	"context"
	"log/slog"
	"strings"
	services "user-service/internal/service"
	"user-service/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	user.UnimplementedUserServiceServer
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
	user.RegisterUserServiceServer(gRPC, NewServerAPI(log, service))
}

func (s *serverAPI) Register(ctx context.Context, req *user.RegisterUserRequest) (*user.UserResponse, error) {
	const op = "handler.RegisterUser"

	serviceReq := &services.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := s.service.Register(ctx, serviceReq)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to register user",
			"op", op,
			"email", req.Email,
			"error", err,
		)
		return nil, s.handleError(err)
	}

	return &user.UserResponse{
		Id:       resp.ID,
		Username: resp.Username,
		Email:    resp.Email,
	}, nil
}

func (s *serverAPI) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.UserResponse, error) {
	const op = "handler.GetUser"

	resp, err := s.service.UserByID(ctx, req.Id)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to get user",
			"op", op,
			"id", req.Id,
			"error", err,
		)
		return nil, s.handleError(err)
	}

	return &user.UserResponse{
		Id:       resp.ID,
		Username: resp.Username,
		Email:    resp.Email,
	}, nil
}

func (s *serverAPI) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UserResponse, error) {
	const op = "handler.UpdateUser"

	serviceReq := &services.UpdateUserRequest{
		ID:       req.Id,
		Username: req.Username,
		Email:    req.Email,
	}

	resp, err := s.service.UpdateUser(ctx, serviceReq)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to update user",
			"op", op,
			"id", req.Id,
			"error", err,
		)
		return nil, s.handleError(err)
	}

	return &user.UserResponse{
		Id:       resp.ID,
		Username: resp.Username,
		Email:    resp.Email,
	}, nil
}

func (s *serverAPI) LoginUser(ctx context.Context, req *user.LoginUserRequest) (*user.LoginUserResponse, error) {
	const op = "handler.LoginUser"

	serviceReq := &services.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := s.service.Login(ctx, serviceReq)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to login",
			"op", op,
			"email", req.Email,
			"error", err,
		)
		return nil, s.handleError(err)
	}

	return &user.LoginUserResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (s *serverAPI) LogoutUser(ctx context.Context, req *user.LogoutUserRequest) (*user.LogoutUserResponse, error) {
	const op = "handler.LogoutUser"

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.log.ErrorContext(ctx, "no metadata in context",
			"op", op,
		)
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		s.log.ErrorContext(ctx, "no authorization header",
			"op", op,
		)
		return nil, status.Error(codes.Unauthenticated, "missing authorization header")
	}

	authHeader := authHeaders[0]
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		s.log.ErrorContext(ctx, "invalid authorization header format",
			"op", op,
			"header", authHeader,
		)
		return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
	}
	token := parts[1]

	err := s.service.Logout(ctx, token)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to logout",
			"op", op,
			"error", err,
		)
		return nil, s.handleError(err)
	}

	return &user.LogoutUserResponse{
		Success: true,
		Message: "logged out successfully",
	}, nil
}

// handleError конвертирует ошибки сервиса в gRPC статусы
func (s *serverAPI) handleError(err error) error {
	// Здесь можно добавить маппинг ваших errors в gRPC codes
	// Например:
	// if errors.IsNotFound(err) {
	//     return status.Error(codes.NotFound, err.Error())
	// }
	// if errors.IsValidation(err) {
	//     return status.Error(codes.InvalidArgument, err.Error())
	// }
	// if errors.IsConflict(err) {
	//     return status.Error(codes.AlreadyExists, err.Error())
	// }
	// if errors.IsUnauthorized(err) {
	//     return status.Error(codes.Unauthenticated, err.Error())
	// }

	return status.Error(codes.Internal, "internal error")
}
