# Mastodon to Markdown - Project Specification

## Overview

A Go CLI tool that queries posts from a Mastodon account via the API for a specified time period and composes their content into a markdown document suitable as a starting point for a blog post.

## Goals

1. Authenticate with Mastodon API (OAuth flow or direct access token)
2. Query posts from a specific account for a given time period
3. Export posts to a well-formatted markdown document
4. Support multiple Mastodon instances
5. Provide filtering options (date ranges, visibility, content warnings, etc.)

## Core Features

### Authentication
- Support direct access token authentication
- Support OAuth flow for better UX (optional, can be added later)
- Store credentials securely in SQLite database
- Support multiple Mastodon instances

### Post Querying
- Fetch posts from authenticated account
- Filter by date range (e.g., "last day", "last week", custom date range)
- Filter by post visibility (public, unlisted, private)
- Include/exclude replies
- Include/exclude boosts

### Markdown Export
- Convert post content to clean markdown
- Include post metadata (timestamp, URL, visibility)
- Handle media attachments (reference with URLs)
- Preserve content warnings
- Format threads appropriately
- Support custom templates using text/template
- Default templates compiled into binary using embed
- Allow custom template files via configuration

## Technical Stack

- **Language**: Go
- **CLI Framework**: Cobra
- **Configuration**: Viper (YAML + CLI flags)
- **Database**: SQLite (for auth tokens, post cache)
- **Mastodon API**: go-mastodon library
- **Logging**: Logrus
- **Templates**: text/template with compiled-in default templates

## Commands

### `init`
Initialize configuration file with prompts for:
- Mastodon instance URL
- Access token (or OAuth credentials)

### `auth` (future)
OAuth authentication flow:
- `auth link` - Generate authorization URL
- `auth code <code>` - Exchange code for token

### `fetch`
Query and export posts to markdown:
- `--since <duration>` - Time period (e.g., "24h", "7d")
- `--start <date>` - Start date (YYYY-MM-DD)
- `--end <date>` - End date (YYYY-MM-DD)
- `--output <file>` - Output file (default: stdout)
- `--exclude-replies` - Skip replies
- `--exclude-boosts` - Skip boosts
- `--visibility <types>` - Filter by visibility (comma-separated)

### `version`
Display version information

## Configuration

YAML configuration file (`mastodon-to-markdown.yaml`):

```yaml
mastodon:
  server: "https://mastodon.social"
  access_token: "..." # Optional if using OAuth

output:
  template: "default" # Future: custom templates
  include_metadata: true
  include_media_urls: true
```

## Output Format

Example markdown output:

```markdown
# Posts from 2025-11-04 to 2025-11-11

## 2025-11-11 10:30

https://mastodon.social/@username/123456

This is the post content with **formatting** preserved.

Media: [image](https://files.mastodon.social/...)

---

## 2025-11-10 15:45

CW: Spoiler

https://mastodon.social/@username/123457

Post content here...

---
```

## Non-Goals (v1)

- Posting to Mastodon (that's feed-to-mastodon's job)
- Fetching other users' posts (focus on own account)
- Advanced filtering (hashtags, mentions, etc.)
- GUI or web interface

## Success Criteria

1. Can authenticate with a Mastodon instance
2. Can fetch posts from the last 7 days
3. Can export posts to a readable markdown file
4. Configuration persists between runs
5. Clear error messages for common issues
