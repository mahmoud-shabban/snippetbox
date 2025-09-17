package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
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

func (app *Application) clientError(w http.ResponseWriter, r *http.Request, status int) {
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

func (app *Application) newTemplateData() templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
	}
}

func humanDate(t time.Time) string {
	return t.Format("02 Feb 2007 at 13:05")
}
