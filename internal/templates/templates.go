package templates

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"text/template"
)

//go:embed default.md
var defaultTemplate string

// GetDefaultTemplate returns the embedded default template content
func GetDefaultTemplate() (string, error) {
	return defaultTemplate, nil
}

// Renderer handles loading and rendering markdown templates
type Renderer struct {
	tmpl *template.Template
}

// NewRenderer creates a new template renderer
// If templatePath is empty, uses the embedded default template
// Otherwise loads the template from the specified file
func NewRenderer(templatePath string) (*Renderer, error) {
	var tmpl *template.Template
	var err error

	if templatePath == "" {
		// Use embedded default template
		tmpl, err = template.New("default").Parse(defaultTemplate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse default template: %w", err)
		}
	} else {
		// Load template from file
		tmpl, err = template.ParseFiles(templatePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load template from %s: %w", templatePath, err)
		}
	}

	return &Renderer{tmpl: tmpl}, nil
}

// Render executes the template with the given data and writes to the writer
func (r *Renderer) Render(w io.Writer, data *TemplateData) error {
	if err := r.tmpl.Execute(w, data); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}
	return nil
}

// RenderToFile executes the template and writes to a file
// If filename is empty or "-", writes to stdout
func (r *Renderer) RenderToFile(filename string, data *TemplateData) error {
	if filename == "" || filename == "-" {
		return r.Render(os.Stdout, data)
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", filename, err)
	}
	defer f.Close()

	if err := r.Render(f, data); err != nil {
		return err
	}

	return nil
}
