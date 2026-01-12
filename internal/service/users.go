package services

import (
	"mess/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type userService interface {
	Register(user *model.User) error
	UserByID(id string) (*model.User, error)
	Auth(email string, password string) error
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
func (m *Service) Auth(email string, password string) error {
	user, err := m.Repo.UserByEmail(email)

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(user.Password))
	if err != nil {
		return err
	}

	return nil
}
