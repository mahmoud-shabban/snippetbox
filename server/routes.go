package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/mahmoud-shabban/snippetbox/ui"
)

func (app *Application) routes() http.Handler {

	mux := http.NewServeMux()

	// dynamic middleware is midllewares that need to work on specific handlers only
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home)) // GET method also works with HEAD only one method is allowed in this form of definition
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))

	// user routes
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// csrf protected endpoints
	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))
	// test endpoint
	mux.HandleFunc("/test", app.test)

	// file server
	fileserver := http.FileServerFS(ui.Files)
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	// return app.recoverPanic(app.logRequest(commonHeaders(mux)))

	// standard is the middleware chain that works on all routes
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
