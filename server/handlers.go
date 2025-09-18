package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mahmoud-shabban/snippetbox/internal/models"
	"github.com/mahmoud-shabban/snippetbox/internal/validator"
)

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

	data := app.newTemplateData(r)
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

	// flash := app.sessionManager.PopString(r.Context(), "flash")

	data := app.newTemplateData(r)
	data.Snippet = snippet
	// data.Flash = flash
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
	data := app.newTemplateData(r)
	form := snippetCreateForm{
		Expires: 1,
	}

	data.Form = form
	app.render(w, r, http.StatusOK, "create", data)
}

type snippetCreateForm struct {
	validator.Validator `form:"-"`
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	// Validations map[string]string
}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(http.StatusCreated)
	// w.Write([]byte("save new snippet to DB...\n"))
	// title := "new snail"
	// content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	// expires := 10

	// err := r.ParseForm()

	// extracting form data and making validations
	// form := snippetCreateForm{
	// 	Validator: validator.Validator{Errors: make(map[string]string)},
	// }

	// form := snippetCreateForm{}
	var form snippetCreateForm

	err := app.decodePostForm(r, &form)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// validations := make(map[string]string)
	// form.Validations = make(map[string]string)
	// form.Title = r.PostForm.Get("title")
	// form.Content = r.PostForm.Get("content")

	// form.Expires, err = strconv.Atoi(r.PostForm.Get("expires"))

	err = app.formDecoder.Decode(&form, r.PostForm)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "title cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "title cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "content cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "allowed values 1, 7 or 365")
	// if strings.TrimSpace(form.Title) == "" {
	// 	form.Validations["title"] = "title cannot be blank"
	// } else if utf8.RuneCountInString(form.Title) > 100 {
	// 	form.Validations["title"] = "title cannot be more than 100 characters long"
	// }

	// if strings.TrimSpace(form.Content) == "" {
	// 	form.Validations["content"] = "content cannot be blank"
	// }

	// if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
	// 	form.Validations["expires"] = "expires must be 1, 7 or 365"
	// }

	// validateForm(&form)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet Successfully Created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
