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
		ID:              string(status.ID),
		CreatedAt:       status.CreatedAt,
		FormattedTime:   timerange.FormatDateTime(status.CreatedAt),
		URL:             status.URL,
		Content:         cleanContent(status.Content),
		ContentWarning:  status.SpoilerText,
		Visibility:      string(status.Visibility),
		IsReply:         status.InReplyToID != nil,
		IsBoost:         status.Reblog != nil,
		RepliesCount:    status.RepliesCount,
		ReblogsCount:    status.ReblogsCount,
		FavouritesCount: status.FavouritesCount,
	}

	// Convert media attachments
	for _, attachment := range status.MediaAttachments {
		post.MediaAttachments = append(post.MediaAttachments, templates.MediaAttachment{
			Type:        string(attachment.Type),
			URL:         attachment.URL,
			PreviewURL:  attachment.PreviewURL,
			Description: attachment.Description,
		})
	}

	return post
}

// ConvertStatuses converts multiple Mastodon statuses
func ConvertStatuses(statuses []*mastodon.Status) []templates.Post {
	posts := make([]templates.Post, 0, len(statuses))
	for _, status := range statuses {
		posts = append(posts, ConvertStatus(status))
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
