package wallet

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type PostgresWalletRepo struct {
	db *sql.DB
}

func NewPostgresWalletRepo(db *sql.DB) *PostgresWalletRepo {
	return &PostgresWalletRepo{db: db}
}

type WalletRepo interface {
	CreateWallet(ctx context.Context, wallet *WalletStruct) (*WalletStruct, error)
	GetWallet(ctx context.Context, id string) (*WalletStruct, error)
	UpdateWallet(ctx context.Context, wallet *WalletStruct) error
}

func (r *PostgresWalletRepo) CreateWallet(ctx context.Context, wallet *WalletStruct) (*WalletStruct, error) {
	// Updated query with RETURNING clause to retrieve the new wallet ID
	query := `INSERT INTO wallets (balance) VALUES ($1) RETURNING id`

	// Declare a variable to store the generated wallet ID
	var walletID string

	// Execute the query and retrieve the new wallet ID
	err := r.db.QueryRowContext(ctx, query, wallet.Balance).Scan(&walletID)
	if err != nil {
		return nil, err
	}

	// Return the newly created wallet
	return &WalletStruct{
		WalletID: uuid.MustParse(walletID),
		Balance:  wallet.Balance,
	}, nil
}

func (r *PostgresWalletRepo) GetWallet(ctx context.Context, id string) (*WalletStruct, error) {
	query := `SELECT * FROM wallets WHERE id = $1`
	wallet := WalletStruct{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&wallet.WalletID, &wallet.Balance, &wallet.Version, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		fmt.Println("Error getting wallet", err)
		return nil, err
	}
	return &wallet, nil
}

func (r *PostgresWalletRepo) UpdateWallet(ctx context.Context, wallet *WalletStruct) error {
	query := `UPDATE wallets SET balance = $1, updated_at = NOW(), version = version + 1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, wallet.Balance, wallet.WalletID)
	if err != nil {
		return err
	}
	return nil
}
