package net

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/eagledb14/paperlink/engagement"
	"github.com/eagledb14/paperlink/types"
	"github.com/gofiber/fiber/v2"
)

func Template(state *types.State, app *fiber.App) {
	app.Get("/template", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(BuildPage("/", "Templates", BuildHtml("tmpl_list.html", state.Templates)))
	})

	app.Post("/template/new", func(c *fiber.Ctx) error {
		title := c.Get("HX-Prompt")
		t := engagement.NewTemplate(strings.Clone(title))
		state.AddTemplate(t)

		return c.Redirect("/template")
	})


	app.Get("/template/section/view/:name", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendString(BuildPage("/ templates / ", "Engagements", BuildHtml("engagement_list.html", state.Engagements)))
		}

		return c.SendString(BuildPage("/ templates / ", name, getTemplateView(state, name)))
	})

	app.Post("/template/section/new/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}
		e, err := state.GetTemplate(name)
		if err != nil {
			return c.Redirect("/template/section/view/"+name)
		}

		title := c.Get("HX-Prompt")
		err = e.InsertSection(title, "")

		return c.Redirect("/template/section/view/"+name)
	})

	app.Delete("/template/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendString(BuildPage("/ templates / ", "Engagements", BuildHtml("engagement_list.html", state.Engagements)))
		}

		state.DeleteTemplate(name)

		return c.SendString(BuildPage("/ templates / ", "Templates", BuildHtml("tmpl_list.html", state.Templates)))
	})

	app.Put("/template/section/body/:name/:key", func(c *fiber.Ctx) error {
		parser := struct {
			Content string `json:"content"`
		}{}

		if err := c.BodyParser(&parser); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid JSON",
			})
		}

		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}
		key := c.Params("key")

		e, err := state.GetTemplate(name)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}
		s := e.GetSection(key)
		e.UpdateSection(s.Key, s.Index, s.Title, parser.Content)

		return c.SendStatus(fiber.StatusOK)	
	})

	app.Delete("/template/section/:name/:key", func(c *fiber.Ctx) error {
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}
		key := c.Params("key")
		e, err  := state.GetTemplate(name)
		if err != nil {
			return c.SendStatus(fiber.StatusNoContent)
		}

		keyInt, err := strconv.Atoi(key)
		if err != nil {
			return c.SendStatus(fiber.StatusNoContent)
		}

		e.DeleteSection(keyInt)

		return c.SendString("")
	})

	app.Post("/template/section/update/:name/:key", func(c *fiber.Ctx) error {
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.Redirect("/template/section/view/"+name)
		}
		key := c.Params("key")
		e, err := state.GetTemplate(name)
		if err != nil {
			return c.Redirect("/template/section/view/"+name)
		}

		title := c.Get("HX-Prompt")
		section := e.GetSection(key)

		keyInt, err := strconv.Atoi(key)
		e.UpdateSection(keyInt, section.Index, title, section.Body)

		return c.Redirect("/template/section/view/"+name)
	})

	app.Put("/template/section/up/:name/:key", func(c *fiber.Ctx) error {
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}
		key := c.Params("key")

		e, err := state.GetTemplate(name)
		if err != nil {
			return c.SendString(BuildPage("/ templates / ", name, getTemplateView(state, name)))
		}

		bottomSection := e.GetSection(key)

		topSection := e.GetSectionFromIndex(bottomSection.Index - 1)
		if topSection.Key == 0 {
			return c.SendString(BuildPage("/ templates / ", name, getTemplateView(state, name)))
		}
		err = e.UpdateSection(bottomSection.Key, topSection.Index, bottomSection.Title, bottomSection.Body)
		err = e.UpdateSection(topSection.Key, bottomSection.Index, topSection.Title, topSection.Body)

		return c.SendString(BuildPage("/ templates / ", name, getTemplateView(state, name)))
	})

	app.Put("/template/section/down/:name/:key", func(c *fiber.Ctx) error {
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}
		key := c.Params("key")

		e, err := state.GetTemplate(name)
		if err != nil {
			return c.SendString(BuildPage("/ templates / ", name, getTemplateView(state, name)))
		}

		bottomSection := e.GetSection(key)

		topSection := e.GetSectionFromIndex(bottomSection.Index + 1)
		if topSection.Key == 0 {
			return c.SendString(BuildPage("/ templates / ", name, getTemplateView(state, name)))
		}
		err = e.UpdateSection(bottomSection.Key, topSection.Index, bottomSection.Title, bottomSection.Body)
		err = e.UpdateSection(topSection.Key, bottomSection.Index, topSection.Title, topSection.Body)

		return c.SendString(BuildPage("/ templates / ", name, getTemplateView(state, name)))
	})


}

func getTemplateView(state *types.State, name string) string {
	e, err := state.GetTemplate(name)
	if err != nil {
		return err.Error()
	}

	data := struct {
		Name string
		Sections []engagement.Section
	}{
		Name: e.Name,
		Sections: e.GetSections(),
	}

	return BuildText("tmpl_view.html", data)
}
