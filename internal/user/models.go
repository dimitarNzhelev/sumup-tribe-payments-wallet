package user

import "github.com/google/uuid"

type UserStruct struct {
	ID           uuid.UUID
	FirstName    string
	LastName     string
	PasswordHash string
	Email        string
	CreatedAt    string
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCreateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
