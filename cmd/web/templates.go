package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"snippetbox.saran.net/internal/models"
	"snippetbox.saran.net/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

func replaceNewLine(str string) template.HTML {
	return template.HTML(strings.ReplaceAll(str, "\\n", "<br>"))
}

var functions = template.FuncMap{
	"humanDate":      humanDate,
	"replaceNewLine": replaceNewLine,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(ui.Files, "html/pages/*.templ")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{
			"html/base.templ",
			"html/partials/*.templ",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
