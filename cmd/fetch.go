package cmd

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/lmorchard/mastodon-to-markdown/internal/mastodon"
	"github.com/lmorchard/mastodon-to-markdown/internal/templates"
	"github.com/lmorchard/mastodon-to-markdown/internal/timerange"
	mastodonAPI "github.com/mattn/go-mastodon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch posts from Mastodon and export to markdown",
	Long: `Fetch posts from your Mastodon account for a specified time period
and export them to a markdown document suitable as a blog post starting point.

Example usage:
  mastodon-to-markdown fetch --since 7d --output posts.md
  mastodon-to-markdown fetch --start 2025-11-01 --end 2025-11-07
  mastodon-to-markdown fetch --since 24h --exclude-replies`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log := GetLogger()
		cfg := GetConfig()

		log.Info("Running fetch command")

		// Load Mastodon config from viper
		cfg.Mastodon.Server = viper.GetString("mastodon.server")
		cfg.Mastodon.AccessToken = viper.GetString("mastodon.access_token")
		cfg.Output.Template = viper.GetString("output.template")

		// Parse time range
		since := viper.GetString("fetch.since")
		start := viper.GetString("fetch.start")
		end := viper.GetString("fetch.end")

		tr, err := timerange.Parse(since, start, end)
		if err != nil {
			return fmt.Errorf("invalid time range: %w", err)
		}

		log.Infof("Fetching posts from %s to %s", timerange.FormatDate(tr.Start), timerange.FormatDate(tr.End))

		// Initialize Mastodon client
		client, err := mastodon.NewClient(cfg)
		if err != nil {
			return fmt.Errorf("failed to create Mastodon client: %w", err)
		}

		// Verify credentials and get account info
		ctx := context.Background()
		account, err := client.VerifyCredentials(ctx)
		if err != nil {
			return fmt.Errorf("failed to verify Mastodon credentials: %w", err)
		}

		log.Infof("Authenticated as @%s", account.Username)

		// Fetch statuses
		log.Info("Fetching statuses...")
		allStatuses := []*mastodonAPI.Status{}
		var maxID mastodonAPI.ID

		// Pagination loop
		for {
			pg := &mastodonAPI.Pagination{
				MaxID: maxID,
				Limit: 40,
			}

			statuses, err := client.GetStatuses(ctx, account.ID, pg)
			if err != nil {
				return fmt.Errorf("failed to fetch statuses: %w", err)
			}

			if len(statuses) == 0 {
				break
			}

			// Filter by time range
			foundInRange := false
			for _, status := range statuses {
				if status.CreatedAt.Before(tr.Start) {
					// We've gone past our time range
					break
				}
				if status.CreatedAt.After(tr.End) {
					// Haven't reached our time range yet
					continue
				}
				foundInRange = true
				allStatuses = append(allStatuses, status)
			}

			// If the last status is before our start time, we're done
			if len(statuses) > 0 && statuses[len(statuses)-1].CreatedAt.Before(tr.Start) {
				break
			}

			// If we didn't find any in range and we're past the end, keep going
			if !foundInRange && len(statuses) > 0 && statuses[len(statuses)-1].CreatedAt.Before(tr.End) {
				break
			}

			maxID = statuses[len(statuses)-1].ID
		}

		log.Infof("Found %d statuses in time range", len(allStatuses))

		// Apply filters
		filtered := filterStatuses(allStatuses,
			viper.GetBool("fetch.exclude_replies"),
			viper.GetBool("fetch.exclude_boosts"),
			viper.GetString("fetch.visibility"),
		)

		log.Infof("After filtering: %d statuses", len(filtered))

		// Convert to template format
		posts := mastodon.ConvertStatuses(filtered)

		// Sort posts based on configuration
		sortOrder := viper.GetString("output.sort_order")
		if sortOrder == "" {
			sortOrder = "asc" // Default to oldest first
		}
		sortPosts(posts, sortOrder)

		// Prepare template data
		data := &templates.TemplateData{
			StartDate: timerange.FormatDate(tr.Start),
			EndDate:   timerange.FormatDate(tr.End),
			Posts:     posts,
		}

		// Initialize template renderer
		templatePath := cfg.Output.Template
		renderer, err := templates.NewRenderer(templatePath)
		if err != nil {
			return fmt.Errorf("failed to initialize template: %w", err)
		}

		// Render to output
		outputFile := viper.GetString("fetch.output")
		if err := renderer.RenderToFile(outputFile, data); err != nil {
			return fmt.Errorf("failed to render output: %w", err)
		}

		if outputFile != "" && outputFile != "-" {
			log.Infof("Output written to %s", outputFile)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	// Time range flags
	fetchCmd.Flags().String("since", "", "Time period to fetch (e.g., '24h', '7d')")
	fetchCmd.Flags().String("start", "", "Start date (YYYY-MM-DD)")
	fetchCmd.Flags().String("end", "", "End date (YYYY-MM-DD)")

	// Output flags
	fetchCmd.Flags().StringP("output", "o", "", "Output file (default: stdout)")
	fetchCmd.Flags().String("sort-order", "asc", "Sort order: 'asc' (oldest first) or 'desc' (newest first)")

	// Filter flags
	fetchCmd.Flags().Bool("exclude-replies", false, "Exclude reply posts")
	fetchCmd.Flags().Bool("exclude-boosts", false, "Exclude boosted posts")
	fetchCmd.Flags().String("visibility", "", "Filter by visibility (comma-separated: public,unlisted,private)")

	// Bind flags to viper
	_ = viper.BindPFlag("fetch.since", fetchCmd.Flags().Lookup("since"))
	_ = viper.BindPFlag("fetch.start", fetchCmd.Flags().Lookup("start"))
	_ = viper.BindPFlag("fetch.end", fetchCmd.Flags().Lookup("end"))
	_ = viper.BindPFlag("fetch.output", fetchCmd.Flags().Lookup("output"))
	_ = viper.BindPFlag("output.sort_order", fetchCmd.Flags().Lookup("sort-order"))
	_ = viper.BindPFlag("fetch.exclude_replies", fetchCmd.Flags().Lookup("exclude-replies"))
	_ = viper.BindPFlag("fetch.exclude_boosts", fetchCmd.Flags().Lookup("exclude-boosts"))
	_ = viper.BindPFlag("fetch.visibility", fetchCmd.Flags().Lookup("visibility"))
}

// filterStatuses applies visibility, reply, and boost filters to statuses
func filterStatuses(statuses []*mastodonAPI.Status, excludeReplies, excludeBoosts bool, visibilityFilter string) []*mastodonAPI.Status {
	filtered := []*mastodonAPI.Status{}

	// Parse visibility filter
	visibilities := map[string]bool{}
	if visibilityFilter != "" {
		for _, v := range strings.Split(visibilityFilter, ",") {
			visibilities[strings.TrimSpace(v)] = true
		}
	}

	for _, status := range statuses {
		// Filter replies
		if excludeReplies && status.InReplyToID != nil {
			continue
		}

		// Filter boosts
		if excludeBoosts && status.Reblog != nil {
			continue
		}

		// Filter by visibility
		if len(visibilities) > 0 && !visibilities[string(status.Visibility)] {
			continue
		}

		filtered = append(filtered, status)
	}

	return filtered
}

// sortPosts sorts posts by creation time
// sortOrder: "asc" for oldest first (forward chronological), "desc" for newest first
func sortPosts(posts []templates.Post, sortOrder string) {
	sort.Slice(posts, func(i, j int) bool {
		if sortOrder == "desc" {
			return posts[i].CreatedAt.After(posts[j].CreatedAt)
		}
		// Default to "asc" - oldest first
		return posts[i].CreatedAt.Before(posts[j].CreatedAt)
	})
}
