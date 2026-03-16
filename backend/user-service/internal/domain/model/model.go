package model

// User represents a user in the system
type User struct {
	ID                int64  `db:"id" redis:"id"`
	Username          string `db:"username" redis:"username"`
	Password          string
	EncryptedPassword string `db:"password"`
	Email             string `db:"email" redis:"email"`
}

type UserRequest struct {
	ID int64 `json:"id"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type RegisterResponse struct {
	ID int64 `json:"id"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type UpdateUserRequest struct {
	ID       int64  `json:"id" validate:"required,min=1"`
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
}

type UpdateUserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LogoutRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type LogoutResponse struct {
	Success bool `json:"success"`
}
