package config

// Config holds application configuration
type Config struct {
	// Core settings
	Database string
	Verbose  bool
	Debug    bool
	LogJSON  bool

	// Mastodon settings
	Mastodon struct {
		Server      string
		AccessToken string
	}

	// Output settings
	Output struct {
		IncludeMetadata  bool
		IncludeMediaURLs bool
		Template         string // Template to use: "default" (built-in) or path to custom file
		SortOrder        string // "asc" (oldest first) or "desc" (newest first)
		PublicOnly       bool   // Only include public posts (exclude direct/private)
	}
}
