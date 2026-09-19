package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	configPath string
	envPath    string
)

var rootCmd = &cobra.Command{
	Use:          "auction-bid",
	Short:        "Auction & bidding backend service",
	SilenceUsage: true,
}

func Execute() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		stop()
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "path ke file config")
	rootCmd.PersistentFlags().StringVarP(&envPath, "env", "e", ".env", "path ke file env")
}
