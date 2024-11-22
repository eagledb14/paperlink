package net

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)


func Run() {
	fmt.Println("hi")

	app := fiber.New()
	_ = app

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(BuildPage(""))
	})

	app.Static("/style.css", "./templates/styles.css")
	app.Listen(":8080")
}


