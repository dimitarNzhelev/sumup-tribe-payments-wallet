package user

import (
	"context"
	"tribe-payments-wallet-golang-interview-assignment/internal/auth"

	"github.com/sumup-oss/go-pkgs/errors"
)

var (
	ErrEmptyField = errors.New("Empty field in user")
	ErrUserExists = errors.New("User already exists")
	ErrUserIsNil  = errors.New("User is nil")
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
	if user == nil {
		return ErrUserIsNil
	}

	if user.Email == "" || user.PasswordHash == "" || user.FirstName == "" || user.LastName == "" {
		return ErrEmptyField
	}

	_, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return ErrUserExists
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
		return nil, ErrEmptyField
	}

	usr, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*User, error) {
	if id == "" {
		return nil, ErrEmptyField
	}

	usr, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return usr, nil
}
