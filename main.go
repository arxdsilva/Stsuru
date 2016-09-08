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
	http.StripPrefix("/out/", http.FileServer(http.Dir("out/")))
}

// Home handles "/" GET request and loads all data
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

// AddLink handles POST request to DB and redirects to Home
func AddLink(w http.ResponseWriter, r *http.Request) {
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

	path := "http://localhost:8080/"
	// URL hashing
	linkShort, dbHash := hash(link, path)

	l := &lines{Link: link, Short: linkShort, Hash: dbHash}

	_, err := findOne(dbHash)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err = insert(l)
	checkError(err)

	http.Redirect(w, r, "/", http.StatusFound)
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
	return
}

// RemoveLink searches db for a certain link & removes It if It exists
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
	idInfo := id["id"]

	l, err := findOne(idInfo)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	http.Redirect(w, r, l, http.StatusFound)
}

func hash(link, path string) (string, string) {
	h := md5.New()
	io.WriteString(h, link)
	hash := string(h.Sum(nil))
	linkShort := fmt.Sprintf("%s%x", path, hash)
	dbHash := fmt.Sprintf("%x", hash)
	return linkShort, dbHash
}

func findOne(dbHash string) (string, error) {
	dbData := lines{}
	session, err := mgo.Dial("localhost")
	checkError(err)
	defer session.Close()
	err = session.DB("tsuru").C("links").Find(bson.M{"hash": dbHash}).One(&dbData)
	return dbData.Link, err
}

func insert(l *lines) error {
	session, err := mgo.Dial("localhost")
	defer session.Close()
	err = session.DB("tsuru").C("links").Insert(l)
	return err
}
