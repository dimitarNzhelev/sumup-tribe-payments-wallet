package httpv1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tribe-payments-wallet-golang-interview-assignment/internal/auth"
	"tribe-payments-wallet-golang-interview-assignment/internal/config"
	walletModule "tribe-payments-wallet-golang-interview-assignment/internal/wallet"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/sumup-oss/go-pkgs/logger"
)

func writeWalletResponse(w http.ResponseWriter, statusCode int, wallet *walletModule.Wallet, log logger.StructuredLogger) {
	w.WriteHeader(statusCode)
	response := walletModule.WalletResponse{
		WalletID:  wallet.WalletID,
		Balance:   float64(wallet.Balance) / 100.0,
		Version:   wallet.Version,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error(fmt.Sprintf("Error encoding response: %s", err))
		http.Error(w, fmt.Sprintf("Error encoding response: %s", err), statusCode)
	}
}

func HandleCreateWallet(log logger.StructuredLogger, walletService *walletModule.WalletService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log.Info("HandleCreateWallet")

		// Get the user id from the context
		userID, ok := auth.GetUserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusUnauthorized)
			return
		}

		wallet := &walletModule.Wallet{
			UserID: uuid.MustParse(userID),
		}

		// Call the service to create the wallet
		err := walletService.CreateWallet(r.Context(), wallet)
		if err != nil {
			log.Error(fmt.Sprintf("Error creating wallet: %s", err))
			http.Error(w, fmt.Sprintf("Error creating wallet: %s", err), http.StatusUnprocessableEntity)
			return
		}

		writeWalletResponse(w, http.StatusCreated, wallet, log)
	}
}

func HandleGetWallet(log logger.StructuredLogger, walletService *walletModule.WalletService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get the wallet ID from the URL
		id := chi.URLParam(r, "id")
		if id == "" {
			log.Error(config.ErrWalletIDEmpty.Error())
			http.Error(w, config.ErrWalletIDEmpty.Error(), http.StatusBadRequest)
			return
		}

		log.Info(fmt.Sprintf("HandleGetWallet for wallet ID: %s", id))

		wallet, err := walletService.GetWallet(r.Context(), id)

		if err != nil {
			log.Error(fmt.Sprintf("Error getting wallet: %s", err))
			http.Error(w, "Wallet not found", http.StatusNotFound)
			return
		}

		writeWalletResponse(w, http.StatusOK, wallet, log)
	}
}

func HandleTransactionInWallet(log logger.StructuredLogger, walletService *walletModule.WalletService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Info("TransactionInWalletHandler")

		// Get the wallet ID from the URL
		id := chi.URLParam(r, "id")
		if id == "" {
			log.Error(config.ErrWalletIDEmpty.Error())
			http.Error(w, config.ErrWalletIDEmpty.Error(), http.StatusBadRequest)
			return
		}

		//Get the amount and type from the request body
		var req walletModule.WalletRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error(config.ErrInvalidRequestPayload.Error())
			http.Error(w, config.ErrInvalidRequestPayload.Error(), http.StatusBadRequest)
			return
		}

		log.Info(fmt.Sprintf("HandleTransactionInWallet for wallet ID: %s, amount: %f, transaction type %s", id, req.Amount, req.TransactionType))

		wallet, err := walletService.GetWallet(r.Context(), id)

		// Convert the amount to cents
		amount := int64(req.Amount * 100)

		if req.TransactionType == "withdraw" {
			err = walletService.WithdrawFromWallet(r.Context(), amount, wallet)
		} else if req.TransactionType == "deposit" {
			err = walletService.DepositInWallet(r.Context(), amount, wallet)
		} else {
			http.Error(w, "Invalid transaction type", http.StatusBadRequest)
			log.Error("Invalid transaction type")
			return
		}

		if err != nil {
			log.Error(fmt.Sprintf("Error %s in wallet: %s", req.TransactionType, err))
			http.Error(w, fmt.Sprintf("Error %s in wallet: %s", req.TransactionType, err), http.StatusUnprocessableEntity)
			return
		}

		writeWalletResponse(w, http.StatusOK, wallet, log)
	}
}
