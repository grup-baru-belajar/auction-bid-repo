package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/config"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/database"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/handlers"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/middlewares"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/repository"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/routes"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/services"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/token"
	wsh "github.com/grup-baru-belajar/auction-bid-repo/internal/websocket"
)

const (
	shutdownTimeout   = 10 * time.Second
	readHeaderTimeout = 10 * time.Second
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Jalankan HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(
			configPath,
			envPath,
			cmd.Flags().Changed("config"),
			cmd.Flags().Changed("env"),
		)
		if err != nil {
			return err
		}

		cmd.Printf("app.env=%s app.port=%d\n", cfg.App.Env, cfg.App.Port)
		cmd.Printf(
			"database.host=%s database.port=%d database.user=%s database.name=%s database.sslmode=%s\n",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Name, cfg.Database.SSLMode,
		)
		cmd.Printf("jwt.expires_in=%s\n", cfg.JWT.ExpiresIn)

		db, err := database.New(cmd.Context(), cfg.Database.DSN())
		if err != nil {
			return err
		}
		defer db.Close()

		cmd.Println("database: connected")

		userRepo := repository.NewUserRepository(db)
		auctionRepo := repository.NewAuctionRepository(db)
		auctionDetailRepo := repository.NewAuctionDetailRepository(db)
		bidRepo := repository.NewBidRepository(db)
		reportingRepo := repository.NewReportingRepository(db)

		tokenManager := token.NewTokenManager(cfg.JWT.Secret, cfg.JWT.ExpiresIn)
		authService := services.NewAuthService(userRepo, tokenManager)
		cloudinaryService := services.NewCloudinaryService(cfg.Cloudinary)
		auctionService := services.NewAuctionService(auctionRepo)
		auctionDetailService := services.NewAuctionDetailService(auctionDetailRepo)
		bidService := services.NewBidService(bidRepo)
		reportingService := services.NewReportingService(reportingRepo)

		wsHub := wsh.NewHub()
		wsHandler := wsh.NewHandler(wsHub, auctionDetailRepo)

		handler := handlers.New(
			authService,
			auctionService,
			auctionDetailService,
			bidService,
			reportingService,
			wsHandler,
			cloudinaryService,
		)

		if cfg.App.Env != "development" {
			gin.SetMode(gin.ReleaseMode)
		}
		router := gin.Default()
		router.Use(middlewares.CORS(cfg.App.CORSOrigins))
		routes.Setup(router, handler, tokenManager, wsHandler)

		srv := &http.Server{
			Addr:              fmt.Sprintf(":%d", cfg.App.Port),
			Handler:           router,
			ReadHeaderTimeout: readHeaderTimeout,
		}

		errCh := make(chan error, 1)
		go func() {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}()
		cmd.Printf("http: listening on %s\n", srv.Addr)

		select {
		case err := <-errCh:
			return fmt.Errorf("http server: %w", err)
		case <-cmd.Context().Done():
			cmd.Println("http: shutting down...")
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("http shutdown: %w", err)
		}
		cmd.Println("http: stopped")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
