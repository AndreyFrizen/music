package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"mess/internal/api/handlers"
	"mess/internal/config"
	"mess/internal/repository/store"
	routes "mess/internal/routes/rout"
	services "mess/internal/service"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// Constants
const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	ctx := context.Background()

	// Init variables
	var storeRepository store.Repository
	var serviceMusic services.InterfaceService
	var handler handlers.HandlerInterface

	// Load configuration from YAML file
	config, err := config.LoadConfig()
	log.Print("Config loaded")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Setup logger
	log := setupLogger(config.Env)
	log.Info("Logger setup")

	// Initialize database connection
	db, err := sql.Open("sqlite3", config.StoragePath)
	if err != nil {
		log.Error("Failed to connect to database")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Error("Failed to ping database")
	}
	log.Info("Database init")

	// Initialize Redis client
	cash := store.NewClient(config)

	err = cash.Ping(ctx).Err()
	if err != nil {
		log.Error("Failed to connect to Redis")
	}
	log.Info("Redis init")

	// Init Repository
	storeRepository = store.NewStore(db, cash)
	log.Info("init reposiory")
	// Init Services
	serviceMusic = services.NewService(storeRepository)
	log.Info("init services")
	// Init Handlers
	handler = handlers.NewHandler(serviceMusic)
	log.Info("init handlers")

	// Initialize router
	r := gin.Default()
	r.LoadHTMLGlob("../static/templ/*")

	r.GET("/", handler.Home)
	r.POST("/register", handler.RegisterUser)
	r.POST("/login", handler.LoginUser)
	// Init Middlewares
	// r.Use(middl.AuthMiddleware())
	r.GET("/user/:id", handler.UserByID)

	routes.AlbumsRoutes(r, handler)

	routes.ArtistRoutes(r, handler)

	routes.TracksRoutes(r, handler)

	routes.PlaylistsRoutes(r, handler)

	// Run router
	r.Run(config.Server.Address)
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
