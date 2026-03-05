package repository

import (
	"context"
	"database/sql"
	"time"
	"user-service/internal/domain/errors"

	modeluser "user-service/internal/domain/model"

	"github.com/redis/go-redis/v9"
)

type store struct {
	db       *sql.DB
	redis    *redis.Client
	cacheTTL time.Duration
}

func NewRepository(db *sql.DB, redis *redis.Client) *store {
	return &store{
		db:       db,
		redis:    redis,
		cacheTTL: 15 * time.Minute,
	}
}

// Create user in database
func (s *store) Register(ctx context.Context, u *modeluser.User) (int64, error) {
	const op = "repository.UserRepository.CreateUser"

	var id int64
	err := s.db.QueryRowContext(ctx,
		"INSERT INTO users (email, username, password) VALUES ($1, $2, $3) RETURNING id",
		u.Email, u.Username, u.EncryptedPassword).Scan(&id)

	if err != nil {
		return 0, errors.DatabaseError(op, err)
	}

	return id, nil
}

// Get user by id
func (s *store) UserByID(ctx context.Context, id int64) (*modeluser.User, error) {
	const op = "repository.UserRepository.UserByID"

	query := "SELECT id, username, email FROM users WHERE id = $1"

	row := s.db.QueryRowContext(ctx, query, id)

	var user modeluser.User

	err := row.Scan(&user.ID, &user.Username, &user.Email)

	if err == sql.ErrNoRows {
		return nil, errors.NotFoundError(op, "user not found")
	}
	if err != nil {
		return nil, errors.DatabaseError(op, err)
	}

	return &user, nil
}

// Get user by email
func (s *store) UserByEmail(ctx context.Context, email string) (*modeluser.User, error) {
	const op = "repository.UserRepository.UserByEmail"

	query := "SELECT id, username, email FROM users WHERE email = $1"

	row := s.db.QueryRowContext(ctx, query, email)

	var user modeluser.User

	err := row.Scan(&user.ID, &user.Username, &user.Email)

	if err == sql.ErrNoRows {
		return nil, errors.NotFoundError(op, "user not found")
	}
	if err != nil {
		return nil, errors.DatabaseError(op, err)
	}

	return &user, nil
}

// Update user in database
func (s *store) UpdateUser(ctx context.Context, u *modeluser.User) error {
	const op = "repository.UserRepository.UpdateUser"

	result, err := s.db.ExecContext(ctx,
		"UPDATE users SET username = $1 WHERE id = $2",
		u.Username, u.ID)

	if err != nil {
		return s.handleError(op, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return s.handleError(op, err)
	}

	if rows == 0 {
		return s.handleError(op, err)
	}

	return nil
}

func (s *store) UpdateUserEmail(ctx context.Context, u *modeluser.User) error {
	const op = "repository.UserRepository.UpdateUserEmail"

	result, err := s.db.ExecContext(ctx,
		"UPDATE users SET email = $1 WHERE id = $2",
		u.Email, u.ID)

	if err != nil {
		return s.handleError(op, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return s.handleError(op, err)
	}

	if rows == 0 {
		return s.handleError(op, err)
	}

	return nil
}

func (s *store) handleError(op string, err error) error {
	if err == sql.ErrNoRows {
		return errors.NotFoundError(op, "user not found")
	}

	return errors.DatabaseError(op, err)
}
