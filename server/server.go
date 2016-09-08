package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Router registers the routes from Stsuru
func Server() {
	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.HandleFunc("/css/", css)
	r.HandleFunc("/link/add", AddLink)
	r.HandleFunc("/redirect/{id}", LinkSolver)
	r.HandleFunc("/link/remove/{id}", RemoveLink)
	http.Handle("/", r)
	fmt.Println("The server is now live @ localhost:8080")
	http.ListenAndServe(":8080", nil)
}
