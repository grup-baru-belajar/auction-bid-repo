package cmd

import (
	"github.com/grup-baru-belajar/auction-bid-repo/internal/config"
	"github.com/spf13/cobra"
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
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
