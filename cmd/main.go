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
	"net/http"
	"os"
	"strconv"

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

	r.POST("/register", handler.RegisterUser)
	r.POST("/login", handler.LoginUser)
	// Init Middlewares
	// r.Use(middl.AuthMiddleware())
	r.GET("/user/:id", handler.UserByID)

	routes.AlbumsRoutes(r, handler)

	routes.ArtistRoutes(r, handler)

	routes.TracksRoutes(r, handler)

	routes.PlaylistsRoutes(r, handler)

	// Handlers get
	r.GET("/play/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		filep := "/home/andrey/projects/music/static/" + filename

		// Проверяем существование файла
		if _, err := os.Stat(filep); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		// Открываем файл
		file, err := os.Open(filep)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		// Получаем информацию о файле
		fileInfo, _ := file.Stat()

		// Определяем Content-Type
		contentType := "audio/mpeg"

		// Устанавливаем заголовки для потоковой передачи
		c.Header("Content-Type", contentType)
		c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
		c.Header("Accept-Ranges", "bytes")

		// Потоковая передача всего файла
		http.ServeContent(c.Writer, c.Request, filename, fileInfo.ModTime(), file)
	})

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
