package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mahmoud-shabban/snippetbox/internal/models"
)

func (app *Application) test(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Server", "snippetBox")
	// w.Header().Set("erver", "GO")
	panic("Server is panicing, don't worry!!")
	w.Write([]byte(r.PathValue("path")))
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {

	// w.Header().Add("server", "GO")

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData()
	data.Snippets = snippets
	// data := templateData{Snippets: snippets}

	app.render(w, r, http.StatusOK, "home", data)

	// templates := []string{
	// 	"./ui/html/base.tmpl.html",
	// 	"./ui/html/pages/home.tmpl.html",
	// 	"./ui/html/partials/nav.tmpl.html",
	// }
	// tmpl, err := template.ParseFiles(templates...) // path relative to root dir snippetbox
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// err = tmpl.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// }

}

func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
	// check valid id
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.logger.Error("invalid snippet id", slog.Any("id", r.PathValue("id")))
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData()
	data.Snippet = snippet
	// data := templateData{Snippet: snippet}

	app.render(w, r, http.StatusOK, "view", data)

	// templates := []string{
	// 	"./ui/html/base.tmpl.html",
	// 	"./ui/html/partials/nav.tmpl.html",
	// 	"./ui/html/pages/view.tmpl.html",
	// }

	// tmpl, err := template.ParseFiles(templates...)

	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// err = tmpl.ExecuteTemplate(w, "base", data)

	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }
	// fmt.Fprintf(w, "%+v", snippet)
	// w.Write([]byte(fmt.Sprintf("view snippet #%d...\n", id)))
}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("display new snippept form...\n"))

}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(http.StatusCreated)
	// w.Write([]byte("save new snippet to DB...\n"))
	title := "new snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 10

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
