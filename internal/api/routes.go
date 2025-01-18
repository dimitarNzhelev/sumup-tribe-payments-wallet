package api

import (
	"tribe-payments-wallet-golang-interview-assignment/internal/api/httpv1"
	"tribe-payments-wallet-golang-interview-assignment/internal/config"
	"tribe-payments-wallet-golang-interview-assignment/internal/user"
	"tribe-payments-wallet-golang-interview-assignment/internal/wallet"

	"github.com/go-chi/chi/v5"
	"github.com/sumup-oss/go-pkgs/logger"
)

func RegisterRoutes(
	mux *chi.Mux,
	log logger.StructuredLogger,
	walletService *wallet.WalletService,
	userService *user.UserService,
	Auth config.JWTConfig,
) {
	mux.Get("/live", Health)

	mux.Route("/v1", func(r chi.Router) {
		r.Post("/user", httpv1.HandlerCreateUser(log, userService))
		r.Post("/login", httpv1.HandlerLoginUser(log, userService))

		// Everything here requires a valid JWT token
		r.Group(func(r chi.Router) {
			r.Use(httpv1.AuthMiddleware(log))

			// Create wallet
			r.Post("/wallet", httpv1.HandleCreateWallet(log, walletService))

			// Access or modify a specific wallet only if user owns it
			r.Route("/wallet/{id}", func(r chi.Router) {
				r.Use(httpv1.WalletOwnershipMiddleware(log, walletService))

				r.Get("/", httpv1.HandleGetWallet(log, walletService))
				r.Post("/transaction", httpv1.HandleTransactionInWallet(log, walletService))
			})
		})
	})

}
