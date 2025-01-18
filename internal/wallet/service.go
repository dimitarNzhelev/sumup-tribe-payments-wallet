package wallet

import (
	"context"
	"tribe-payments-wallet-golang-interview-assignment/internal/config"
	"tribe-payments-wallet-golang-interview-assignment/internal/transactions"

	"github.com/google/uuid"
)

type WalletService struct {
	repo                WalletRepo
	transactionsService *transactions.TransactionService
}

func NewWalletService(repo WalletRepo, transactionsService *transactions.TransactionService) *WalletService {
	return &WalletService{repo: repo, transactionsService: transactionsService}
}

func (s *WalletService) CreateWallet(ctx context.Context, wallet *Wallet) error {
	if wallet == nil {
		return config.ErrWalletNotFound
	}

	if wallet.UserID == uuid.Nil {
		return config.ErrUserIDEmpty
	}

	err := s.repo.CreateWallet(ctx, wallet)
	if err != nil {
		return err
	}
	return nil
}

func (s *WalletService) GetWallet(ctx context.Context, id string) (*Wallet, error) {
	if id == "" {
		return nil, config.ErrWalletIDEmpty
	}

	wallet, err := s.repo.GetWallet(ctx, id)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *WalletService) updateWallet(ctx context.Context, wallet *Wallet) error {
	if wallet == nil {
		return config.ErrWalletNotFound
	}

	err := s.repo.UpdateWallet(ctx, wallet)
	if err != nil {
		return err
	}

	return nil
}

func (s *WalletService) DepositInWallet(ctx context.Context, money int64, wallet *Wallet) error {
	if money <= 0 {
		return config.ErrDepositNegative
	}

	if wallet == nil {
		return config.ErrWalletNotFound
	}

	wallet.Balance += money

	err := s.updateWallet(ctx, wallet)
	if err != nil {
		return err
	}

	tr := transactions.Transaction{
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

func (s *WalletService) WithdrawFromWallet(ctx context.Context, money int64, wallet *Wallet) error {
	if wallet == nil {
		return config.ErrWalletNotFound
	}

	if money <= 0 {
		return config.ErrWithdrawalZero
	}

	if wallet.Balance < money {
		return config.ErrInsufficientFunds
	}
	wallet.Balance -= money

	err := s.updateWallet(ctx, wallet)
	if err != nil {
		return err
	}

	tr := transactions.Transaction{
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
