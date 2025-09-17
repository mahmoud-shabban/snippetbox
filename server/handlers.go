package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/mahmoud-shabban/snippetbox/internal/models"
)

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	Validations map[string]string
}

func (app *Application) test(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Server", "snippetBox")
	// w.Header().Set("erver", "GO")
	// panic("Server is panicing, don't worry!!")
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
	// w.Write([]byte("display new snippept form...\n"))
	data := app.newTemplateData()
	form := snippetCreateForm{
		Expires: 1,
	}

	data.Form = form
	app.render(w, r, http.StatusOK, "create", data)
}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(http.StatusCreated)
	// w.Write([]byte("save new snippet to DB...\n"))
	// title := "new snail"
	// content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	// expires := 10

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// extracting form data and making validations
	formData := snippetCreateForm{}
	// validations := make(map[string]string)
	formData.Title = r.PostForm.Get("title")
	formData.Content = r.PostForm.Get("content")
	formData.Expires, err = strconv.Atoi(r.PostForm.Get("expires"))
	formData.Validations = make(map[string]string)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(formData.Title) == "" {
		formData.Validations["title"] = "title cannot be blank"
	} else if utf8.RuneCountInString(formData.Title) > 100 {
		formData.Validations["title"] = "title cannot be more than 100 characters long"
	}

	if strings.TrimSpace(formData.Content) == "" {
		formData.Validations["content"] = "content cannot be blank"
	}

	if formData.Expires != 1 && formData.Expires != 7 && formData.Expires != 365 {
		formData.Validations["expires"] = "expires must be 1, 7 or 365"
	}

	if len(formData.Validations) > 0 {
		data := app.newTemplateData()
		data.Form = formData
		app.render(w, r, http.StatusUnprocessableEntity, "create", data)
		return
	}

	id, err := app.snippets.Insert(formData.Title, formData.Content, formData.Expires)
	if err != nil {
		app.serverError(w, r, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
