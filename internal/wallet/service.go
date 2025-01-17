package wallet

import (
	"context"
	"tribe-payments-wallet-golang-interview-assignment/internal/transactions"

	"github.com/sumup-oss/go-pkgs/errors"
)

var (
	ErrWalletNotFound    = errors.New("Wallet not found")
	ErrBalanceNegative   = errors.New("Balance cannot be negative")
	ErrWalledIDEmpty     = errors.New("Wallet ID is empty")
	ErrDepositNegative   = errors.New("Deposit amount must be positive")
	ErrWithdrawalZero    = errors.New("Withdrawal amount must be positive")
	ErrInsufficientFunds = errors.New("Insufficient funds")
)

type WalletService struct {
	repo                WalletRepo
	transactionsService *transactions.TransactionService
}

func NewWalletService(repo WalletRepo, transactionsService *transactions.TransactionService) *WalletService {
	return &WalletService{repo: repo, transactionsService: transactionsService}
}

func (s *WalletService) CreateWallet(ctx context.Context, wallet *WalletStruct) (*WalletStruct, error) {
	if wallet == nil {
		return nil, ErrWalletNotFound
	}

	if wallet.Balance < 0 {
		return nil, ErrBalanceNegative
	}

	wallet, err := s.repo.CreateWallet(ctx, wallet)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *WalletService) GetWallet(ctx context.Context, id string) (*WalletStruct, error) {
	if id == "" {
		return nil, ErrWalledIDEmpty
	}

	wallet, err := s.repo.GetWallet(ctx, id)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *WalletService) updateWallet(ctx context.Context, wallet *WalletStruct) error {
	if wallet == nil {
		return ErrWalletNotFound
	}

	err := s.repo.UpdateWallet(ctx, wallet)
	if err != nil {
		return err
	}

	return nil
}

func (s *WalletService) DepositInWallet(ctx context.Context, money int64, wallet *WalletStruct) error {
	if money <= 0 {
		return ErrDepositNegative
	}

	if wallet == nil {
		return ErrWalletNotFound
	}

	wallet.Balance += money

	err := s.updateWallet(ctx, wallet)
	if err != nil {
		return err
	}

	tr := transactions.TransactionStruct{
		WalletID:        wallet.WalletID,
		Amount:          money,
		TransactionType: "deposit",
		BalanceSnapshot: wallet.Balance,
	}

	err = s.transactionsService.CreateTransaction(ctx, &tr)
	if err != nil {
		return err
	}
	return nil
}

func (s *WalletService) WithdrawFromWallet(ctx context.Context, money int64, wallet *WalletStruct) error {
	if wallet == nil {
		return ErrWalletNotFound
	}

	if money <= 0 {
		return ErrWithdrawalZero
	}

	if wallet.Balance < money {
		return ErrInsufficientFunds
	}
	wallet.Balance -= money

	err := s.updateWallet(ctx, wallet)
	if err != nil {
		return err
	}

	tr := transactions.TransactionStruct{
		WalletID:        wallet.WalletID,
		Amount:          money,
		TransactionType: "withdrawal",
		BalanceSnapshot: wallet.Balance,
	}

	err = s.transactionsService.CreateTransaction(ctx, &tr)
	if err != nil {
		return err
	}
	return nil
}
