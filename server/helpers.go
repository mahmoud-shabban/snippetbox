package main

import (
	"log/slog"
	"net/http"
	"os"
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
