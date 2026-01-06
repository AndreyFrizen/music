package main

import (
	"database/sql"
	"log"
	"log/slog"
	"mess/internal/api/handlers"
	"mess/internal/config"
	"mess/internal/repository/store"
	services "mess/internal/service"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	var storeRepository store.Repository
	var serviceMusic services.MusicServices
	var handler handlers.Handler
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
		log.Error("Failed to connect to database")
	}

	if err := db.Ping(); err != nil {
		log.Error("Failed to ping database")
	}
	defer db.Close()

	storeRepository = store.NewStore(db)
	serviceMusic = services.NewMusicService(storeRepository)
	handler = handlers.NewHandler(serviceMusic)

	// Initialize router
	r := gin.Default()
	r.POST("/register", func(ctx *gin.Context) {
		err := handler.RegisterUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})
	// r.GET("/track", func(ctx *gin.Context) {
	// 	handlers.GetTrackStream(ctx)
	// })
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
