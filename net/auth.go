package net

import (
	"fmt"
	"strings"
	"time"

	"github.com/eagledb14/paperlink/auth"
	"github.com/eagledb14/paperlink/types"
	"github.com/gofiber/fiber/v2"
)


func Auth(state *types.State, app *fiber.App) {

	app.Get("/login", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(BuildHtml("login.html", ""))
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		username := strings.Clone(c.FormValue("username"))
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
		cookie := c.Cookies("session")
		adminUsername := state.Auth.Cookies[cookie]
		adminUser, err := state.Auth.GetUser(adminUsername)
		if !adminUser.Admin || err != nil {
			return c.SendStatus(404)
		}

		username := strings.Clone(c.FormValue("username"))
		tempPassword, _ := state.Auth.GenerateCookie()

		newAdmin := c.FormValue("admin") == "on"

		_, err = state.Auth.NewUser(username, tempPassword, newAdmin)
		if err != nil {
			return c.SendStatus(404)
		}

		return c.SendString("Temporary Password: " + tempPassword)
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

	app.Delete("/account", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		cookie := c.Cookies("session")
		username := state.Auth.Cookies[cookie]
		password := c.Get("HX-Prompt")

		valid, err := state.Auth.ValidateUser(username, password)
		if !valid || err != nil {
			return c.SendStatus(404)
		}

		state.Auth.DeleteUser(username)
		delete(state.Auth.Cookies, cookie)

		return c.SendString(BuildHtml("login.html", ""))
	})

	app.Delete("/account-admin", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		cookie := c.Cookies("session")
		adminUsername := state.Auth.Cookies[cookie]
		adminUser, err := state.Auth.GetUser(adminUsername)

		if !adminUser.Admin || err != nil {
			return c.SendStatus(404)
		}

		username := c.FormValue("username")
		state.Auth.DeleteUser(username)

		return c.SendString("")
	})

	app.Put("/toggle-admin", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		cookie := c.Cookies("session")
		adminUsername := state.Auth.Cookies[cookie]
		adminUser, err := state.Auth.GetUser(adminUsername)

		if !adminUser.Admin || err != nil {
			return c.SendStatus(404)
		}

		username := c.FormValue("username")
		user, err := state.Auth.GetUser(username)
		if err != nil {

		}

		state.Auth.UpdateAdmin(user.Username, !user.Admin)

		users := state.Auth.GetUsers()
		data := struct {
			User auth.User
			Users []auth.User
		} {
			User: adminUser,
			Users: users,
		}

		return c.SendString(BuildPage("/ Profile /", adminUser.Username, BuildText("profile.html", data)))
	})

}
