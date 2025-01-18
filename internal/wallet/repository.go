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
	CreateWallet(ctx context.Context, wallet *Wallet) (*Wallet, error)
	GetWallet(ctx context.Context, id string) (*Wallet, error)
	UpdateWallet(ctx context.Context, wallet *Wallet) error
}

func (r *PostgresWalletRepo) CreateWallet(ctx context.Context, wallet *Wallet) (*Wallet, error) {
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
	return &Wallet{
		WalletID: uuid.MustParse(walletID),
		Balance:  wallet.Balance,
	}, nil
}

func (r *PostgresWalletRepo) GetWallet(ctx context.Context, id string) (*Wallet, error) {
	query := `SELECT * FROM wallets WHERE id = $1`
	wallet := Wallet{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&wallet.WalletID, &wallet.Balance, &wallet.Version, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		fmt.Println("Error getting wallet", err)
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
	var currentBalance int64
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

	// Check if the balance is still the same
	if wallet.Balance != currentBalance {
		return fmt.Errorf("Wallet balance mismatch")
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
