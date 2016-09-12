package main

import (
	"fmt"
	"net/http"

	"github.com/arxdsilva/Stsuru/web"
	"github.com/gorilla/mux"
)

// Registers the routes used by Stsuru
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", web.Home)
	// r.HandleFunc("/css/", web.CSS)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../../css"))))
	r.HandleFunc("/r/{id}", web.Redirect)
	r.HandleFunc("/link/add", web.AddLink)
	r.HandleFunc("/link/remove/{id}", web.RemoveLink)
	http.Handle("/", r)
	fmt.Println("The server is now live @ localhost:8080")
	http.ListenAndServe(":8080", nil)
}
