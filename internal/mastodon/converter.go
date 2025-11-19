package mastodon

import (
	"html"
	"strings"

	"github.com/lmorchard/mastodon-to-markdown/internal/templates"
	"github.com/lmorchard/mastodon-to-markdown/internal/timerange"
	"github.com/mattn/go-mastodon"
)

// ConvertStatus converts a Mastodon status to our template Post format
func ConvertStatus(status *mastodon.Status) templates.Post {
	post := templates.Post{
		ID:                string(status.ID),
		CreatedAt:         status.CreatedAt,
		FormattedTime:     timerange.FormatDateTime(status.CreatedAt),
		FormattedDate:     timerange.FormatDate(status.CreatedAt),
		FormattedTimeOnly: status.CreatedAt.Format("15:04"),
		URL:               status.URL,
		Content:           cleanContent(status.Content),
		ContentWarning:    status.SpoilerText,
		Visibility:        string(status.Visibility),
		IsReply:           status.InReplyToID != nil,
		IsBoost:           status.Reblog != nil,
		IsFavorited:       false, // Will be set by ConvertFavourite
		RepliesCount:      status.RepliesCount,
		ReblogsCount:      status.ReblogsCount,
		FavouritesCount:   status.FavouritesCount,
	}

	// If this is a boost, extract the original post and any commentary
	if status.Reblog != nil {
		post.BoostCommentary = cleanContent(status.Content)
		post.OriginalPost = extractOriginalPost(status.Reblog)
	} else {
		// Convert media attachments for non-boost posts
		for _, attachment := range status.MediaAttachments {
			post.MediaAttachments = append(post.MediaAttachments, templates.MediaAttachment{
				Type:        string(attachment.Type),
				URL:         attachment.URL,
				PreviewURL:  attachment.PreviewURL,
				Description: attachment.Description,
			})
		}
	}

	return post
}

// ConvertFavourite converts a favorited Mastodon status to our template Post format
func ConvertFavourite(status *mastodon.Status) templates.Post {
	post := templates.Post{
		ID:                string(status.ID),
		CreatedAt:         status.CreatedAt,
		FormattedTime:     timerange.FormatDateTime(status.CreatedAt),
		FormattedDate:     timerange.FormatDate(status.CreatedAt),
		FormattedTimeOnly: status.CreatedAt.Format("15:04"),
		URL:               status.URL,
		IsFavorited:       true,
		OriginalPost:      extractOriginalPost(status),
	}

	return post
}

// extractOriginalPost extracts original post details from a status
func extractOriginalPost(status *mastodon.Status) *templates.OriginalPost {
	original := &templates.OriginalPost{
		AuthorName:     status.Account.DisplayName,
		AuthorUsername: string(status.Account.Username),
		AuthorURL:      status.Account.URL,
		Content:        cleanContent(status.Content),
		ContentWarning: status.SpoilerText,
		URL:            status.URL,
	}

	// Convert media attachments
	for _, attachment := range status.MediaAttachments {
		original.MediaAttachments = append(original.MediaAttachments, templates.MediaAttachment{
			Type:        string(attachment.Type),
			URL:         attachment.URL,
			PreviewURL:  attachment.PreviewURL,
			Description: attachment.Description,
		})
	}

	return original
}

// ConvertStatuses converts multiple Mastodon statuses
func ConvertStatuses(statuses []*mastodon.Status) []templates.Post {
	posts := make([]templates.Post, 0, len(statuses))
	for _, status := range statuses {
		posts = append(posts, ConvertStatus(status))
	}
	return posts
}

// ConvertFavourites converts multiple favorited Mastodon statuses
func ConvertFavourites(statuses []*mastodon.Status) []templates.Post {
	posts := make([]templates.Post, 0, len(statuses))
	for _, status := range statuses {
		posts = append(posts, ConvertFavourite(status))
	}
	return posts
}

// cleanContent removes HTML tags and unescapes HTML entities
func cleanContent(content string) string {
	// Simple HTML tag removal - for more robust handling, consider using a proper HTML parser
	re := strings.NewReplacer(
		"<p>", "\n\n",
		"</p>", "",
		"<br>", "\n",
		"<br/>", "\n",
		"<br />", "\n",
	)
	content = re.Replace(content)

	// Remove remaining HTML tags
	for strings.Contains(content, "<") && strings.Contains(content, ">") {
		start := strings.Index(content, "<")
		end := strings.Index(content, ">")
		if start < end {
			content = content[:start] + content[end+1:]
		} else {
			break
		}
	}

	// Unescape HTML entities
	content = html.UnescapeString(content)

	// Clean up multiple newlines
	content = strings.TrimSpace(content)

	return content
}
