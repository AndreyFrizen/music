package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	"user-service/config"
	handlers "user-service/internal/api"
	"user-service/internal/pkg/jwt"
	"user-service/internal/repository"
	services "user-service/internal/service"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	config     *config.Config
}

func NewApp(log *slog.Logger, config *config.Config) *App {
	tokenManager := jwt.NewTokenManager(nil)
	repo := repository.NewRepository(nil, nil)
	userService := services.NewService(repo, log, tokenManager)

	grpcServer := grpc.NewServer()

	handlers.Register(grpcServer, log, userService)

	return &App{
		log:        log,
		gRPCServer: grpcServer,
		config:     config,
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)
	log.Info("gRPC server started")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %s %w", op, err)
	}
	defer l.Close()

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("failed to serve: %s %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("gRPC server stopped")

	a.gRPCServer.GracefulStop()
}
