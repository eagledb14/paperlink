package main

import (
	"github.com/eagledb14/paperlink/net"
	"github.com/eagledb14/paperlink/types"
	"github.com/gofiber/fiber/v2"
)


func Run() {
	state := types.NewState()

	app := fiber.New()


	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.Redirect("/engagement")
	})

	net.Engagement(state, app)
	net.Section(state, app)
	net.Finding(state, app)
	net.Asset(state, app)
	net.Code(state, app)
	net.Dictionary(state, app)
	net.Template(state, app)

	app.Static("/style.css", "./tmpl/styles.css")

	app.Static("/tinymce.min.js", "./node_modules/tinymce/tinymce.min.js")
	app.Static("/themes/silver/theme.min.js", "./node_modules/tinymce/themes/silver/theme.min.js")
	app.Static("/plugins", "./node_modules/tinymce/plugins")
	app.Static("/models", "./node_modules/tinymce/models")
	app.Static("/icons", "./node_modules/tinymce/icons")
	app.Static("/skins", "./node_modules/tinymce/skins")

	app.Listen(":8080")
}


