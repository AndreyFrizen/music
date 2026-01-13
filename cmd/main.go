package main

import (
	"database/sql"
	"log"
	"log/slog"
	"mess/internal/api/handlers"
	"mess/internal/config"
	middl "mess/internal/lib/middlware"
	"mess/internal/repository/store"
	services "mess/internal/service"
	storage_redis "mess/storage/redis"
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
	var serviceMusic services.InterfaceService
	var handler handlers.HandlerInterface
	var cash storage_redis.CashInterface
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
	log.Info("init reposiory")
	storeRepository = store.NewStore(db)
	cash = storage_redis.NewClient(config)
	_ = cash
	log.Info("init services")
	serviceMusic = services.NewService(storeRepository)
	log.Info("init handlers")
	handler = handlers.NewHandler(serviceMusic)

	// Initialize router
	r := gin.Default()

	// Handlers post

	r.POST("/register", handler.RegisterUser)

	r.POST("/login", func(ctx *gin.Context) {
		handler.LoginUser(ctx)
	})

	r.Use(middl.AuthMiddleware())

	r.POST("/login", handler.LoginUser)

	r.POST("/addartist", handler.CreateArtist)

	r.POST("/addalbum", handler.AddAlbum)

	r.POST("/addtrack", func(ctx *gin.Context) {
		handler.AddTrack(ctx)
	})

	r.POST("/addalbum", func(ctx *gin.Context) {
		handler.AddAlbum(ctx)
	})

	r.POST("/addplaylist", func(ctx *gin.Context) {
		handler.CreatePlaylist(ctx)

	})

	r.POST("/addtracktoplaylist", func(ctx *gin.Context) {
		handler.AddTrackToPlaylist(ctx)
	})

	r.POST("/addplaylist", handler.CreatePlaylist)

	r.POST("/addtracktoplaylist", handler.AddTrackToPlaylist)

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

	r.GET("/artists", handler.Artists)

	// Run router
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
