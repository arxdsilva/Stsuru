package handlers

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/alecthomas/template"
	"github.com/arxdsilva/Stsuru/mngo"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

// Lines is the structure that is added inside Db
type Lines struct {
	Link  string
	Short string
	Hash  string
}

// AddLink ...
func AddLink(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	link := r.Form["user_link"][0]
	// checking the URL
  validateURL(link)
	path := "http://localhost:8080/"
	// URL hashing
	linkShort, dbHash := Hash(link, path)
	l := &Lines{Link: link, Short: linkShort, Hash: dbHash}
	_, err := mgo.FindOne(dbHash)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	err = mgo.Insert(l)
	CheckError(err)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Home ...
func Home(w http.ResponseWriter, r *http.Request) {
	Data := []Lines{}

	t, err := template.ParseFiles("tmpl/index.html")
	CheckError(err)

	t.Execute(w, Data)
}

// CSS ...
func CSS(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/out/", http.FileServer(http.Dir("out/")))
}

// RemoveLink searches db for a certain link & removes It if It exists
func RemoveLink(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	idInfo := id["id"]
	mngo.Delete(idInfo)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Redirect ...
func Redirect(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	idInfo := id["id"]
	l, err := findOne(idInfo)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	http.Redirect(w, r, l, http.StatusFound)
}

// CheckError ...
func CheckError(err error) {
	if err != nil {
		log.Panic(err)
	}
	return
}

// Hash ....
func Hash(link, path string) (string, string) {
	h := md5.New()
	io.WriteString(h, link)
	hash := string(h.Sum(nil))
	linkShort := fmt.Sprintf("%s%x", path, hash)
	dbHash := fmt.Sprintf("%x", hash)
	return linkShort, dbHash
}

func validateURL(l string) {
  isURL := govalidator.IsURL(link)
	if isURL != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	validateURL := govalidator.IsRequestURL(link)
	if validateURL != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}
