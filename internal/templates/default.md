# Posts from {{.StartDate}} to {{.EndDate}}

{{range .Posts}}
## {{.FormattedTime}}

{{if .ContentWarning}}CW: {{.ContentWarning}}

{{end}}{{.URL}}

{{.Content}}
{{if .MediaAttachments}}

{{range .MediaAttachments}}Media: [{{.Type}}]({{.URL}}){{if .Description}} - {{.Description}}{{end}}
{{end}}{{end}}
---

{{end}}
