package cmd

import (
	"context"
	"fmt"

	"github.com/lmorchard/mastodon-to-markdown/internal/mastodon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// validateAuthCmd is the machine-friendly companion to `whoami`:
// it answers "do my credentials work?" with a single-line stdout and
// a clean exit code, suitable for orchestrators (or scripts) that
// want to gate behavior on auth health.
var validateAuthCmd = &cobra.Command{
	Use:   "validate-auth",
	Short: "Check whether configured Mastodon credentials are accepted",
	Long: `Run a minimal authenticated request against the configured Mastodon
instance and exit 0 if the credentials are accepted, non-zero otherwise.

Reads ` + "`server`" + ` and ` + "`access_token`" + ` from the usual
sources (flags / MASTODON_SERVER + MASTODON_ACCESS_TOKEN env vars /
config file).`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig()
		cfg.Server = viper.GetString("server")
		cfg.AccessToken = viper.GetString("access_token")

		if cfg.Server == "" {
			return fmt.Errorf("mastodon server not configured (set MASTODON_SERVER or `server:`)")
		}
		if cfg.AccessToken == "" {
			return fmt.Errorf("mastodon access token not configured (set MASTODON_ACCESS_TOKEN or `access_token:`)")
		}

		client, err := mastodon.NewClient(cfg)
		if err != nil {
			return err
		}
		account, err := client.VerifyCredentials(context.Background())
		if err != nil {
			return fmt.Errorf("verify credentials: %w", err)
		}

		fmt.Printf("validate-auth: ok (@%s on %s)\n", account.Username, cfg.Server)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateAuthCmd)
}
