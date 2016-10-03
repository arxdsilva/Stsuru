package server

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/arxdsilva/Stsuru/web/persist"
	"github.com/asaskevich/govalidator"

	"github.com/alecthomas/template"
	"github.com/gorilla/mux"
)

// Server ...
type Server struct {
	Storage persist.Storage
	URL     string
}

// Listen Registers the routes used by Stsuru and redirects traffic
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
	link := r.Form["user_link"][0]
	v := validateURL(link)
	if !v {
		http.Redirect(w, r, "/", http.StatusNotModified)
		return
	}
	linkshort, dbHash := hash(link, s.URL)
	_, err := s.Storage.FindHash(dbHash)
	if err != nil {
		err = s.Storage.Save(link, linkshort, dbHash)
		checkError(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusNotModified)
}

// Home querys Storage for all It's elements and calls the specified HTML to load them into the page.
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
	idHash := id["id"]
	s.Storage.Remove(idHash)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Redirect takes the hashed URL and checks Mongo If It exists;
func (s *Server) Redirect(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	idHash := id["id"]
	l, err := s.Storage.FindHash(idHash)
	if err != nil {
		http.Redirect(w, r, l, http.StatusNotFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)

}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
	return
}

func hash(link, path string) (string, string) {
	h := md5.New()
	io.WriteString(h, link)
	hash := string(h.Sum(nil))
	linkShort := fmt.Sprintf("%s%x", path, hash)
	dbHash := fmt.Sprintf("%x", hash)
	return linkShort, dbHash
}

func validateURL(l string) bool {
	isURL := govalidator.IsURL(l)
	validURL := govalidator.IsRequestURL(l)
	if isURL == false || validURL == false {
		return false
	}
	return true
}
