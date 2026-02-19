package services

import (
	"context"
	"log"
	"mess/internal/model"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type userService interface {
	Register(user *model.User, ctx context.Context) error
	UserByID(id int, ctx context.Context) (*model.User, error)
	Auth(user *model.User, ctx context.Context) (string, error)
}

// Register registers a new user
func (m *Service) Register(user *model.User, ctx context.Context) error {
	err := validate.Struct(user)
	err = user.EncryptPassword()
	err = m.Repo.CreateUser(user, ctx)

	m.Auth(user, ctx)

	return err
}

// UserByID retrieves a user by ID
func (m *Service) UserByID(id int, ctx context.Context) (*model.User, error) {
	return m.Repo.UserByID(id, ctx)
}

// GetUserByEmail retrieves a user by email
func (m *Service) Auth(u *model.User, ctx context.Context) (string, error) {
	log.Println("Authenticating user...", u)
	user, err := m.Repo.UserByEmail(u.Email, ctx)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(u.Password))
	if err != nil {
		return "", err
	}

	token, err := generateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

type TokenClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

func generateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		},
		UserID: strconv.Itoa(userID),
	})
	return token.SignedString(model.SecretKey)
}
