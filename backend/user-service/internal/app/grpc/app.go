package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	config "user-service/config/grpc_server"
	handlers "user-service/internal/api"
	"user-service/internal/app/database"
	"user-service/internal/pkg/jwt"
	"user-service/internal/repository"
	services "user-service/internal/service"

	"google.golang.org/grpc"
)

type App struct {
	log            *slog.Logger
	externalServer *grpc.Server
	eternalServer  *grpc.Server
	config         *config.Config
}

func NewApp(log *slog.Logger, config *config.Config, db *database.DB) *App {
	tokenManager := jwt.NewTokenManager(config)
	repo := repository.NewRepository(db)
	userService := services.NewService(repo, log, tokenManager)

	externalServer := grpc.NewServer()
	handlers.Register(externalServer, log, userService)

	internalServer := grpc.NewServer()
	handlers.RegisterInternal(internalServer, log, userService)

	return &App{
		log:            log,
		externalServer: externalServer,
		eternalServer:  internalServer,
		config:         config,
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.config.GRPCPort),
	)
	log.Info("gRPC server started")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.config.GRPCPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %s %w", op, err)
	}
	defer l.Close()

	go func() {
		if err := a.externalServer.Serve(l); err != nil {
			log.Error("failed to serve", slog.String("error", err.Error()))
		}
	}()

	l, err = net.Listen("tcp", fmt.Sprintf(":%d", a.config.InternalGRPCPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %s %w", op, err)
	}
	defer l.Close()

	if err := a.eternalServer.Serve(l); err != nil {
		log.Error("failed to serve", slog.String("error", err.Error()))
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("gRPC server stopped")

	a.externalServer.GracefulStop()
	a.eternalServer.GracefulStop()
}
