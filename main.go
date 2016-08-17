package main

import (
	"crypto/md5"
	"fmt"
	"io"

	"github.com/iris-contrib/template/django"
	"github.com/kataras/iris"
	"github.com/tsuru/Stsuru/solver"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Lines struct {
	Link  string
	Short string
}

func main() {
	iris.UseTemplate(django.New()).Directory("./templates", ".html")
	iris.Post("/link/add", addLink)
	iris.Get("/remove/link/:id", remover)
	iris.Get("/", homer)
	iris.Get("/:id", solver.Redirect)
	iris.Listen(":8080")
}

func homer(ctx *iris.Context) {
	data := []Lines{}

	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		panic(err)
	}

	c := session.DB("tsuru").C("links")
	err = c.Find(bson.M{}).All(&data)
	if err != nil {
		panic(err)
	}

	context := map[string]interface{}{}
	context["array"] = data
	ctx.Render("mypage.html", context)
}

func addLink(ctx *iris.Context) {
	link := ctx.FormValueString("user_link")

	if link == "" {
		ctx.Redirect("/")
		return
	}

	h := md5.New()
	io.WriteString(h, link)
	hash := string(h.Sum(nil))
	linkshort := fmt.Sprintf("tsu.ru:8080/%x", hash)

	linha := &Lines{Link: link, Short: linkshort}
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		panic(err)
	}
	err = session.DB("tsuru").C("links").Insert(linha)
	if err != nil {
		panic(err)
	}
	ctx.Redirect("/")
}

func remover(ctx *iris.Context) {
	id := ctx.Param("id")
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		panic(err)
	}

	c := session.DB("tsuru").C("links")
	err = c.Remove(bson.M{"link": id})
	if err != nil {
		panic(err)
	}
	ctx.Redirect("/")
}
