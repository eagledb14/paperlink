package net

import (

	"github.com/eagledb14/paperlink/types"
	"github.com/gofiber/fiber/v2"
)


func Run() {
	state := types.NewState()
	_ = state

	app := fiber.New()
	_ = app

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		// return c.SendString(BuildPage(""))
		return c.Redirect("/engagement")
	})

	Engagement(state, app)
	Section(state, app)
	Template(state, app)

	app.Static("/style.css", "./tmpl/styles.css")

	app.Static("/tinymce.min.js", "./resources/tinymce/tinymce.min.js")
	app.Static("/themes/silver/theme.min.js", "./resources/tinymce/themes/silver/theme.min.js")
	app.Static("/plugins", "./resources/tinymce/plugins")
	app.Static("/models", "./resources/tinymce/models")
	app.Static("/icons", "./resources/tinymce/icons")
	app.Static("/skins", "./resources/tinymce/skins")

	app.Listen(":8080")
}


