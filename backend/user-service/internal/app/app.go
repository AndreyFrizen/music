package app

import (
	"log/slog"
	config "user-service/config/grpc_server"
	grpcapp "user-service/internal/app/grpc"
)

type App struct {
	gRPCServer *grpcapp.App
}

func NewApp(log *slog.Logger, config *config.Config) *App {

	grpcApp := grpcapp.NewApp(log, config)
	return &App{
		gRPCServer: grpcApp,
	}
}
