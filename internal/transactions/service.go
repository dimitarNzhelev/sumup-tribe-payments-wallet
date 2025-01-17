package transactions

import (
	"context"

	"github.com/sumup-oss/go-pkgs/errors"
)

type TransactionService struct {
	repo TransactionRepo
}

func NewTransactionService(repo TransactionRepo) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, transaction *TransactionStruct) error {
	if transaction == nil {
		return errors.New("Transaction is nil")
	}

	err := s.repo.CreateTransaction(ctx, transaction)
	if err != nil {
		return err
	}
	return nil
}
