package httpv1

import (
	"fmt"
	"net/http"
	"strings"
	"tribe-payments-wallet-golang-interview-assignment/internal/auth"
	"tribe-payments-wallet-golang-interview-assignment/internal/config"
	"tribe-payments-wallet-golang-interview-assignment/internal/wallet"

	"github.com/go-chi/chi/v5"
	"github.com/sumup-oss/go-pkgs/logger"
)

// Middleware that checks for a valid JWT
func AuthMiddleware(log logger.StructuredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, config.ErrInvalidAuthHeader.Error(), http.StatusUnauthorized)
				return
			}
			// Typically, the header is "Authorization: Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, config.ErrInvalidAuthHeader.Error(), http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]

			userID, _, err := auth.ValidateJWT(tokenString)
			if err != nil {
				log.Error(config.ErrInvalidToken.Error())
				http.Error(w, config.ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			// Store the user ID in the context so handlers can retrieve it
			ctx := auth.WriteUserIDToContext(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WalletOwnershipMiddleware(log logger.StructuredLogger, walletService *wallet.WalletService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			if id == "" {
				http.Error(w, config.ErrWalletIDEmpty.Error(), http.StatusBadRequest)
				return
			}

			walletObj, err := walletService.GetWallet(r.Context(), id)
			if err != nil {
				log.Error(config.ErrWalletNotFound.Error())
				http.Error(w, config.ErrWalletNotFound.Error(), http.StatusNotFound)
				return
			}

			// Retrieve userID from context, cause AuthMiddleware should be used before this
			userID, ok := auth.GetUserIDFromContext(r.Context())

			if !ok || userID == "" {
				// This scenario would only happen if AuthMiddleware wasn't used, or context wasn't set
				http.Error(w, config.ErrMissingUserIDContext.Error(), http.StatusUnauthorized)
				return
			}

			if walletObj.UserID.String() != userID {
				log.Error(fmt.Sprintf("Wallet %s does not belong to user %s", walletObj.WalletID.String(), userID))
				http.Error(w, "Wallet does not belong to user", http.StatusForbidden)
				return
			}

			// If all checks pass, call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
