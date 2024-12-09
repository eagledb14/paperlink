package net

import (
	"fmt"

	"github.com/eagledb14/paperlink/types"

	"github.com/gofiber/fiber/v2"
)

func Profile(state *types.State, app *fiber.App) {
	app.Get("/profile", func(c *fiber.Ctx) error {
		cookie := c.Cookies("session")

		username := state.Auth.Cookies[cookie]
		user, err := state.Auth.GetUser(username)
		fmt.Println(username, cookie)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(404)
		}

		return c.SendString(BuildPage("/ Profile /", user.Username, BuildText("profile.html", user)))
	})

}
