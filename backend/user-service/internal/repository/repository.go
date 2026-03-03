package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"user-service/internal/model"
)

type Store struct {
	db *sql.DB
}

func NewRepository(host, port, user, password, dbname string) (*Store, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

// Create user in database
func (s *Store) CreateUser(u *model.User, ctx context.Context) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx,
		"INSERT INTO users (email, username, password) VALUES ($1, $2, $3) RETURNING id",
		u.Email, u.Username, u.EncryptedPassword).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Get user by id
func (s *Store) UserByID(id int, ctx context.Context) (*model.User, error) {
	query := fmt.Sprintf("SELECT * FROM users WHERE id = '%d'", id)

	row := s.db.QueryRowContext(ctx, query)

	var user model.User

	err := row.Scan(&user.ID, &user.Username, &user.EncryptedPassword, &user.Email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Get user by email
func (s *Store) UserByEmail(email string, ctx context.Context) (*model.User, error) {
	query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)

	row := s.db.QueryRowContext(ctx, query)

	var user model.User

	err := row.Scan(&user.ID, &user.Username, &user.EncryptedPassword, &user.Email)
	log.Println("User found:", user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
