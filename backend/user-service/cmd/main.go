package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"user-service/config"
	handlers "user-service/internal/api"
	"user-service/internal/app"
	"user-service/internal/repository"
	services "user-service/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx := context.Background()

	// Init variables
	var repository repository.UserRepository
	var service services.UserService
	var handler handlers.HandlerInterface

	// Load configuration from YAML file
	config, err := config.LoadConfig()
	log.Print("Config loaded")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Setup logger
	log := setupLogger(config.Env)
	log.Info("Logger setup", slog.Any("cfg", cfg))

	application := app.NewApp(log, config.GRPCPort, config.TokenExpiration)
	application.Run()

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
	// r.LoadHTMLGlob("../static/templ/*")

	r.GET("/", handler.Home)
	r.POST("/register", handler.RegisterUser)
	r.POST("/login", handler.LoginUser)
	// Init Middlewares
	// r.Use(middl.AuthMiddleware())
	r.GET("/user/:id", handler.UserByID)

	// Run router
	r.Run(config.Server.Address)
}
