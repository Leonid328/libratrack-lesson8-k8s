package web

import (
	"embed"
	"html/template"
	"io/fs"
)

//go:embed templates/*.html templates/partials/*.html
var templateFS embed.FS

func ParseTemplates() (*template.Template, error) {
	return template.ParseFS(templateFS, "templates/*.html", "templates/partials/*.html")
}

func TemplateFS() fs.FS {
	return templateFS
}
