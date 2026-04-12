package main

import (
	config "collection-service/config/grpc_server"
	"collection-service/internal/app/database"
	grpcapp "collection-service/internal/app/grpc"
	"context"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		return
	}

	log := logger.SetupLogger(cfg.Env)
	log.Info("logger set up")

	db, err := database.NewDB(log, cfg)
	if err != nil {
		log.Error("Error connecting to database", "error", err)
		return
	}
	defer db.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	grpcServer := grpcapp.NewApp(log, cfg, db)

	go func() {
		<-ctx.Done()
		grpcServer.Stop()
	}()
	err = grpcServer.Run()
	if err != nil {
		log.Error("Error running gRPC server", "error", err)
		return
	}

}
