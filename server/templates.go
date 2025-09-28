package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/mahmoud-shabban/snippetbox/internal/models"
	"github.com/mahmoud-shabban/snippetbox/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func newTemplateCache() (map[string]*template.Template, error) {

	var funcMap = template.FuncMap{
		"humanDate": humanDate,
	}
	cache := make(map[string]*template.Template)

	files, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	partials, err := fs.Glob(ui.Files, "html/partials/*.tmpl.html")

	if err != nil {
		return nil, err
	}

	for _, f := range files {
		fname := filepath.Base(f)
		name := strings.Split(fname, ".")[0]
		temps := []string{
			"html/base.tmpl.html",
			f,
		}

		t, err := template.New("").Funcs(funcMap).ParseFS(
			ui.Files,
			append(temps, partials...)...,
		)

		t = t.Funcs(funcMap)

		if err != nil {
			return nil, err
		}

		cache[name] = t
	}

	return cache, nil
}
