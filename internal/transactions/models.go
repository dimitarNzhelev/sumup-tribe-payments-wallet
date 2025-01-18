package transactions

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id              uuid.UUID
	WalletID        uuid.UUID
	Amount          int64
	TransactionType string
	BalanceSnapshot int64
	Created_At      time.Time
}
