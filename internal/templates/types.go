package templates

import "time"

// TemplateData represents the complete data structure passed to templates
type TemplateData struct {
	StartDate string
	EndDate   string
	Posts     []Post
	Days      []DayGroup // Posts grouped by day
}

// DayGroup represents all posts for a specific day, organized by type
type DayGroup struct {
	Date            string
	OwnPosts        []Post
	BoostedPosts    []Post
	FavoritedPosts  []Post
}

// Post represents a Mastodon post with all relevant fields for templating
type Post struct {
	ID                string
	CreatedAt         time.Time
	FormattedTime     string // Full date and time (e.g., "2025-11-11 14:30")
	FormattedDate     string // Date only (e.g., "2025-11-11")
	FormattedTimeOnly string // Time only (e.g., "14:30")
	URL               string
	Content           string
	ContentWarning    string
	Visibility        string
	IsReply           bool
	IsBoost           bool
	IsFavorited       bool   // This post was favorited by the user (from favourites endpoint)
	MediaAttachments  []MediaAttachment
	RepliesCount      int64
	ReblogsCount      int64
	FavouritesCount   int64

	// For boosted posts
	BoostCommentary   string // User's commentary when boosting
	OriginalPost      *OriginalPost // Details of the original boosted/favorited post
}

// OriginalPost represents the original post that was boosted or favorited
type OriginalPost struct {
	AuthorName        string
	AuthorUsername    string
	AuthorURL         string
	Content           string
	ContentWarning    string
	URL               string
	MediaAttachments  []MediaAttachment
}

// MediaAttachment represents a media file attached to a post
type MediaAttachment struct {
	Type        string // "image", "video", "gifv", "audio", "unknown"
	URL         string
	PreviewURL  string
	Description string
}
