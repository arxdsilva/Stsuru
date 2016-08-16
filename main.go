package main

import (
	"crypto/md5"
	"fmt"
	"io"

	"github.com/iris-contrib/template/django"
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type lines struct {
	Link  string
	Short string
}

func main() {
	iris.UseTemplate(django.New()).Directory("./templates", ".html")
	iris.Post("/link/add", addLink)
	iris.Get("/remove/link/:id", remover)
	iris.Get("/", homer)
	iris.Get("/:id", redirect)
	iris.Listen(":8080")
}

// Faz display da pagina home e carrega o html/DB
func homer(ctx *iris.Context) {
	// cria um slice de dicionarios do(a) tipo(estrutura) 'lines'
	data := []lines{}

	// DB conn
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		panic(err)
	}

	// DB display all data
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

	// prevents empty links
	if link == "" {
		ctx.Redirect("/")
		return
	}

	// inicia criacao e insercao do hash na url
	h := md5.New()
	io.WriteString(h, link)
	hash := string(h.Sum(nil))
	linkshort := fmt.Sprintf("tsu.ru:8080/%x", hash)

	linha := &lines{Link: link, Short: linkshort}
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
	// parsear na url > pegar id p/ query no mongo
	// como fazer o handler funcionar?

	id := ctx.Param("id")
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		panic(err)
	}
	// deleta do banco o valor 'id' recebido
	c := session.DB("tsuru").C("links")
	err = c.Remove(bson.M{"link": id})
	if err != nil {
		panic(err)
	}
	ctx.Redirect("/")
}

func redirect(ctx *iris.Context, id string) {
	// id := ctx.Param("id")
	// session, err := mgo.Dial("localhost")
	// defer session.Close()
	// if err != nil {
	// 	panic(err)
	// }
}
