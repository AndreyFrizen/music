package services

import (
	"context"
	"log/slog"
	"track-service/internal/domain/errors"
	"track-service/internal/domain/model"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo         TrackRepository
	log          *slog.Logger
	validate     *validator.Validate
	tokenManager *jwt.TokenManager
}

func NewService(repo TrackRepository, log *slog.Logger, jwt *jwt.TokenManager) *service {
	return &service{
		repo:         repo,
		log:          log,
		validate:     validator.New(),
		tokenManager: jwt,
	}
}

type TrackRepository interface {
	CreateTrack(ctx context.Context, req *CreateTrackRequest) (CreateTrackResponse, error)
	TrackByID(ctx context.Context, req *GetTrackRequest) (GetTrackResponse, error)
	UpdateTrack(ctx context.Context, req *UpdateTrackRequest) (UpdateTrackResponse, error)
	DeleteTrack(ctx context.Context, req *DeleteTrackRequest) (DeleteTrackResponse, error)
}

// CreateTrack creates a track in the database
func (s *service) CreateTrack(track *model.Track, ctx context.Context) (int64, error) {
	const op = "service.UserService.Register"

	if err := s.validate.Struct(req); err != nil {
		return errors.ValidationError(op, map[string]string{
			"username": "invalid username",
			"email":    "invalid email",
			"password": "invalid password",
		})
	}

	existing, err := s.repo.UserByEmail(ctx, req.Email)
	if err != nil && !errors.IsNotFound(err) {
		s.log.ErrorContext(ctx, "failed to check existing user",
			"op", op,
			"email", req.Email,
			"error", err,
		)
		return nil, errors.InternalError(op, err)
	}
	if existing != nil {
		return nil, errors.ConflictError(op, "email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to hash password",
			"op", op,
			"error", err,
		)
		return nil, errors.InternalError(op, err)
	}

	user := &modeluser.User{
		Username:          req.Username,
		Email:             req.Email,
		EncryptedPassword: string(hashedPassword),
	}

	id, err := s.repo.Register(ctx, user)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to create user",
			"op", op,
			"error", err,
		)
		return nil, errors.InternalError(op, err)
	}

	s.log.InfoContext(ctx, "user registered successfully",
		"op", op,
		"user_id", id,
		"email", req.Email,
	)

	return &UserResponse{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
	}, nil
}
