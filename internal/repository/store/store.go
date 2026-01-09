package store

import (
	"database/sql"
)

type Repository interface {
	userRepository
	artistRepository
	trackRepository
	albumRepository
	playlistRepository
}
type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}
