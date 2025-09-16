package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
)

func (app *Application) test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "snippetBox")
	w.Header().Set("erver", "GO")
	w.Write([]byte(r.PathValue("path")))
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("server", "GO")

	templates := []string{
		"./ui/html/pages/home.tmpl",
		"./ui/html/pages/base.tmpl",
		"./ui/html/partials/nav.tmpl",
	}
	tmpl, err := template.ParseFiles(templates...) // path relative to root dir snippetbox
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
	// check valid id
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.logger.Error("invalid snippet id", slog.Any("id", r.PathValue("id")))
		http.NotFound(w, r)
		return
	}

	w.Write([]byte(fmt.Sprintf("view snippet #%d...\n", id)))
}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("display new snippept form...\n"))
}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("save new snippet to DB...\n"))
}
