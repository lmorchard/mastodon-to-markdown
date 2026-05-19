package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/lmorchard/mastodon-to-markdown/internal/timerange"
	"github.com/lmorchard/mastodon-to-markdown/internal/timewindow"
	"github.com/spf13/cobra"
)

// exportCmd is the orchestrator-facing entry point with the canonical
// `--since/--until/-o` flag shape shared across all *-to-markdown tools.
// It composes the same pipeline as `fetch`; filter / sort / visibility
// options stay in the user's config file rather than the command line.
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Orchestrator-friendly export of Mastodon posts to markdown",
	Long: `Fetch posts and render them to a markdown document over a canonical
time window.

The --since/--until flag shape matches the contract used by me-to-markdown
and the rest of the *-to-markdown tools. Filter options (--exclude-replies,
visibility, sort order, etc.) are read from the config file or environment;
this subcommand exposes only the orchestrator-facing flags.

Example usage:
  mastodon-to-markdown export --since 168h
  mastodon-to-markdown export --since 2026-05-11 --until 2026-05-18 -o posts.md`,
	RunE: func(cmd *cobra.Command, args []string) error {
		since, _ := cmd.Flags().GetString("since")
		until, _ := cmd.Flags().GetString("until")
		output, _ := cmd.Flags().GetString("output")

		now := time.Now()

		sinceTime, err := timewindow.Parse(since, now, false)
		if err != nil {
			return fmt.Errorf("--since: %w", err)
		}

		untilTime := now
		if until != "" {
			untilTime, err = timewindow.Parse(until, now, true)
			if err != nil {
				return fmt.Errorf("--until: %w", err)
			}
		}

		if !untilTime.After(sinceTime) {
			return fmt.Errorf("--until (%s) must be after --since (%s)",
				untilTime.Format(time.RFC3339), sinceTime.Format(time.RFC3339))
		}

		tr := &timerange.TimeRange{Start: sinceTime, End: untilTime}
		return runFetchPipeline(context.Background(), tr, output)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().String("since", "", "Start of time window (YYYY-MM-DD or Go duration like 168h) — required")
	exportCmd.Flags().String("until", "", "End of time window (YYYY-MM-DD, defaults to now)")
	exportCmd.Flags().StringP("output", "o", "", "Output file (default: stdout)")
	_ = exportCmd.MarkFlagRequired("since")
}
