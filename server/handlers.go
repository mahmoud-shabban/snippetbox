package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "snippetBox")
	w.Header().Set("erver", "GO")
	w.Write([]byte(r.PathValue("path")))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home Page...\n"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	// check valid id
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		log.Printf("Error: invalid snippet id {%s}\n", r.PathValue("id"))
		http.NotFound(w, r)
		return
	}

	w.Write([]byte(fmt.Sprintf("view snippet #%d...\n", id)))
}
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("display new snippept form...\n"))
}
func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("save new snippet to DB...\n"))
}
