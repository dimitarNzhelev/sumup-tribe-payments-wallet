package httpv1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tribe-payments-wallet-golang-interview-assignment/internal/wallet"

	"github.com/go-chi/chi/v5"

	"github.com/sumup-oss/go-pkgs/logger"
)

func CreateWalletHandler(log logger.StructuredLogger, walletService *wallet.WalletService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log.Info("CreateWalletHandler")

		//Get the balance from the request body
		var req wallet.WalletRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		wallet := &wallet.WalletStruct{
			Balance: req.Balance,
		}

		// Call the service to create the wallet
		newWallet, err := walletService.CreateWallet(r.Context(), wallet)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating wallet: %s", err), http.StatusInternalServerError)
			return
		}

		// Set status and write the response as JSON
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(newWallet); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %s", err), http.StatusInternalServerError)
			return
		}
	}
}

func GetWalletHandler(log logger.StructuredLogger, walletService *wallet.WalletService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get the wallet ID from the URL
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "Wallet ID is required", http.StatusBadRequest)
			return
		}

		wallet, err := walletService.GetWallet(r.Context(), id)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting wallet: %s", err), http.StatusInternalServerError)
			return
		}

		// Set status and write the response as JSON
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(wallet); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %s", err), http.StatusInternalServerError)
			return
		}
	}
}

func DepositInWalletHandler(log logger.StructuredLogger, walletService *wallet.WalletService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Info("DepositInWalletHandler")

		// Get the wallet ID from the URL
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "Wallet ID is required", http.StatusBadRequest)
			return
		}

		//Get the amount from the request body
		var req wallet.WalletRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid request payload: %s", err), http.StatusBadRequest)
			return
		}

		wallet, err := walletService.GetWallet(r.Context(), id)

		// Call the service to deposit in the wallet
		err = walletService.DepositInWallet(r.Context(), req.Amount, wallet)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error depositing in wallet: %s", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(wallet); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %s", err), http.StatusInternalServerError)
			return
		}
	}
}

func WithdrawFromWalletHandler(log logger.StructuredLogger, walletService *wallet.WalletService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Info("WithdrawFromWalletHandler")

		// Get the wallet ID from the URL
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "Wallet ID is required", http.StatusBadRequest)
			return
		}

		//Get the amount from the request body
		var req wallet.WalletRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid request payload: %s", err), http.StatusBadRequest)
			return
		}

		wallet, err := walletService.GetWallet(r.Context(), id)

		// Call the service to deposit in the wallet
		err = walletService.WithdrawFromWallet(r.Context(), req.Amount, wallet)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error withdrawing from wallet: %s", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(wallet); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %s", err), http.StatusInternalServerError)
			return
		}
	}
}
