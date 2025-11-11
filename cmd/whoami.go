package cmd

import (
	"context"
	"fmt"

	"github.com/lmorchard/mastodon-to-markdown/internal/mastodon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// whoamiCmd represents the whoami command
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show information about the authenticated Mastodon account",
	Long: `Display details about the Mastodon account associated with the configured
access token. This is useful for verifying your configuration and credentials.

Example:
  mastodon-to-markdown whoami`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log := GetLogger()
		cfg := GetConfig()

		// Load Mastodon config from viper
		cfg.Mastodon.Server = viper.GetString("mastodon.server")
		cfg.Mastodon.AccessToken = viper.GetString("mastodon.access_token")

		// Validate configuration
		if cfg.Mastodon.Server == "" {
			return fmt.Errorf("mastodon server not configured (set mastodon.server in config file)")
		}
		if cfg.Mastodon.AccessToken == "" {
			return fmt.Errorf("mastodon access token not configured (set mastodon.access_token in config file)")
		}

		// Initialize Mastodon client
		client, err := mastodon.NewClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create Mastodon client: %w", err)
		}

		// Verify credentials and get account info
		ctx := context.Background()
		account, err := client.VerifyCredentials(ctx)
		if err != nil {
			return fmt.Errorf("failed to verify credentials: %w", err)
		}

		// Display account information
		fmt.Printf("\n✅ Successfully authenticated!\n\n")
		fmt.Printf("Server:        %s\n", cfg.Mastodon.Server)
		fmt.Printf("Username:      @%s\n", account.Username)
		fmt.Printf("Display Name:  %s\n", account.DisplayName)
		fmt.Printf("Account ID:    %s\n", account.ID)
		fmt.Printf("Account URL:   %s\n", account.URL)
		fmt.Printf("Created:       %s\n", account.CreatedAt.Format("2006-01-02"))
		fmt.Printf("\nStatistics:\n")
		fmt.Printf("  Posts:       %d\n", account.StatusesCount)
		fmt.Printf("  Following:   %d\n", account.FollowingCount)
		fmt.Printf("  Followers:   %d\n", account.FollowersCount)
		fmt.Printf("\n")

		log.Debug("Account verification successful")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
