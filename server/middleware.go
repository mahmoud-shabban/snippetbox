package main

import (
	"fmt"
	"net/http"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Code will execute when the request is recieved before it is passed to serveMux or next handler
		// Note: This is split across multiple lines for readability. You don't
		// need to do this in your own code.
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Server", "Go")
		// return // will cause the chain to go back instead of moving to next handler it doesn't go to servmux nor to app handler it goes back to client

		next.ServeHTTP(w, r)

		// Code that will be executed in the way of response after handler is done and response is in its way back
		// demonstration: commonHeaders → servemux → application handler → servemux → commonHeaders
	})
}

func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip        = r.RemoteAddr
			method    = r.Method
			uri       = r.URL.RequestURI()
			proto     = r.Proto
			UserAgent = r.UserAgent()
		)
		app.logger.Info("recieved request", "ip", ip, "method", method, "proto", proto, "uri", uri, "userAgent", UserAgent)
		next.ServeHTTP(w, r)
	})
}

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
