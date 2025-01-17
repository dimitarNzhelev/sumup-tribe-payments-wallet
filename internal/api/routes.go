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
		r.Post("/wallet/{id}/withdraw", httpv1.WithdrawFromWalletHandler(log, walletService))
		r.Post("/wallet/{id}/deposit", httpv1.DepositInWalletHandler(log, walletService))
		r.Get("/wallet/{id}", httpv1.GetWalletHandler(log, walletService))
		r.Post("/wallet", httpv1.CreateWalletHandler(log, walletService))
	})

}
