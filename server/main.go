package main

import (
	"fmt"
	"net/http"

	"github.com/arxdsilva/Stsuru/handlers"
	"github.com/gorilla/mux"
)

// Router registers the routes from Stsuru
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Home)
	r.HandleFunc("/css/", handlers.Css)
	r.HandleFunc("/link/add", handlers.AddLink)
	r.HandleFunc("/redirect/{id}", handlers.LinkSolver)
	r.HandleFunc("/link/remove/{id}", handlers.RemoveLink)
	http.Handle("/", r)
	fmt.Println("The server is now live @ localhost:8080")
	http.ListenAndServe(":8080", nil)
}
