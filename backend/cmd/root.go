package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	configPath string
	envPath    string
)

var rootCmd = &cobra.Command{
	Use:   "auction-bid",
	Short: "Auction & bidding backend service",
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "path ke file config")
	rootCmd.PersistentFlags().StringVarP(&envPath, "env", "e", ".env", "path ke file env")
}
