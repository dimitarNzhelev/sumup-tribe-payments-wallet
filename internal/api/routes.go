package api

import (
	"tribe-payments-wallet-golang-interview-assignment/internal/api/httpv1"
	"tribe-payments-wallet-golang-interview-assignment/internal/wallet"

	"github.com/go-chi/chi/v5"
	"github.com/sumup-oss/go-pkgs/logger"
)

func RegisterRoutes(
	mux *chi.Mux,
	log logger.StructuredLogger,
	walletService *wallet.WalletService,
) {
	mux.Get("/live", Health)

	mux.Route("/v1", func(r chi.Router) {
		r.Post("/wallet/{id}/transaction", httpv1.HandleTransactionInWallet(log, walletService))
		r.Get("/wallet/{id}", httpv1.HandleGetWallet(log, walletService))
		r.Post("/wallet", httpv1.HandleCreateWallet(log, walletService))
	})

}
