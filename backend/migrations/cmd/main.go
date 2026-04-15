package main

import (
	"log/slog"
	"mess/backend/migrations/config"
	"mess/backend/migrations/internal/migrator"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.LoadConfig()
	if err != nil {

	}

	m := migrator.NewMigrator(logger, configs)

	if err := m.Run(); err != nil {
		logger.Error("Migration failed", "error", err)
		os.Exit(1)
	}

	logger.Info("All migrations completed successfully")
}
