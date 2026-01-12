package services

import (
	"mess/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type userService interface {
	Register(user *model.User) error
	UserByID(id string) (*model.User, error)
	Auth(user *model.User) error
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
func (m *Service) Auth(u *model.User) error {
	user, err := m.Repo.UserByEmail(u.Email)

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(u.Password))
	if err != nil {
		return err
	}

	return nil
}
