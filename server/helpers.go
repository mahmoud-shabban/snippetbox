package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

func (app *Application) check(err error) {
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}

func (app *Application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), slog.Any("method", method), slog.Any("uri", uri))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func (app *Application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {

	ts, ok := app.templateCache[page]

	if !ok {
		err := fmt.Errorf("template %s does not exist", page)
		app.serverError(w, r, err)
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	// err := ts.ExecuteTemplate(w, "base", data)

	_, err = buf.WriteTo(w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *Application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		CSRFToken:       nosurf.Token(r),
		IsAuthenticated: app.isAuthenticated(r),
	}
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// func validateForm(formData *snippetCreateForm) {
// 	if strings.TrimSpace(formData.Title) == "" {
// 		formData.Validations["title"] = "title cannot be blank"
// 	} else if utf8.RuneCountInString(formData.Title) > 100 {
// 		formData.Validations["title"] = "title cannot be more than 100 characters long"
// 	}

// 	if strings.TrimSpace(formData.Content) == "" {
// 		formData.Validations["content"] = "content cannot be blank"
// 	}

// 	if formData.Expires != 1 && formData.Expires != 7 && formData.Expires != 365 {
// 		formData.Validations["expires"] = "expires must be 1, 7 or 365"
// 	}

// }

func (app *Application) decodePostForm(r *http.Request, dst any) error {

	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)

	if err != nil {

		var invalidDecodeError *form.InvalidDecoderError
		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}

		return err
	}

	return nil
}

func (app *Application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}
