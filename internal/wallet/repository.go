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
	CreateWallet(ctx context.Context, wallet *Wallet) error
	GetWallet(ctx context.Context, id string) (*Wallet, error)
	UpdateWallet(ctx context.Context, wallet *Wallet) error
}

func (r *PostgresWalletRepo) CreateWallet(ctx context.Context, wallet *Wallet) error {
	// Updated query with RETURNING clause to retrieve the new wallet ID
	query := `INSERT INTO wallets (user_id) VALUES ($1) RETURNING id, version, created_at, updated_at`

	// Declare a variable to store the generated wallet ID
	var walletID string

	// Execute the query and retrieve the new wallet ID
	err := r.db.QueryRowContext(ctx, query, wallet.UserID).Scan(&walletID, &wallet.Version, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return err
	}

	wallet.WalletID = uuid.MustParse(walletID)

	return nil
}

func (r *PostgresWalletRepo) GetWallet(ctx context.Context, id string) (*Wallet, error) {
	query := `SELECT * FROM wallets WHERE id = $1`

	wallet := Wallet{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&wallet.WalletID, &wallet.UserID, &wallet.Balance, &wallet.Version, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *PostgresWalletRepo) UpdateWallet(ctx context.Context, wallet *Wallet) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted, // TODO: Adjust isolation level as needed
	})
	if err != nil {
		return err
	}

	// Handle the commit/rollback at the end using a defer
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Lock the wallet row
	var currentVersion int
	lockQuery := `SELECT version FROM wallets WHERE id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, lockQuery, wallet.WalletID).Scan(&currentVersion)
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
