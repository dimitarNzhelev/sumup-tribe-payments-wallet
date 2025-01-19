package config

import "github.com/sumup-oss/go-pkgs/errors"

var (
	ErrInvalidRequestPayload   = errors.New("Invalid request payload")
	ErrWalletIDEmpty           = errors.New("Wallet ID is required")
	ErrInvalidCredentials      = errors.New("Invalid credentials")
	ErrWalletNotFound          = errors.New("Wallet not found")
	ErrBalanceNegative         = errors.New("Balance cannot be negative")
	ErrUserIDEmpty             = errors.New("User ID is empty")
	ErrDepositNegative         = errors.New("Deposit amount must be positive")
	ErrWithdrawalZero          = errors.New("Withdrawal amount must be positive")
	ErrInsufficientFunds       = errors.New("Insufficient funds")
	ErrEmptyField              = errors.New("Empty field in user")
	ErrUserExists              = errors.New("User already exists")
	ErrUserIsNil               = errors.New("User is nil")
	ErrTransactionNotFound     = errors.New("Transaction not found")
	ErrInvalidAuthHeader       = errors.New("Invalid Authorization header format")
	ErrInvalidToken            = errors.New("Invalid JWT token")
	ErrMissingUserIDContext    = errors.New("Missing user ID in context")
	ErrFailedToCreateUser      = errors.New("Failed to create user")
	ErrGeneratingJWT           = errors.New("Error generating JWT")
	ErrInvalidClaimsFormat     = errors.New("Invalid claims format")
	ErrUnexpectedSigningMethod = errors.New("Unexpected signing method")
	ErrInvalidUserData         = errors.New("Invalid user data")
)
