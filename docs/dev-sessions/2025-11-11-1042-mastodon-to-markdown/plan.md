# Implementation Plan - Mastodon to Markdown

## Overview

Build a Go CLI tool to fetch Mastodon posts and export them to markdown, following the established patterns from the go-cli-builder skill.

## Phase 1: Project Scaffold (COMPLETED)

- [x] Create dev session directory structure
- [x] Write spec.md
- [x] Scaffold project using go-cli-builder
- [x] Update go.mod with proper module name and dependencies
- [x] Add fetch command skeleton
- [x] Update configuration structures
- [x] Update example YAML file

## Phase 2: Core Infrastructure (CURRENT)

### 2.1 Template System
- [ ] Create `internal/templates/` package
- [ ] Define template data structures (Post, MediaAttachment, etc.)
- [ ] Create default markdown template using text/template syntax
- [ ] Embed default templates using `//go:embed`
- [ ] Implement template loading (default vs custom file)
- [ ] Implement template rendering function

### 2.2 Mastodon Client
- [ ] Create `internal/mastodon/` package
- [ ] Implement client initialization from config
- [ ] Implement authentication verification
- [ ] Add error handling for common API errors

## Phase 3: Post Fetching

### 3.1 Time Range Parsing
- [ ] Create `internal/timerange/` package
- [ ] Parse "since" duration strings (e.g., "24h", "7d")
- [ ] Parse start/end date strings (YYYY-MM-DD)
- [ ] Validate time range logic

### 3.2 Post Retrieval
- [ ] Implement account lookup (verify authenticated user)
- [ ] Implement statuses fetching with pagination
- [ ] Apply time range filtering
- [ ] Apply visibility filtering
- [ ] Handle exclude-replies flag
- [ ] Handle exclude-boosts flag

### 3.3 Data Structures
- [ ] Define internal Post structure
- [ ] Convert Mastodon API status to internal Post
- [ ] Extract relevant metadata
- [ ] Handle content warnings
- [ ] Handle media attachments

## Phase 4: Markdown Export

### 4.1 Markdown Generation
- [ ] Create `internal/export/` package
- [ ] Implement markdown formatter
- [ ] Apply templates to posts
- [ ] Handle special markdown characters in content
- [ ] Format timestamps
- [ ] Format URLs

### 4.2 Output Handling
- [ ] Implement file output
- [ ] Implement stdout output
- [ ] Add output validation

## Phase 5: Integration

### 5.1 Wire Up fetch Command
- [ ] Initialize Mastodon client
- [ ] Parse and validate flags
- [ ] Fetch posts
- [ ] Export to markdown
- [ ] Add progress logging
- [ ] Add error handling

### 5.2 Database Integration (Optional for v1)
- [ ] Define schema for cached posts
- [ ] Implement post caching
- [ ] Check cache before API calls

## Phase 6: Testing & Polish

### 6.1 Testing
- [ ] Add unit tests for time range parsing
- [ ] Add unit tests for template rendering
- [ ] Add unit tests for post filtering
- [ ] Add integration test (may need mocked API)

### 6.2 Documentation
- [ ] Update README with usage examples
- [ ] Document template format
- [ ] Document configuration options
- [ ] Add troubleshooting section

### 6.3 Build & Release
- [ ] Test make targets (build, lint, format)
- [ ] Verify GitHub Actions workflows
- [ ] Test cross-compilation
- [ ] Create initial release

## Implementation Notes

### Key Dependencies
- `github.com/mattn/go-mastodon` - Mastodon API client
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration
- `text/template` - Template rendering
- `embed` - Embed default templates

### Template Design

Default template should produce output like:

```markdown
# Posts from {{.StartDate}} to {{.EndDate}}

{{range .Posts}}
## {{.FormattedTime}}

{{if .ContentWarning}}CW: {{.ContentWarning}}{{end}}

{{.URL}}

{{.Content}}

{{range .MediaAttachments}}
Media: [{{.Type}}]({{.URL}})
{{end}}

---

{{end}}
```

### Directory Structure

```
internal/
├── config/          # Configuration struct (done)
├── database/        # Database layer (done)
├── mastodon/        # Mastodon API client wrapper
├── timerange/       # Time range parsing utilities
├── templates/       # Template management
│   └── default.md   # Default template (embedded)
└── export/          # Markdown export logic
```

### Open Questions

1. Should we cache posts in the database to avoid repeated API calls?
   - Pros: Faster, less API load
   - Cons: More complexity, stale data
   - Decision: Skip for v1, add later if needed

2. How to handle threads? Export as separate posts or group them?
   - Decision: Export as separate posts for v1, add thread detection later

3. Should we support multiple output formats (JSON, HTML)?
   - Decision: Markdown only for v1

## Success Criteria

- [ ] Can authenticate with Mastodon instance
- [ ] Can fetch posts from last 7 days
- [ ] Can export posts to markdown file
- [ ] Can customize output with custom template
- [ ] Configuration persists between runs
- [ ] Clear error messages
- [ ] Passes linting and formatting checks
- [ ] Has basic test coverage

## Next Steps

1. Implement template system first (foundational)
2. Implement Mastodon client wrapper
3. Implement time range parsing
4. Wire everything together in fetch command
5. Test and polish
