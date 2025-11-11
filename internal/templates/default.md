# Posts from {{.StartDate}} to {{.EndDate}}
{{$currentDate := ""}}
{{range .Posts}}
{{if ne .FormattedDate $currentDate}}
{{$currentDate = .FormattedDate}}

## {{.FormattedDate}}
{{end}}

### {{.FormattedTimeOnly}}

{{if .ContentWarning}}CW: {{.ContentWarning}}

{{end}}{{.URL}}

{{.Content}}
{{if .MediaAttachments}}

{{range .MediaAttachments}}Media: [{{.Type}}]({{.URL}}){{if .Description}} - {{.Description}}{{end}}
{{end}}{{end}}
---

{{end}}
