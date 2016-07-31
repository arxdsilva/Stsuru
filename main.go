package main

import (
	"fmt"

	"crypto/md5"
	"io"

	"github.com/iris-contrib/template/django"
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type lines struct {
	Number int
	Link   string
	Short  string
}

func main() {
	iris.UseTemplate(django.New()).Directory("./templates", ".html")
	iris.Post("/link/add", addLink)
	iris.Post("/link/remove", removeLink)
	iris.Get("/", homer)
	iris.Listen(":8080")
}

func homer(ctx *iris.Context) {
	data := []lines{}
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
	h := md5.New()
	io.WriteString(h, link)
	hash := string(h.Sum(nil))
	linkshort := fmt.Sprintf("tsu.ru/%x", hash)
	number := 0
	linha := &lines{Number: number, Link: link, Short: linkshort}
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

func removeLink(ctx *iris.Context) {
	a := ctx.PostValue("remove")
	fmt.Print(a)
	// linha := &lines{Number: number, Link: link, Short: linkshort}
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		panic(err)
	}
	err = session.DB("tsuru").C("links").Remove(a)
	if err != nil {
		panic(err)
	}
	ctx.Redirect("/")
}
