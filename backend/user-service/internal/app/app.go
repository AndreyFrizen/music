package app

import (
	"log/slog"
	grpcapp "user-service/internal/app/grpc"
)

type App struct {
	gRPCServer *grpcapp.App
}

func NewApp(log *slog.Logger, port int, token string) *App {

	grpcApp := grpcapp.NewApp(log, port)
	return &App{
		gRPCServer: grpcApp,
	}
}
