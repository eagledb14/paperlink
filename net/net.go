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

	app.Static("/style.css", "./tmpl/styles.css")
	app.Listen(":8080")
}


