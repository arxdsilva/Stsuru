package main

import (
	"github.com/kataras/iris"
)

func main() {
	iris.Get("/home/", home)
	iris.Listen(":8080")
}

func home(ctx *iris.Context) {
	ctx.Render("home.html", struct{ Link string }{Link: "www.globo.com"})
}
