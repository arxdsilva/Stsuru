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
	Number int
	Link   string
	Short  string
}

func main() {
	iris.UseTemplate(django.New()).Directory("./templates", ".html")
	iris.Post("/link/add", addLink)
	iris.Get("/remove/link", remover)
	iris.Get("/", homer)
	iris.Listen(":8080")
}

// Faz display da pagina home e carrega o html/DB
func homer(ctx *iris.Context) {
	// cria um array de dicionarios do(a) tipo(estrutura) 'lines'
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

func remover(ctx *iris.Context) {
	a := ctx.Param("remover")
	fmt.Println(a)
	// linha := &lines{Number: number, Link: link, Short: linkshort}
	// conexao com o db
	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		panic(err)
	}
	// deleta do banco o valor 'a' recebido
	err = session.DB("tsuru").C("links").Remove(bson.M{"link": a})
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(bson.M{"link": a})
		ctx.NotFound()
	}
	ctx.Redirect("/")
}
