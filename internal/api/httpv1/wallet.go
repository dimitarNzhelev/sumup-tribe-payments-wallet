package httpv1

import (
	"encoding/json"
	"fmt"
	"net/http"
	walletModule "tribe-payments-wallet-golang-interview-assignment/internal/wallet"

	"github.com/go-chi/chi/v5"

	"github.com/sumup-oss/go-pkgs/errors"
	"github.com/sumup-oss/go-pkgs/logger"
)

var (
	ErrInvalidPayload = errors.New("Invalid payload")
	ErrWalledIDEmpty  = errors.New("Wallet ID is required")
)

func CreateWalletHandler(log logger.StructuredLogger, walletService *walletModule.WalletService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log.Info("CreateWalletHandler")

		//Get the balance from the request body
		var req walletModule.WalletRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, ErrInvalidPayload.Error(), http.StatusBadRequest)
			return
		}

		wallet := &walletModule.WalletStruct{
			Balance: int64(req.Balance * 100),
		}

		// Call the service to create the wallet
		newWallet, err := walletService.CreateWallet(r.Context(), wallet)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating wallet: %s", err), http.StatusInternalServerError)
			return
		}

		// Set status and write the response as JSON
		w.WriteHeader(http.StatusCreated)

		response := walletModule.WalletResponse{
			WalletID:  newWallet.WalletID,
			Balance:   float64(newWallet.Balance) / 100.0,
			Version:   newWallet.Version,
			CreatedAt: newWallet.CreatedAt,
			UpdatedAt: newWallet.UpdatedAt,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %s", err), http.StatusInternalServerError)
			return
		}
	}
}

func GetWalletHandler(log logger.StructuredLogger, walletService *walletModule.WalletService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get the wallet ID from the URL
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, ErrWalledIDEmpty.Error(), http.StatusBadRequest)
			return
		}

		wallet, err := walletService.GetWallet(r.Context(), id)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting wallet: %s", err), http.StatusInternalServerError)
			return
		}

		// Set status and write the response as JSON
		w.WriteHeader(http.StatusOK)
		response := walletModule.WalletResponse{
			WalletID:  wallet.WalletID,
			Balance:   float64(wallet.Balance) / 100.0,
			Version:   wallet.Version,
			CreatedAt: wallet.CreatedAt,
			UpdatedAt: wallet.UpdatedAt,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %s", err), http.StatusInternalServerError)
			return
		}
	}
}

func TransactionInWalletHandler(log logger.StructuredLogger, walletService *walletModule.WalletService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Info("TransactionInWalletHandler")

		// Get the wallet ID from the URL
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, ErrWalledIDEmpty.Error(), http.StatusBadRequest)
			return
		}

		//Get the amount from the request body
		var req walletModule.WalletRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, ErrInvalidPayload.Error(), http.StatusBadRequest)
			return
		}

		wallet, err := walletService.GetWallet(r.Context(), id)

		// Convert the amount to cents
		amount := int64(req.Amount * 100)

		if req.TransactionType == "withdraw" {
			err = walletService.WithdrawFromWallet(r.Context(), amount, wallet)
		} else if req.TransactionType == "deposit" {
			err = walletService.DepositInWallet(r.Context(), amount, wallet)
		} else {
			http.Error(w, "Invalid transaction type", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("Error %s in wallet: %s", req.TransactionType, err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := walletModule.WalletResponse{
			WalletID:  wallet.WalletID,
			Balance:   float64(wallet.Balance) / 100.0,
			Version:   wallet.Version,
			CreatedAt: wallet.CreatedAt,
			UpdatedAt: wallet.UpdatedAt,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %s", err), http.StatusInternalServerError)
			return
		}
	}
}
