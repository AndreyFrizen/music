package grpcapp

import (
	config "catalog-service/config/grpc_server"
	handlers "catalog-service/internal/api"
	"catalog-service/internal/app/database"
	"catalog-service/internal/repository"
	services "catalog-service/internal/service"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	config     *config.Config
}

func NewApp(log *slog.Logger, config *config.Config, db *database.DB) *App {
	repo := repository.NewRepository(db)
	userService := services.NewService(repo, log)

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
		slog.Int("port", a.config.GRPCPort),
	)
	log.Info("gRPC server started")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.config.GRPCPort))
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
