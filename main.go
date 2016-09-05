package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

type lines struct {
	Link  string
	Short string
	Hash  string
}

func main() {
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

func css(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/out/", http.FileServer(http.Dir("/out/home.css")))
}

// Home ...
func Home(w http.ResponseWriter, r *http.Request) {
	Data := []lines{}
	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)
	c := session.DB("tsuru").C("links")
	err = c.Find(bson.M{}).All(&Data)
	checkError(err)
	t, err := template.ParseFiles("tmpl/index.html")
	checkError(err)

	t.Execute(w, Data)
}

// AddLink ...
func AddLink(w http.ResponseWriter, r *http.Request) {
	// gets the URL inserted in Form
	r.ParseForm()
	link := r.Form["user_link"][0]

	// checking the URL
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

	// URL hashing
	h := md5.New()
	io.WriteString(h, link)
	hash := string(h.Sum(nil))
	linkshort := fmt.Sprintf("http://localhost:8080/%x", hash)
	dbHash := fmt.Sprintf("%x", hash)

	// URL storage
	linha := &lines{Link: link, Short: linkshort, Hash: dbHash}
	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)

	dbData := lines{}
	err = session.DB("tsuru").C("links").Find(bson.M{"hash": dbHash}).One(&dbData)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err = session.DB("tsuru").C("links").Insert(linha)
	checkError(err)

	http.Redirect(w, r, "/", http.StatusFound)
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
	return
}

// RemoveLink ...
func RemoveLink(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	idInfo := id["id"]

	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)
	c := session.DB("tsuru").C("links")
	err = c.Remove(bson.M{"hash": idInfo})
	checkError(err)
	http.Redirect(w, r, "/", http.StatusFound)
}

// LinkSolver ...
func LinkSolver(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	dbData := lines{}
	idInfo := id["id"]

	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)
	err = session.DB("tsuru").C("links").Find(bson.M{"hash": idInfo}).One(&dbData)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	http.Redirect(w, r, dbData.Link, http.StatusFound)
}
