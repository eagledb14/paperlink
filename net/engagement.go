package net

import (
	"net/url"

	"github.com/eagledb14/paperlink/engagement"
	"github.com/eagledb14/paperlink/types"
	"github.com/gofiber/fiber/v2"
)


func Engagement(state *types.State, app *fiber.App) {
	app.Get("/engagement", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		return c.SendString(BuildPage("Engagements", BuildHtml("engagement_list.html", state.Engagements)))
	})

	app.Get("/engagement/new", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(BuildPage("Engagements", BuildHtml("new_engagement.html", struct{}{})))
	})

	app.Post("/engagement/new", func(c *fiber.Ctx) error {
		name := c.FormValue("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}
		contact := c.FormValue("contact")
		email := c.FormValue("email")

		newEngagement := engagement.NewEngagement(name, contact, email)

		state.AddEngagement(newEngagement)
		return c.Redirect("/section/view/"+newEngagement.Name)
	})

	app.Get("/engagement/template", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(BuildPage("Engagements", BuildHtml("template_engagement.html", state.Templates)))
	})

	app.Post("/engagement/template", func(c *fiber.Ctx) error {
		name := c.FormValue("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}
		contact := c.FormValue("contact")
		email := c.FormValue("email")
		template := c.FormValue("template")

		newEngagement := engagement.NewEngagementFromTemplate(template, name, contact, email)

		state.AddEngagement(newEngagement)
		return c.Redirect("/section/view/"+newEngagement.Name)
	})

	app.Delete("/engagement/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendString(BuildPage("Engagements", BuildHtml("engagement_list.html", state.Engagements)))
		}

		state.DeleteEnagement(name)

		return c.SendString(BuildPage("Engagements", BuildHtml("engagement_list.html", state.Engagements)))
	})


}

