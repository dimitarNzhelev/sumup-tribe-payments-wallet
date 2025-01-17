package wallet

import (
	"context"
	"tribe-payments-wallet-golang-interview-assignment/internal/transactions"

	"github.com/sumup-oss/go-pkgs/errors"
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
		return nil, errors.New("Wallet is nil")
	}

	if wallet.Balance < 0 {
		return nil, errors.New("Balance cannot be negative")
	}

	wallet, err := s.repo.CreateWallet(ctx, wallet)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *WalletService) GetWallet(ctx context.Context, id string) (*WalletStruct, error) {
	if id == "" {
		return nil, errors.New("Wallet ID is empty")
	}

	wallet, err := s.repo.GetWallet(ctx, id)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *WalletService) updateWallet(ctx context.Context, wallet *WalletStruct) error {
	if wallet == nil {
		return errors.New("Wallet is nil")
	}

	err := s.repo.UpdateWallet(ctx, wallet)
	if err != nil {
		return err
	}

	return nil
}

func (s *WalletService) DepositInWallet(ctx context.Context, money int64, wallet *WalletStruct) error {
	if money <= 0 {
		return errors.New("Deposit amount must be positive")
	}

	if wallet == nil {
		return errors.New("Wallet not found")
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
		return errors.New("Wallet not found")
	}

	if money <= 0 {
		return errors.New("Withdrawal amount must be positive")
	}

	if wallet.Balance < money {
		return errors.New("Insufficient funds")
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
