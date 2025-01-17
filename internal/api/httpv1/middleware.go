package httpv1

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"tribe-payments-wallet-golang-interview-assignment/internal/auth"
	"tribe-payments-wallet-golang-interview-assignment/internal/wallet"

	"github.com/go-chi/chi/v5"
	"github.com/sumup-oss/go-pkgs/logger"
)

// Middleware that checks for a valid JWT
func CheckAuth(log logger.StructuredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}
			// Typically, the header is "Authorization: Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]

			userID, _, err := auth.ValidateJWT(tokenString)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid token: %s", err), http.StatusUnauthorized)
				return
			}

			// Store the user ID in the context so handlers can retrieve it
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CheckIfWalletIsOwnedByUser(log logger.StructuredLogger, walletService *wallet.WalletService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			if id == "" {
				http.Error(w, "Wallet ID is required", http.StatusBadRequest)
				return
			}

			walletObj, err := walletService.GetWallet(r.Context(), id)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error getting wallet: %s", err), http.StatusInternalServerError)
				return
			}

			// Retrieve userID from context, cause CheckAuth should be used before this
			userID, ok := r.Context().Value("userID").(string)

			if !ok || userID == "" {
				// This scenario would only happen if CheckAuth wasn't used, or context wasn't set
				http.Error(w, "Missing user in context", http.StatusUnauthorized)
				return
			}

			if walletObj.UserID.String() != userID {
				http.Error(w, "Wallet does not belong to user", http.StatusForbidden)
				return
			}

			// If all checks pass, call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
