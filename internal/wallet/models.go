package wallet

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	WalletID  uuid.UUID
	UserID    uuid.UUID
	Balance   int64
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WalletRequest struct {
	Balance         float64 `json:"balance"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
}

type WalletResponse struct {
	WalletID  uuid.UUID `json:"wallet_id"`
	Balance   float64   `json:"balance"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
