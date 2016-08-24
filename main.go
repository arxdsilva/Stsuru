package main

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
)

type lines struct {
	Link  string
	Short string
	Hash  string
}

var t = pongo2.Must(pongo2.FromFile("my.html"))

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("NOT REAL")
	_, err := t.Execute(pongo2.Context{})
	if err != nil {
		log.Panic(err)
	}

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
}

// ____________________________________________________________________________________________
//
// func addLink(w http.ResponseWriter, r *http.Request) {
// 	// link := ctx.FormValueString("user_link")
// 	// if link == "" {
// 	// 	// ctx.Redirect("/")
// 	// }
//
// 	h := md5.New()
// 	// io.WriteString(h, link)
// 	hash := string(h.Sum(nil))
// 	linkshort := fmt.Sprintf("http://tsu.ru:8080/%x", hash)
// 	dbHash := fmt.Sprintf("%x", hash)
//
// 	// linha := &lines{Link: link, Short: linkshort, Hash: dbHash}
// 	session, err := mgo.Dial("localhost")
// 	defer session.Close()
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	err = session.DB("tsuru").C("links").Insert(linha)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	http.Redirect(w, r, "/", http.StatusFound)
// }
//
// func removeLink(w http.ResponseWriter, r *http.Request) {
// 	id := mux.Vars(r)
// 	// id := ctx.Param("id")
// 	session, err := mgo.Dial("localhost")
// 	defer session.Close()
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	c := session.DB("tsuru").C("links")
// 	err = c.Remove(bson.M{"hash": id})
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	// ctx.Redirect("/")
// }
//
// func linkSolver(w http.ResponseWriter, r *http.Request) {
// 	// id := ctx.Param("id")
// 	dbData := lines{}
// 	// id = fmt.Sprintf("%s", id)
//
// 	session, err := mgo.Dial("localhost")
// 	defer session.Close()
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	c := session.DB("tsuru").C("links").Find(bson.M{"hash": id}).One(&dbData)
// 	if c != nil {
// 		// ctx.Redirect("/")
// 	}
// 	// ctx.Redirect(dbData.Link)
// }
