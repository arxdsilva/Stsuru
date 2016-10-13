package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/arxdsilva/Stsuru/shortener"
	"github.com/arxdsilva/Stsuru/web/persist"
	"github.com/arxdsilva/Stsuru/web/persist/data"

	"github.com/alecthomas/template"
	"github.com/gorilla/mux"
)

// Server is the user's way to customize which storage & Host will be used
type Server struct {
	Storage    persist.Storage
	CustomHost string
}

// Listen Registers the routes used by Stsuru and redirects traffic
func (s *Server) Listen() {
	r := mux.NewRouter()
	r.HandleFunc("/", s.home)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	r.HandleFunc("/r/{id}", s.redirectLink)
	r.HandleFunc("/link/add", s.addLink)
	r.HandleFunc("/l/r/{id}", s.removeLink)
	http.Handle("/", r)
	fmt.Println("The server is now live @ localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func (s *Server) addLink(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	link := r.Form["user_link"][0]
	u, err := url.Parse(link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotModified)
		return
	}
	newShort := shortener.NewShorten{
		U: u,
	}
	n, err := newShort.Shorten()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotModified)
		return
	}
	linkshort := n.String()
	dbHash := n.Path
	_, err = s.Storage.FindHash(dbHash)
	if err != nil {
		Data := data.LinkData{
			Link:       link,
			Hash:       dbHash,
			Short:      linkshort,
			CustomHost: s.CustomHost,
		}
		err = s.Storage.Save(&Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotModified)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusNotModified)
	return
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	path := "tmpl/index.html"
	d, err := s.Storage.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) removeLink(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	idHash := id["id"]
	s.Storage.Remove(idHash)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) redirectLink(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	idHash := id["id"]
	l, err := s.Storage.FindHash(idHash)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, l, http.StatusFound)
}
