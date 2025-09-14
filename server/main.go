package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	wd, _ := os.Getwd()
	log.Printf("working dir: %s", wd)
	log.Println("Snipppetbox server started on :8080")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home) // GET method also works with HEAD only one method is allowed in this form of definition
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	// test endpoint
	mux.HandleFunc("/test", test)

	// file server
	fileserver := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	check(http.ListenAndServe(":8080", mux))
}
