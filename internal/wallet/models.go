package wallet

import (
	"time"

	"github.com/google/uuid"
)

type WalletStruct struct {
	WalletID  uuid.UUID
	Balance   int64
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WalletRequest struct {
	Balance float64 `json:"balance"`
	Amount  float64 `json:"amount"`
}

type WalletResponse struct {
	WalletID  uuid.UUID `json:"wallet_id"`
	Balance   float64   `json:"balance"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
