package services

import (
	"context"
	"log/slog"
	"user-service/internal/domain/errors"
	modeluser "user-service/internal/domain/model"
	"user-service/internal/pkg/jwt"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo         UserRepository
	log          *slog.Logger
	validate     *validator.Validate
	tokenManager *jwt.TokenManager
}

func NewService(repo UserRepository, log *slog.Logger, jwt *jwt.TokenManager) *service {
	return &service{
		repo:         repo,
		log:          log,
		validate:     validator.New(),
		tokenManager: jwt,
	}
}

type UserRepository interface {
	Register(ctx context.Context, user *modeluser.User) (int64, error)
	UserByID(ctx context.Context, id int64) (*modeluser.User, error)
	UserByEmail(ctx context.Context, email string) (*modeluser.User, error)
	UpdateUser(ctx context.Context, user *modeluser.User) error
	UpdateUserEmail(ctx context.Context, user *modeluser.User) error
}

// UserByID retrieves a user by ID
func (s *service) UserByID(ctx context.Context, id int64) (*UserResponse, error) {
	const op = "service.UserService.GetUserByID"

	if id <= 0 {
		return nil, errors.ValidationError(op, map[string]string{
			"id": "must be greater than 0",
		})
	}

	user, err := s.repo.UserByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundError(op, "user not found")
		}
		s.log.ErrorContext(ctx, "failed to get user",
			"op", op,
			"id", id,
			"error", err,
		)
		return nil, errors.InternalError(op, err)
	}

	if user == nil {
		return nil, errors.NotFoundError(op, "user not found")
	}

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// Register registers a new user
func (s *service) Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error) {
	const op = "service.UserService.Register"

	if err := s.validate.Struct(req); err != nil {
		return nil, errors.ValidationError(op, map[string]string{
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

// Login authenticates a user and returns tokens
func (s *service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	const op = "service.UserService.Login"

	if err := s.validate.Struct(req); err != nil {
		return nil, errors.ValidationError(op, map[string]string{
			"email":    "not a valid email",
			"password": "not a valid password",
		})
	}

	user, err := s.repo.UserByEmail(ctx, req.Email)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.UnauthorizedError(op, "invalid email or password")
		}
		s.log.ErrorContext(ctx, "failed to get user by email",
			"op", op,
			"email", req.Email,
			"error", err,
		)
		return nil, errors.InternalError(op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(req.Password))
	if err != nil {
		s.log.WarnContext(ctx, "invalid password attempt",
			"op", op,
			"email", req.Email,
		)
		return nil, errors.UnauthorizedError(op, "invalid email or password")
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(user.ID, user.Email, "user")
	if err != nil {
		s.log.ErrorContext(ctx, "failed to generate access token",
			"op", op,
			"user_id", user.ID,
			"error", err,
		)
		return nil, errors.InternalError(op, err)
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(user.ID)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to generate refresh token",
			"op", op,
			"user_id", user.ID,
			"error", err,
		)
		return nil, errors.InternalError(op, err)
	}

	s.log.InfoContext(ctx, "user logged in successfully",
		"op", op,
		"user_id", user.ID,
	)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Logout invalidates user's token (if you have token blacklist)
func (s *service) Logout(ctx context.Context, token string) error {
	const op = "service.UserService.Logout"

	_, err := s.tokenManager.ValidateToken(token)
	if err != nil {
		return errors.UnauthorizedError(op, "invalid token")
	}

	s.log.InfoContext(ctx, "user logged out",
		"op", op,
	)

	return nil
}

// UpdateUser updates user information
func (s *service) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UserResponse, error) {
	const op = "service.UserService.UpdateUser"

	if req.ID <= 0 {
		return nil, errors.ValidationError(op, map[string]string{
			"id": "must be greater than 0",
		})
	}

	if err := s.validate.Struct(req); err != nil {
		return nil, errors.ValidationError(op, map[string]string{
			"username": "not valid",
			"email":    "not valid",
		})
	}

	user, err := s.repo.UserByID(ctx, req.ID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundError(op, "user not found")
		}
		s.log.ErrorContext(ctx, "failed to get user for update",
			"op", op,
			"id", req.ID,
			"error", err,
		)
		return nil, errors.InternalError(op, err)
	}

	if req.Username != "" {
		user.Username = req.Username
	}

	if req.Email != "" && req.Email != user.Email {
		existing, err := s.repo.UserByEmail(ctx, req.Email)
		if err != nil && !errors.IsNotFound(err) {
			s.log.ErrorContext(ctx, "failed to check email availability",
				"op", op,
				"email", req.Email,
				"error", err,
			)
			return nil, errors.InternalError(op, err)
		}
		if existing != nil && existing.ID != req.ID {
			return nil, errors.ConflictError(op, "email already taken")
		}

		user.Email = req.Email
		err = s.repo.UpdateUserEmail(ctx, user)
		if err != nil {
			s.log.ErrorContext(ctx, "failed to update user email",
				"op", op,
				"id", req.ID,
				"error", err,
			)
			return nil, errors.InternalError(op, err)
		}
	} else {
		err = s.repo.UpdateUser(ctx, user)
		if err != nil {
			s.log.ErrorContext(ctx, "failed to update user",
				"op", op,
				"id", req.ID,
				"error", err,
			)
			return nil, errors.InternalError(op, err)
		}
	}

	s.log.InfoContext(ctx, "user updated successfully",
		"op", op,
		"user_id", req.ID,
	)

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
