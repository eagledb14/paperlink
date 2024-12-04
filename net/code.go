package net

import (
	"net/url"
	"strconv"

	"github.com/eagledb14/paperlink/engagement"
	"github.com/eagledb14/paperlink/types"
	"github.com/gofiber/fiber/v2"
)

func Code(state *types.State, app *fiber.App) {
	app.Get("/code/list/:name", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}

		e, err := state.GetEngagement(name)
		if err != nil {
			return c.SendStatus(404)
		}
		codes := e.GetCodes()

		data := struct {
			EngagementName string
			Codes []engagement.Code
		} {
			EngagementName: name,
			Codes: codes,
		}

		body := BuildText("code_list.html", data)

		return c.SendString(BuildPage("/ engagements / short codes / ", name, getCodeView(name, body)))
	})

	app.Get("/code/new/:name", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}

		body := BuildHtml("code_new.html", name)
		return c.SendString(BuildPage("/ engagements / short codes / ", name, getCodeView(name, body)))
	})

	app.Post("/code/new/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}

		e, err := state.GetEngagement(name)
		if err != nil {
			return c.SendStatus(404)
		}
		e.InsertCode(c.FormValue("code"), c.FormValue("paste"))

		codes := e.GetCodes()
		data := struct {
			EngagementName string
			Codes []engagement.Code
		} {
			EngagementName: name,
			Codes: codes,
		}
		body := BuildText("code_list.html", data)

		return c.SendString(BuildPage("/ engagements / short codes / ", name, getCodeView(name, body)))
	})

	app.Delete("/code/:name/:key", func(c *fiber.Ctx) error {
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}

		e, err := state.GetEngagement(name)
		if err != nil {
			return c.SendStatus(404)
		}

		key := c.Params("key")
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			return c.SendStatus(404)
		}

		e.DeleteCode(keyInt)

		return c.SendString("")
	})
}

func getCodeView(name string, body string) string {
	data := struct {
		Body string
		Name string
		Prompt string
		New string
	} {
		Body: body,
		New: "/code/new/" + name,
		Prompt: "New Short Code",
		Name: name ,
	}

	return BuildText("engagement_code_view.html", data)
}
