# mastodon-to-markdown

Export your Mastodon posts to markdown for blog post creation and content curation.

## Overview

`mastodon-to-markdown` is a command-line tool that fetches posts from your Mastodon account and exports them to well-formatted markdown documents. It's perfect as a starting point for blog posts, content curation, or archiving your thoughts.

## Features

- **Flexible Time Ranges**: Fetch posts from the last N hours/days/weeks, or specify exact date ranges
- **Smart Filtering**: Exclude replies, boosts, or private posts
- **Favorites & Boosts**: Include posts you've favorited and boosted, organized by day
- **Customizable Output**: Use the built-in template or create your own
- **Multiple Sort Orders**: Forward chronological (oldest first) or reverse (newest first)
- **Content Preservation**: Keeps content warnings, media attachments, and post metadata
- **Configuration Flexibility**: Configure via YAML file, environment variables, or CLI flags

## Installation

### From Source

```bash
git clone https://github.com/lmorchard/mastodon-to-markdown.git
cd mastodon-to-markdown
make build
```

The binary will be created as `./mastodon-to-markdown`.

### Prerequisites

- Go 1.23 or later
- A Mastodon account with an access token

## Quick Start

1. **Generate configuration files**:
   ```bash
   ./mastodon-to-markdown init
   ```

2. **Edit the configuration file** (`mastodon-to-markdown.yaml`):
   - Add your Mastodon server URL
   - Add your access token (generate at: Settings > Development > New Application)
   - Required scope: `read:statuses`

3. **Verify your credentials**:
   ```bash
   ./mastodon-to-markdown whoami
   ```

4. **Fetch your posts**:
   ```bash
   ./mastodon-to-markdown fetch --since 7d --output posts.md
   ```

## Configuration

### Generating Access Token

1. Go to your Mastodon instance's settings
2. Navigate to Development > New Application
3. Give it a name (e.g., "mastodon-to-markdown")
4. Required scope: `read:statuses`
5. Copy the access token

### Configuration File

The `init` command creates a `mastodon-to-markdown.yaml` file with sensible defaults:

```yaml
# Logging configuration
verbose: false
debug: false
log_json: false

# Mastodon configuration
mastodon:
  server: "https://mastodon.social"
  access_token: "your-token-here"

# Output configuration
output:
  include_metadata: true
  include_media_urls: true
  template: ""  # Empty = use built-in template
  sort_order: "asc"  # "asc" (oldest first) or "desc" (newest first)
  public_only: true  # Exclude private and direct messages

# Fetch configuration
fetch:
  exclude_replies: false      # Exclude reply posts
  exclude_boosts: false       # Exclude boosted posts
  exclude_favorites: false    # Exclude favorited posts
  visibility: ""              # Filter by visibility
```

### Environment Variables

All configuration options can be set via environment variables using the `MASTODON_TO_MARKDOWN_` prefix:

```bash
export MASTODON_TO_MARKDOWN_MASTODON_SERVER="https://your-instance.social"
export MASTODON_TO_MARKDOWN_MASTODON_ACCESS_TOKEN="your-token"
export MASTODON_TO_MARKDOWN_OUTPUT_SORT_ORDER="desc"
export MASTODON_TO_MARKDOWN_FETCH_EXCLUDE_FAVORITES="true"   # Exclude favorites
export MASTODON_TO_MARKDOWN_FETCH_EXCLUDE_REPLIES="true"     # Exclude replies
```

## Usage

### Commands

#### `init` - Initialize configuration

Create default configuration and template files:

```bash
# Basic init
mastodon-to-markdown init

# Custom template filename
mastodon-to-markdown init --template-file my-template.md

# Force overwrite existing files
mastodon-to-markdown init --force
```

#### `whoami` - Verify credentials

Show information about the authenticated account:

```bash
mastodon-to-markdown whoami
```

Output includes server, username, display name, follower counts, and account URL.

#### `fetch` - Export posts

Fetch and export posts to markdown:

```bash
# Fetch last 7 days
mastodon-to-markdown fetch --since 7d --output posts.md

# Fetch specific date range
mastodon-to-markdown fetch --start 2025-11-01 --end 2025-11-07 --output posts.md

# Fetch last 24 hours, exclude replies
mastodon-to-markdown fetch --since 24h --exclude-replies --output today.md

# Exclude favorited posts (only your own posts and boosts)
mastodon-to-markdown fetch --since 7d --exclude-favorites --output posts.md

# Fetch to stdout (for piping)
mastodon-to-markdown fetch --since 7d

# Reverse chronological order (newest first)
mastodon-to-markdown fetch --since 7d --sort-order desc --output posts.md

# Include private posts
mastodon-to-markdown fetch --since 7d --public-only=false --output all-posts.md
```

#### `version` - Show version

Display version information:

```bash
mastodon-to-markdown version
```

### Fetch Options

| Flag | Description | Default |
|------|-------------|---------|
| `--since` | Time period (e.g., '24h', '7d', '2w') | - |
| `--start` | Start date (YYYY-MM-DD) | - |
| `--end` | End date (YYYY-MM-DD) | - |
| `--output`, `-o` | Output file | stdout |
| `--exclude-replies` | Exclude reply posts | false |
| `--exclude-boosts` | Exclude boosted posts | false |
| `--exclude-favorites` | Exclude favorited posts | false |
| `--public-only` | Only public posts | true |
| `--sort-order` | Sort: 'asc' or 'desc' | asc |
| `--visibility` | Filter by visibility (comma-separated) | - |

### Global Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Config file path | ./mastodon-to-markdown.yaml |
| `--verbose`, `-v` | Verbose output | false |
| `--debug` | Debug output | false |
| `--log-json` | JSON log format | false |

## Template Customization

### Using the Default Template

Leave the `template` option empty or unset to use the built-in template.

### Creating a Custom Template

1. Generate the default template:
   ```bash
   mastodon-to-markdown init
   ```

2. Edit `mastodon-to-markdown.md` to customize the output format

3. The template uses Go's `text/template` syntax with the following data structure:

```go
type TemplateData struct {
    StartDate string    // Formatted start date
    EndDate   string    // Formatted end date
    Posts     []Post    // Array of posts
}

type Post struct {
    ID               string
    CreatedAt        time.Time
    FormattedTime    string
    URL              string
    Content          string        // Cleaned HTML content
    ContentWarning   string
    Visibility       string
    IsReply          bool
    IsBoost          bool
    MediaAttachments []MediaAttachment
}
```

### Example Custom Template

```markdown
# My Posts ({{.StartDate}} - {{.EndDate}})

Total posts: {{len .Posts}}

{{range .Posts}}
### {{.FormattedTime}}

{{if .ContentWarning}}⚠️ **CW**: {{.ContentWarning}}{{end}}

{{.Content}}

🔗 [View on Mastodon]({{.URL}})

{{if .MediaAttachments}}**Attachments**: {{range .MediaAttachments}}[{{.Type}}]({{.URL}}) {{end}}{{end}}

---
{{end}}
```

## Examples

### Blog Post Draft

Fetch your posts from the last week and create a blog post draft:

```bash
mastodon-to-markdown fetch --since 7d \
  --exclude-replies \
  --exclude-boosts \
  --exclude-favorites \
  --output blog-draft-$(date +%Y-%m-%d).md
```

### Weekly Roundup

Create a weekly roundup of all your public activity (posts, boosts, and favorites):

```bash
mastodon-to-markdown fetch --since 7d \
  --sort-order desc \
  --output weekly-roundup.md
```

### Content Curation

Create a curated list of interesting posts you've favorited:

```bash
mastodon-to-markdown fetch --since 30d \
  --exclude-replies \
  --exclude-boosts \
  --output curated-favorites.md
```

### Archive Month

Archive all posts from a specific month:

```bash
mastodon-to-markdown fetch \
  --start 2025-11-01 \
  --end 2025-11-30 \
  --public-only=false \
  --output archive-2025-11.md
```

### Thread Export

Export just your original posts (no replies, boosts, or favorites):

```bash
mastodon-to-markdown fetch --since 30d \
  --exclude-replies \
  --exclude-boosts \
  --exclude-favorites \
  --output my-threads.md
```

## Development

### Building

```bash
make build
```

### Linting

```bash
make setup  # Install development tools
make lint   # Run linters
```

### Formatting

```bash
make format
```

### Cleaning

```bash
make clean
```

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

See LICENSE file for details.
