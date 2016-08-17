package solver

import (
	"fmt"

	"github.com/arxdsilva/Stsuru/main"
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func redirect(ctx *iris.Context) {
	id := ctx.Param("id")
	dbData := []main.Lines{}

	session, err := mgo.Dial("localhost")
	defer session.Close()
	if err != nil {
		panic(err)
	}
	c := session.DB("tsuru").C("links")
	e := c.Find(bson.M{"link": id}).One(&dbData)
	if e != nil {
		ctx.Redirect("/")
	}
	fmt.Println(dbData["Link"])
	ctx.Redirect("/")
}
