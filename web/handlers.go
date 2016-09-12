package web

import (
	"log"
	"net/http"

	"github.com/alecthomas/template"
	"github.com/arxdsilva/Stsuru/web/mngo"
	"github.com/gorilla/mux"
)

// AddLink validates the request's URL and asks Mongo to add It on list
func AddLink(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	l := r.Form["user_link"][0]
	err := mngo.Insert(l, w, r)
	checkError(err)
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

// Home querys Mongo for all It's elements and calls the specified HTML to load them into the page.
func Home(w http.ResponseWriter, r *http.Request) {
	Data, err := mngo.GetAll()
	checkError(err)
	t, err := template.ParseFiles("../../tmpl/index.html")
	checkError(err)
	err = t.Execute(w, Data)
	checkError(err)
}

// CSS SHOULD load style into the page :p
func CSS(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/css/", http.FileServer(http.Dir("../../css/")))
}

// RemoveLink searches db for a certain link & removes It if It exists
func RemoveLink(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	idInfo := id["id"]
	mngo.Delete(idInfo)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Redirect takes the hashed URL and checks Mongo If It exists;
func Redirect(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	idInfo := id["id"]
	l, err := mngo.FindHash(idInfo)
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
