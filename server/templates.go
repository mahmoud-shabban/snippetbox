package main

import (
	"html/template"
	"path/filepath"
	"strings"

	"github.com/mahmoud-shabban/snippetbox/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
}

// var funcMap = template.FuncMap{
// 	"humanDate": humanDate,
// }

func newTemplateCache() (map[string]*template.Template, error) {

	var funcMap = template.FuncMap{
		"humanDate": humanDate,
	}
	cache := make(map[string]*template.Template)

	files, err := filepath.Glob("./ui/html/pages/*.tmpl.html")

	if err != nil {
		return nil, err
	}

	partials, err := filepath.Glob("./ui/html/partials/*.tmpl.html")

	if err != nil {
		return nil, err
	}

	for _, f := range files {
		fname := filepath.Base(f)
		name := strings.Split(fname, ".")[0]
		temps := []string{
			"./ui/html/base.tmpl.html",
			f,
		}

		// temps = append(temps, partials...)
		t, err := template.New("").Funcs(funcMap).ParseFiles(
			//"./ui/html/partials/nav.tmpl.html",
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
