package main

import (
	"database/sql"
	"log"
	"log/slog"
	"mess/internal/api/handlers"
	"mess/internal/config"
	"mess/internal/repository/store"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	// Load configuration from YAML file
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup logger
	log := setupLogger(config.Env)

	// Initialize database connection
	db, err := sql.Open("sqlite3", config.StoragePath)
	if err != nil {
		log.Error("Failed to connect to database: %v", err)
	}
	defer db.Close()

	store := store.NewStore(db)

	// Initialize router
	r := gin.Default()
	r.POST("/register", func(ctx *gin.Context) {
		handlers.RegisterUser(ctx, store)
	})

	// r.POST("/login", handlers.LoginUser())
	// r.POST("/artists", handlers.CreateArtist())
	// r.POST("/tracks", handlers.CreateTrack())
	// r.POST("/albums", handlers.CreateAlbum())
	// r.POST("/playlists", handlers.CreatePlaylist())
	// r.POST("/playlists/:playlist_id/tracks", handlers.AddTrackToPlaylist())
	r.Run("localhost:8080")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
