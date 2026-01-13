package store

import (
	"database/sql"
	"mess/internal/config"

	"github.com/redis/go-redis/v9"
)

type Repository interface {
	userRepository
	artistRepository
	trackRepository
	albumRepository
	playlistRepository
}
type Store struct {
	db   *sql.DB
	cash *redis.Client
}

func NewStore(db *sql.DB, cash *redis.Client) *Store {
	return &Store{
		db:   db,
		cash: cash,
	}
}

func NewClient(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis,
		Password: "",
		DB:       0,
	})
	return client
}
