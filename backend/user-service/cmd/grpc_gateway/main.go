package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"
	configate "user-service/config/gateway"
	gateway "user-service/internal/app/gprc-gateway"
	"user-service/internal/lib/logger"
)

func main() {
	cfg, err := configate.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	fmt.Println(cfg)

	log := logger.SetupLogger(cfg.Env)
	log.Info("logger setup", "port", cfg.HTTPPort)

	app, err := gateway.NewGatewayApp(log, cfg)
	if err != nil {
		log.Error("Error creating gateway app", "error", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := app.Stop(shutdownCtx); err != nil {
			log.Error("Error during shutdown", "error", err)
		}
	}()

	err = app.Run()
	if err != nil {
		log.Error("Error running gateway app", "error", err)
	}
}
