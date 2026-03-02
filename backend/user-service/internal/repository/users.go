package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"user-service/internal/model"
)

type userRepository interface {
	CreateUser(u *model.User, ctx context.Context) error
	UserByID(id int, ctx context.Context) (*model.User, error)
	UserByEmail(email string, ctx context.Context) (*model.User, error)
}

type PostgresRepo struct {
	db *sql.DB
}

// Post

// Create user in database
func (s *PostgresRepo) CreateUser(u *model.User, ctx context.Context) error {
	query := fmt.Sprintf("INSERT INTO users (username, password, email) VALUES ('%s', '%s', '%s')",
		u.Username, u.EncryptedPassword, u.Email,
	)

	_, err := s.db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	return nil
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
