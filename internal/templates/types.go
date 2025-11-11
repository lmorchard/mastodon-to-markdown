package templates

import "time"

// TemplateData represents the complete data structure passed to templates
type TemplateData struct {
	StartDate string
	EndDate   string
	Posts     []Post
}

// Post represents a Mastodon post with all relevant fields for templating
type Post struct {
	ID               string
	CreatedAt        time.Time
	FormattedTime    string
	URL              string
	Content          string
	ContentWarning   string
	Visibility       string
	IsReply          bool
	IsBoost          bool
	MediaAttachments []MediaAttachment
	RepliesCount     int64
	ReblogsCount     int64
	FavouritesCount  int64
}

// MediaAttachment represents a media file attached to a post
type MediaAttachment struct {
	Type        string // "image", "video", "gifv", "audio", "unknown"
	URL         string
	PreviewURL  string
	Description string
}
