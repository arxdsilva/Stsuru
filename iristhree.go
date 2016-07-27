package main

import "github.com/kataras/iris"

type shortener struct {
	Input  string
	Button string
}

func m() {
	iris.Get("/hi", home)
	iris.Listen(":8080")
}

func home(ctx *iris.Context) {
	ctx.Render("home1.html", shortener{}{})
}
