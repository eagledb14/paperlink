package net

import (
	"github.com/eagledb14/paperlink/auth"
	"github.com/eagledb14/paperlink/types"

	"github.com/gofiber/fiber/v2"
)

func Profile(state *types.State, app *fiber.App) {
	app.Get("/profile", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		cookie := c.Cookies("session")

		username := state.Auth.Cookies[cookie]
		user, err := state.Auth.GetUser(username)
		if err != nil {
			return c.SendStatus(404)
		}

		users := state.Auth.GetUsers()
		data := struct {
			User auth.User
			Users []auth.User
		} {
			User: user,
			Users: users,
		}

		return c.SendString(BuildPage("/ Profile /", user.Username, BuildText("profile.html", data)))
	})

	app.Post("/profile/reset", func(c *fiber.Ctx) error {
		cookie := c.Cookies("session")
		name := state.Auth.Cookies[cookie]

		oldPass := c.FormValue("current")
		new1 := c.FormValue("new1")
		new2 := c.FormValue("new2")
		if new1 != new2 {
			return c.SendStatus(404)
		}
		if len(new1) < 12 {
			return c.SendStatus(404)
		}

		valid, err := state.Auth.ValidateUser(name, oldPass)
		if !valid || err != nil {
			return c.SendStatus(404)
		}

		user, err := state.Auth.GetUser(name)
		if err != nil {
			return c.SendStatus(404)
		}

		user.PassHash = new1
		newUser, err := state.Auth.UpdatePassword(user, new1)

		users := state.Auth.GetUsers()
		data := struct {
			User auth.User
			Users []auth.User
		} {
			User: newUser,
			Users: users,
		}

		return c.SendString(BuildPage("/ Profile /", newUser.Username, BuildText("profile.html", data)))
	})

}
