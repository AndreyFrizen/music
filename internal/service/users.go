package services

import (
	"mess/internal/model"
)

type userService interface {
	Register(user *model.User) error
	UserByID(id string) (*model.User, error)
}

// Register registers a new user
func (m *Service) Register(user *model.User) error {
	err := validate.Struct(user)
	err = user.EncryptPassword()
	err = m.repo.CreateUser(user)

	return err
}

// UserByID retrieves a user by ID
func (m *Service) UserByID(id string) (*model.User, error) {
	return m.repo.UserByID(id)
}
