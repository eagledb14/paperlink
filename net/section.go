package net

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/eagledb14/paperlink/engagement"
	"github.com/eagledb14/paperlink/types"
	"github.com/gofiber/fiber/v2"
)

func Section(state *types.State, app *fiber.App) {
	app.Get("/section/view/:name", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}

		return c.SendString(BuildPage("/ engagements / narrative / ", name, getSectionView(state, name)))
	})

	app.Put("/section/body/:name/:key", func(c *fiber.Ctx) error {
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

		e, err := state.GetEngagement(name)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}
		s := e.GetSection(key)
		e.UpdateSection(s.Key, s.Index, s.Title, parser.Content)

		return c.SendStatus(fiber.StatusOK)	
	})

	app.Delete("/section/:name/:key", func(c *fiber.Ctx) error {
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}
		key := c.Params("key")
		e, err  := state.GetEngagement(name)
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			return c.SendStatus(fiber.StatusNoContent)
		}

		e.DeleteSection(keyInt)

		return c.SendString("")
	})

	app.Post("/section/new/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}
		e, err := state.GetEngagement(name)
		if err != nil {
			return c.Redirect("/section/view/"+name)
		}

		title := c.Get("HX-Prompt")
		e.InsertSection(strings.Clone(title), "")

		return c.Redirect("/section/view/"+name)
	})

	app.Post("/section/update/:name/:key", func(c *fiber.Ctx) error {
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}
		key := c.Params("key")
		e, err := state.GetEngagement(name)
		if err != nil {
			return c.Redirect("/section/view/"+name)
		}

		title := c.Get("HX-Prompt")
		section := e.GetSection(key)

		keyInt, err := strconv.Atoi(key)
		e.UpdateSection(keyInt, section.Index, title, section.Body)

		return c.Redirect("/section/view/"+name)
	})

	app.Put("/section/up/:name/:key", func(c *fiber.Ctx) error {
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}
		key := c.Params("key")

		e, err := state.GetEngagement(name)
		if err != nil {
			return c.SendString(BuildPage("/ engagements / narrative / ", name, getSectionView(state, name)))
		}

		bottomSection := e.GetSection(key)
		topSection := e.GetSectionFromIndex(bottomSection.Index - 1)
		if topSection.Key == 0 {
			return c.SendString(BuildPage("/ engagements / narrative / ", name, getSectionView(state, name)))
		}
		err = e.UpdateSection(bottomSection.Key, topSection.Index, bottomSection.Title, bottomSection.Body)
		err = e.UpdateSection(topSection.Key, bottomSection.Index, topSection.Title, topSection.Body)

		return c.SendString(BuildPage("/ engagements / narrative / ", name, getSectionView(state, name)))
	})

	app.Put("/section/down/:name/:key", func(c *fiber.Ctx) error {
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}
		key := c.Params("key")

		e, err := state.GetEngagement(name)
		if err != nil {
			return c.SendString(BuildPage("/ engagements / narrative / ", name, getSectionView(state, name)))
		}

		topSection := e.GetSection(key)

		bottomSection := e.GetSectionFromIndex(topSection.Index + 1)
		if bottomSection.Key == 0 {
			return c.SendString(BuildPage("/ engagements / narrative / ", name, getSectionView(state, name)))
		}
		err = e.UpdateSection(topSection.Key, bottomSection.Index, topSection.Title, topSection.Body)
		err = e.UpdateSection(bottomSection.Key, topSection.Index, bottomSection.Title, bottomSection.Body)

		return c.SendString(BuildPage("/ engagements / narrative / ", name, getSectionView(state, name)))
	})
}

func getSectionView(state *types.State, name string) string {
	e, err := state.GetEngagement(name)
	if err != nil {
		return err.Error()
	}

	sectionData := struct {
		EngagementName string
		Sections []engagement.Section
	}{
		EngagementName: e.Name,
		Sections: e.GetSections(),
	}

	data := struct {
		Body string
		Name string
		Prompt string
		New string
	} {
		Body: BuildHtml("section_view.html", sectionData),
		New: "/section/new/" + e.Name,
		Prompt: "New Section",
		Name: e.Name,
	}

	return BuildText("engagement_view.html", data)
}
