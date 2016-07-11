package main

import "github.com/kataras/iris"

func main() {
	iris.Config.Render.Template.Engine = iris.PongoEngine
	iris.Get("/home", home)
	iris.Listen(":8080")
}

func home(ctx *iris.Context) {
	ctx.Render("home.html", map[string]interface{}{
		"Link":  "iris",
		"Short": "ISIS",
	})
}
