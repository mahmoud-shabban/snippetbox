package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Snipppetbox server started on :8080")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home) // GET method also works with HEAD only one method is allowed in this form of definition
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	// test endpoint
	mux.HandleFunc("/test", test)

	check(http.ListenAndServe(":8080", mux))
}
