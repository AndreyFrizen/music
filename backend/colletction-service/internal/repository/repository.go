package repository

import (
	"collection-service/internal/app/database"
)

type store struct {
	db database.DB
}

func NewStore(db database.DB) *store {
	return &store{db: db}
}

func (s *store)
