package user

import (
	"context"
	"tribe-payments-wallet-golang-interview-assignment/internal/auth"

	"github.com/sumup-oss/go-pkgs/errors"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *UserStruct) error {
	if user == nil {
		return errors.New("User is nil")
	}

	_, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return errors.New("User already exists")
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

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*UserStruct, error) {
	if email == "" {
		return nil, errors.New("Email is empty")
	}

	usr, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*UserStruct, error) {
	if id == "" {
		return nil, errors.New("ID is empty")
	}

	usr, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return usr, nil
}
