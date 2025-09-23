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

type snippetCreateForm struct {
	validator.Validator `form:"-"`
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	// Validations map[string]string
}

type userSignupForm struct {
	validator.Validator `form:"-"`
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	// Validations map[string]string
}

type userLoginForm struct {
	validator.Validator `form:"-"`
	Email               string `form:"email"`
	Password            string `form:"password"`
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

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup", data)

}
func (app *Application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "Name cannot be empty")
	form.CheckField(validator.NotBlank(form.Email), "email", "email cannot be empty")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "email not valid")
	form.CheckField(validator.NotBlank(form.Password), "password", "password cannot be empty")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "password must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "signup", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)

	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			data := app.newTemplateData(r)
			form.AddFieldError("email", "Email address already in use")

			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup", data)

		} else {
			app.serverError(w, r, err)
			return
		}
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please login.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login", data)
}

func (app *Application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "Email cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Please enter valid email address")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Invalid Email or Password")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login", data)
			return
		} else {
			app.serverError(w, r, err)
			return
		}
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *Application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	err := app.sessionManager.RenewToken(r.Context())

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "You have been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)

}
