package cmd

import (
	"fmt"
	"os"

	"github.com/lmorchard/mastodon-to-markdown/internal/templates"
	"github.com/spf13/cobra"
)

const defaultConfigContent = `# Configuration file for mastodon-to-markdown
# Copy this to mastodon-to-markdown.yaml and customize as needed

# Database configuration
database: "mastodon-to-markdown.db"

# Logging configuration
verbose: false
debug: false
log_json: false

# Mastodon configuration
mastodon:
  # Your Mastodon instance URL (required)
  server: "https://mastodon.social"

  # Access token for authentication
  # Generate at: Settings > Development > New Application
  # Required scope: read:statuses
  access_token: ""

# Output configuration
output:
  # Include post metadata (timestamp, URL, visibility)
  include_metadata: true

  # Include URLs for media attachments
  include_media_urls: true

  # Template to use for output
  # Leave empty or omit to use built-in default template
  # Set to a filename to use a custom template file (e.g., "mastodon-to-markdown.md")
  template: ""

  # Sort order for posts: "asc" (oldest first, forward chronological) or "desc" (newest first)
  # Default: "asc"
  sort_order: "asc"
`

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration and template files",
	Long: `Create default configuration file and custom template file for customization.

This command generates:
  - mastodon-to-markdown.yaml (configuration file)
  - mastodon-to-markdown.md (customizable template, or use --template-file to specify)

Use --force to overwrite existing files.

Example:
  mastodon-to-markdown init
  mastodon-to-markdown init --template-file my-template.md
  mastodon-to-markdown init --force`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log := GetLogger()
		force, _ := cmd.Flags().GetBool("force")
		templateFile, _ := cmd.Flags().GetString("template-file")

		configFile := "mastodon-to-markdown.yaml"

		// Check if config file exists
		configExists := fileExists(configFile)
		if configExists && !force {
			return fmt.Errorf("config file %s already exists (use --force to overwrite)", configFile)
		}

		// Check if template file exists
		templateExists := fileExists(templateFile)
		if templateExists && !force {
			return fmt.Errorf("template file %s already exists (use --force to overwrite)", templateFile)
		}

		// Create config file
		if err := os.WriteFile(configFile, []byte(defaultConfigContent), 0644); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}

		if configExists {
			log.Infof("Overwrote %s", configFile)
		} else {
			log.Infof("Created %s", configFile)
		}

		// Get default template content
		templateContent, err := templates.GetDefaultTemplate()
		if err != nil {
			return fmt.Errorf("failed to get default template: %w", err)
		}

		// Create template file
		if err := os.WriteFile(templateFile, []byte(templateContent), 0644); err != nil {
			return fmt.Errorf("failed to create template file: %w", err)
		}

		if templateExists {
			log.Infof("Overwrote %s", templateFile)
		} else {
			log.Infof("Created %s", templateFile)
		}

		fmt.Printf("\n✅ Initialization complete!\n\n")
		fmt.Printf("Next steps:\n")
		fmt.Printf("  1. Edit %s and add your Mastodon server and access token\n", configFile)
		fmt.Printf("  2. (Optional) Customize %s for your preferred output format\n", templateFile)
		fmt.Printf("  3. Run: mastodon-to-markdown fetch --since 7d --output posts.md\n\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().Bool("force", false, "Overwrite existing files")
	initCmd.Flags().String("template-file", "mastodon-to-markdown.md", "Name of custom template file to create")
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
