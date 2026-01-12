package services

import (
	"mess/internal/model"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("my_secret_key")

type userService interface {
	Register(user *model.User) error
	UserByID(id string) (*model.User, error)
	Auth(user *model.User) (string, error)
}

// Register registers a new user
func (m *Service) Register(user *model.User) error {
	err := validate.Struct(user)
	err = user.EncryptPassword()
	err = m.Repo.CreateUser(user)

	return err
}

// UserByID retrieves a user by ID
func (m *Service) UserByID(id string) (*model.User, error) {
	return m.Repo.UserByID(id)
}

// GetUserByEmail retrieves a user by email
func (m *Service) Auth(u *model.User) (string, error) {
	user, err := m.Repo.UserByEmail(u.Email)

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

func generateToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Срок действия — 24 часа
		},
		UserID: userID.String(),
	})

	return token.SignedString(secretKey)
}
