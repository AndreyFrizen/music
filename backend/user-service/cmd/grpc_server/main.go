package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	config "user-service/config/grpc_server"
	"user-service/internal/app/database"
	grpcapp "user-service/internal/app/grpc"
	"user-service/internal/lib/logger"
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
