package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type lines struct {
	Link  string
	Short string
	Hash  string
}

func main() {
	r := mux.NewRouter()
	// iris.UseTemplate(django.New()).Directory("./templates", ".html")
	r.HandleFunc("/", homer)
	r.HandleFunc("/links", addLink)
	r.HandleFunc("/links/remove/{id}", removeLink)
	r.HandleFunc("/links/{id}", linkSolver)
	http.Handle("/", r)
	http.ListenAndServe("tsu.ru:8080", nil)
}

func homer(w http.ResponseWriter, r *http.Request) {
	data := []lines{}

	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		log.Panic(err)
	}

	c := session.DB("tsuru").C("links")
	err = c.Find(bson.M{}).All(&data)
	if err != nil {
		log.Panic(err)
	}

	context := map[string]interface{}{}
	context["array"] = data
	// ctx.Render("mypage.html", context)
}

func addLink(w http.ResponseWriter, r *http.Request) {
	// link := ctx.FormValueString("user_link")
	if link == "" {
		// ctx.Redirect("/")
	}

	h := md5.New()
	io.WriteString(h, link)
	hash := string(h.Sum(nil))
	linkshort := fmt.Sprintf("http://tsu.ru:8080/%x", hash)
	dbHash := fmt.Sprintf("%x", hash)

	linha := &lines{Link: link, Short: linkshort, Hash: dbHash}
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		log.Panic(err)
	}
	err = session.DB("tsuru").C("links").Insert(linha)
	if err != nil {
		log.Panic(err)
	}
	// ctx.Redirect("/")
}

func removeLink(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	// id := ctx.Param("id")
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		log.Panic(err)
	}
	c := session.DB("tsuru").C("links")
	err = c.Remove(bson.M{"hash": id})
	if err != nil {
		log.Panic(err)
	}
	// ctx.Redirect("/")
}

func linkSolver(w http.ResponseWriter, r *http.Request) {
	// id := ctx.Param("id")
	dbData := lines{}
	id = fmt.Sprintf("%s", id)

	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		log.Panic(err)
	}
	c := session.DB("tsuru").C("links").Find(bson.M{"hash": id}).One(&dbData)
	if c != nil {
		// ctx.Redirect("/")
	}
	// ctx.Redirect(dbData.Link)
}
