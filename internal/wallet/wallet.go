package wallet

import (
	"time"

	"github.com/google/uuid"
)

type WalletStruct struct {
	WalletID  uuid.UUID
	Balance   float64
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WalletRequest struct {
	Balance float64 `json:"balance"`
	Amount  float64 `json:"amount"`
}
