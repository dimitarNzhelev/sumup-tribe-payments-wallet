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
	query := `INSERT INTO wallets (balance, user_id) VALUES ($1, $2) RETURNING id, created_at, updated_at`

	// Declare a variable to store the generated wallet ID
	var walletID string

	// Execute the query and retrieve the new wallet ID
	err := r.db.QueryRowContext(ctx, query, wallet.Balance, wallet.UserID).Scan(&walletID, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Return the newly created wallet
	return &WalletStruct{
		WalletID:  uuid.MustParse(walletID),
		UserID:    wallet.UserID,
		Balance:   wallet.Balance,
		Version:   1,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}, nil
}

func (r *PostgresWalletRepo) GetWallet(ctx context.Context, id string) (*WalletStruct, error) {
	query := `SELECT * FROM wallets WHERE id = $1`
	wallet := WalletStruct{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&wallet.WalletID, &wallet.UserID, &wallet.Balance, &wallet.Version, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *PostgresWalletRepo) UpdateWallet(ctx context.Context, wallet *WalletStruct) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted, // TODO: Adjust isolation level as needed
	})
	if err != nil {
		return err
	}

	// Handle the commit/rollback at the end using a defer
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Lock the wallet row
	var currentBalance float64
	var currentVersion int
	lockQuery := `SELECT balance, version FROM wallets WHERE id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, lockQuery, wallet.WalletID).Scan(&currentBalance, &currentVersion)
	if err != nil {
		return err
	}

	// Check if the version is still the same
	if wallet.Version != currentVersion {
		return fmt.Errorf("Wallet version mismatch")
	}

	// Perform the update
	updateQuery := `
        UPDATE wallets
        SET balance = $1,
            updated_at = NOW(),
            version = version + 1
        WHERE id = $2
    `
	_, err = tx.ExecContext(ctx, updateQuery, wallet.Balance, wallet.WalletID)
	if err != nil {
		return err
	}

	// The transaction will be committed by the deferred function if no errors occurred
	return nil
}
