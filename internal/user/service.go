package user

import (
	"context"
	"tribe-payments-wallet-golang-interview-assignment/internal/auth"
	"tribe-payments-wallet-golang-interview-assignment/internal/config"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
	if user == nil {
		return config.ErrUserIsNil
	}

	if user.Email == "" || user.PasswordHash == "" || user.FirstName == "" || user.LastName == "" {
		return config.ErrEmptyField
	}

	_, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return config.ErrUserExists
	}

	user.PasswordHash, err = auth.GeneratePasswordHash(user.PasswordHash)
	if err != nil {
		return err
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	if email == "" {
		return nil, config.ErrEmptyField
	}

	usr, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*User, error) {
	if id == "" {
		return nil, config.ErrEmptyField
	}

	usr, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return usr, nil
}
