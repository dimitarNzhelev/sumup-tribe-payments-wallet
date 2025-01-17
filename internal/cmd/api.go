package cmd

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"tribe-payments-wallet-golang-interview-assignment/internal/api"
	"tribe-payments-wallet-golang-interview-assignment/internal/config"
	"tribe-payments-wallet-golang-interview-assignment/internal/http"
	"tribe-payments-wallet-golang-interview-assignment/internal/transactions"
	"tribe-payments-wallet-golang-interview-assignment/internal/wallet"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"github.com/sumup-oss/go-pkgs/errors"
	"github.com/sumup-oss/go-pkgs/logger"
	"github.com/sumup-oss/go-pkgs/os"
	"github.com/sumup-oss/go-pkgs/task"
	"moul.io/chizap"
)

//nolint:gocognit
func NewApiCmd(osExecutor os.OsExecutor) *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "Run application server",
		Long:  "Run application server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			cfg, err := config.NewServerConfig()
			if err != nil {
				return errors.Wrap(err, "failed to create runtime config")
			}

			log, err := logger.NewZapLogger(
				logger.Configuration{
					Level:         cfg.Log.Level,
					Encoding:      logger.EncodingJSON,
					StdoutEnabled: cfg.Log.StdoutEnabled,
				},
			)
			if err != nil {
				return errors.Wrap(err, "failed to create logger")
			}

			defer log.Sync() //nolint:errcheck

			shutdownTask := newShutdownTask(
				log,
				osExecutor,
				cfg.GracefulShutdownTimeout,
			)

			connectionString := fmt.Sprintf(
				"postgres://%s:%s@%s/%s?sslmode=%s",
				cfg.Database.Username,
				cfg.Database.Password,
				cfg.Database.Host,
				cfg.Database.Database,
				cfg.Database.SSLMode,
			)

			log.Info(fmt.Sprintf("Connecting to database: %s", connectionString))

			db, err := sql.Open("postgres", connectionString)
			if err != nil {
				errors.New(fmt.Sprintf("Failed to connect to database: %s", err))
			}
			defer db.Close()

			err = db.Ping()
			if err != nil {
				errors.New(fmt.Sprintf("Failed to ping database: %s", err))
			}

			walletRepo := wallet.NewPostgresWalletRepo(db)
			transactionRepo := transactions.NewPostgresTransactionRepo(db)
			transactionService := transactions.NewTransactionService(transactionRepo)
			walletService := wallet.NewWalletService(walletRepo, transactionService)

			mux := chi.NewRouter()
			mux.Use(middleware.StripSlashes)
			mux.Use(
				http.Recovery(
					log,
					api.WritePanicResponse(log),
				),
				chizap.New(
					log.Logger,
					&chizap.Opts{
						WithReferer:   true,
						WithUserAgent: true,
					},
				),
			)

			api.RegisterRoutes(mux, log, walletService)

			httpServer := http.NewServer(
				log,
				cfg.ListenAddress,
				mux,
				http.WithMaxHeaderBytes(cfg.MaxHeaderBytes),
				http.WithReadHeaderTimeout(cfg.ReadHeaderTimeout),
				http.WithReadTimeout(cfg.ReadTimeout),
				http.WithServerShutdownTimeout(cfg.GracefulShutdownTimeout),
				http.WithWriteTimeout(cfg.WriteTimeout),
			)

			taskGroup := task.NewGroup()
			taskGroup.Go(
				shutdownTask.Run,
				httpServer.Run,
			)

			err = taskGroup.Wait(ctx)
			if err != nil {
				return errors.Wrap(err, "task group exits with error")
			}

			log.Info("taskGroup successfully shutdown")

			return nil
		},
	}
}
