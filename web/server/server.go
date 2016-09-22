package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arxdsilva/Stsuru/web/persist"

	"github.com/alecthomas/template"
	"github.com/gorilla/mux"
)

// Server is a struct that implements ...
type Server struct {
	Storage persist.Storage
}

// Listen Registers the routes used by Stsuru
func (s *Server) Listen() {
	r := mux.NewRouter()
	r.HandleFunc("/", s.Home)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	r.HandleFunc("/r/{id}", s.Redirect)
	r.HandleFunc("/link/add", s.AddLink)
	r.HandleFunc("/l/r/{id}", s.RemoveLink)
	http.Handle("/", r)
	fmt.Println("The server is now live @ localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// AddLink validates the request's URL and asks Mongo to add It on list
func (s *Server) AddLink(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	l := r.Form["user_link"][0]
	err := s.Storage.Save(l)
	checkError(err)
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

// Home querys Mongo for all It's elements and calls the specified HTML to load them into the page.
func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	path := "tmpl/index.html"
	d, err := s.Storage.List()
	checkError(err)
	t, err := template.ParseFiles(path)
	checkError(err)
	err = t.Execute(w, d)
	checkError(err)
}

// CSS loads style into the page
func CSS(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/css/", http.FileServer(http.Dir("css/")))
}

// RemoveLink searches db for a certain link & removes It if It exists
func (s *Server) RemoveLink(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	idInfo := id["id"]
	s.Storage.Remove(idInfo)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Redirect takes the hashed URL and checks Mongo If It exists;
func (s *Server) Redirect(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	idInfo := id["id"]
	l, err := s.Storage.FindHash(idInfo)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
	}
	http.Redirect(w, r, l, http.StatusFound)
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
	return
}
