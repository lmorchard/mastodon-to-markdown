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
		Template         string
	}
}
