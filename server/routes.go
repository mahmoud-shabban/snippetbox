package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *Application) routes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.home) // GET method also works with HEAD only one method is allowed in this form of definition
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// test endpoint
	mux.HandleFunc("/test", app.test)

	// file server
	fileserver := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	// return app.recoverPanic(app.logRequest(commonHeaders(mux)))

	middlewares := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return middlewares.Then(mux)
}
