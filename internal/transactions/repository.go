package transactions

import (
	"context"
	"database/sql"
)

type PostgresTransactionRepo struct {
	db *sql.DB
}

func NewPostgresTransactionRepo(db *sql.DB) *PostgresTransactionRepo {
	return &PostgresTransactionRepo{db: db}
}

type TransactionRepo interface {
	CreateTransaction(ctx context.Context, transaction *Transaction) error
}

func (r *PostgresTransactionRepo) CreateTransaction(ctx context.Context, transaction *Transaction) error {
	query := `INSERT INTO transactions ( wallet_id, amount, transaction_type, balance_snapshot) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, transaction.WalletID, transaction.Amount, transaction.TransactionType, transaction.BalanceSnapshot)
	if err != nil {
		return err
	}
	return nil
}
