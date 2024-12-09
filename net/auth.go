package net

import (
	"time"

	"github.com/eagledb14/paperlink/types"
	"github.com/gofiber/fiber/v2"
)


func Auth(state *types.State, app *fiber.App) {

	app.Get("/login", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(BuildHtml("login.html", ""))
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		valid, err := state.Auth.ValidateUser(username, password)
		if !valid || err != nil {
			return c.SendStatus(404)
		}
		cookieString, _ := state.Auth.GenerateCookie()

		cookie := new(fiber.Cookie)
		cookie.Name = "session"
		cookie.Value = cookieString
		cookie.Expires = time.Now().Add(12 * time.Hour)
		cookie.HTTPOnly = true
		cookie.Secure = true
		cookie.SameSite = "Strict"

		state.Auth.Cookies[cookieString] = username

		c.Cookie(cookie)

		return c.SendString(BuildPage("/", "Engagements", BuildHtml("engagement_list.html", state.Engagements)))
	})

	app.Post("/create-user", func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		if len(password) < 12 {
			return c.SendStatus(404)
		}

		_, err := state.Auth.NewUser(username, password, false)
		if err != nil {
			return c.SendStatus(404)
		}
		cookieString, _ := state.Auth.GenerateCookie()

		cookie := new(fiber.Cookie)
		cookie.Name = "session"
		cookie.Value = cookieString
		cookie.Expires = time.Now().Add(12 * time.Hour)
		cookie.HTTPOnly = true
		cookie.Secure = true
		cookie.SameSite = "Strict"

		c.Cookie(cookie)
		
		state.Auth.Cookies[cookieString] = username

		return c.SendString(BuildPage("/", "Engagements", BuildHtml("engagement_list.html", state.Engagements)))
	})

	app.Post("/logout", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		setCookie := c.Cookies("session")

		c.ClearCookie()
		delete(state.Auth.Cookies, setCookie)

		return c.SendString(BuildHtml("login.html", ""))
	})

	app.Use(func(c *fiber.Ctx) error {
		cookie := c.Cookies("session")

		_, ok := state.Auth.Cookies[cookie]

		if !ok || cookie == "" {
			return c.Redirect("/login")
		}

		return c.Next()
	})
}
