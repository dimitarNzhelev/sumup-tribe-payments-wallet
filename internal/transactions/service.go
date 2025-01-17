package transactions

import (
	"context"
)

type TransactionService struct {
	repo TransactionRepo
}

func NewTransactionService(repo TransactionRepo) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, transaction *TransactionStruct) error {
	err := s.repo.CreateTransaction(ctx, transaction)
	if err != nil {
		return err
	}
	return nil
}
