package transactions

import (
	"context"
	"tribe-payments-wallet-golang-interview-assignment/internal/config"
)

type TransactionService struct {
	repo TransactionRepo
}

func NewTransactionService(repo TransactionRepo) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, transaction *Transaction) error {
	if transaction == nil {
		return config.ErrTransactionNotFound
	}

	err := s.repo.CreateTransaction(ctx, transaction)
	if err != nil {
		return err
	}
	return nil
}
