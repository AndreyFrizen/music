package model

import (
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID                int    `json:"id" db:"id" redis:"id"`
	Username          string `json:"username" db:"username" redis:"username"`
	Password          string `json:"password" validate:"required,min=6,max=100"`
	EncryptedPassword string `db:"password"`
	Email             string `json:"email" db:"email" validate:"required,email" redis:"email"`
}

func (u *User) EncryptPassword() error {

	if len(u.Password) > 0 {
		enc, err := encryptedPassword(u.Password)
		if err != nil {
			return err
		}
		u.EncryptedPassword = string(enc)
	}

	return nil
}

func encryptedPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
