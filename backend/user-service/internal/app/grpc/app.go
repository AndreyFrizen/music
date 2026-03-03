package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	usergrpc "user-service/internal/api"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(log *slog.Logger, port int) *App {
	grpcServer := grpc.NewServer()
	usergrpc.Register(grpcServer)
	return &App{
		log:        log,
		gRPCServer: grpcServer,
		port:       port,
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
