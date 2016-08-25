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

	"github.com/gorilla/mux"
)

type lines struct {
	Link  string
	Short string
	Hash  string
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/link/add", addLink)
	r.HandleFunc("/link/remove", removeLink)
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
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

func addLink(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	link := r.Form["user_link"][0]

	if link == "" {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	// cria o Hash
	h := md5.New()
	io.WriteString(h, link)
	hash := string(h.Sum(nil))
	linkshort := fmt.Sprintf("http://tsu.ru:8080/%x", hash)
	dbHash := fmt.Sprintf("%x", hash)

	linha := &lines{Link: link, Short: linkshort, Hash: dbHash}
	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)

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

func removeLink(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)

	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)

	c := session.DB("tsuru").C("links")
	err = c.Remove(bson.M{"hash": id})
	checkError(err)

	http.Redirect(w, r, "/", http.StatusFound)
}

// func linkSolver(w http.ResponseWriter, r *http.Request) {
// 	// id := ctx.Param("id")
// 	dbData := lines{}
// 	// id = fmt.Sprintf("%s", id)
//
// 	session, err := mgo.Dial("localhost")
// 	defer session.Close()
// 	checkError(err)
//
// 	c := session.DB("tsuru").C("links").Find(bson.M{"hash": id}).One(&dbData)
// 	if c != nil {
// 		http.Redirect(w, r, "/", http.StatusFound)
// 	}
// 	// ctx.Redirect(dbData.Link)
// }
