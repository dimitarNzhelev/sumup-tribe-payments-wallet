package transactions

import (
	"time"

	"github.com/google/uuid"
)

type TransactionStruct struct {
	Id              uuid.UUID
	WalletID        uuid.UUID
	Amount          float64
	TransactionType string
	BalanceSnapshot float64
	Created_At      time.Time
}
