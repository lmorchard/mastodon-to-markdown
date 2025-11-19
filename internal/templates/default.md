# Posts from {{.StartDate}} to {{.EndDate}}
{{range .Days}}
## {{.Date}}
{{if .OwnPosts}}
### My Posts
{{range .OwnPosts}}
#### {{.FormattedTimeOnly}}

{{if .ContentWarning}}CW: {{.ContentWarning}}

{{end}}{{.URL}}

{{.Content}}
{{if .MediaAttachments}}

{{range .MediaAttachments}}Media: [{{.Type}}]({{.URL}}){{if .Description}} - {{.Description}}{{end}}
{{end}}{{end}}
---

{{end}}{{end}}
{{if .BoostedPosts}}
### Posts I Boosted
{{range .BoostedPosts}}
#### {{.FormattedTimeOnly}}

{{if .BoostCommentary}}My commentary: {{.BoostCommentary}}

{{end}}{{if .OriginalPost}}**{{.OriginalPost.AuthorName}}** ([@{{.OriginalPost.AuthorUsername}}]({{.OriginalPost.AuthorURL}}))

{{if .OriginalPost.ContentWarning}}CW: {{.OriginalPost.ContentWarning}}

{{end}}{{.OriginalPost.URL}}

{{.OriginalPost.Content}}
{{if .OriginalPost.MediaAttachments}}

{{range .OriginalPost.MediaAttachments}}Media: [{{.Type}}]({{.URL}}){{if .Description}} - {{.Description}}{{end}}
{{end}}{{end}}{{end}}
---

{{end}}{{end}}
{{if .FavoritedPosts}}
### Posts I Favorited
{{range .FavoritedPosts}}
#### {{.FormattedTimeOnly}}

{{if .OriginalPost}}**{{.OriginalPost.AuthorName}}** ([@{{.OriginalPost.AuthorUsername}}]({{.OriginalPost.AuthorURL}}))

{{if .OriginalPost.ContentWarning}}CW: {{.OriginalPost.ContentWarning}}

{{end}}{{.OriginalPost.URL}}

{{.OriginalPost.Content}}
{{if .OriginalPost.MediaAttachments}}

{{range .OriginalPost.MediaAttachments}}Media: [{{.Type}}]({{.URL}}){{if .Description}} - {{.Description}}{{end}}
{{end}}{{end}}{{end}}
---

{{end}}{{end}}
{{end}}
