package net

import (
	"net/url"
	"strconv"

	"github.com/eagledb14/paperlink/engagement"
	"github.com/eagledb14/paperlink/types"

	"github.com/gofiber/fiber/v2"
)

func Asset(state *types.State, app *fiber.App) {
	app.Get("/asset/list/:name", func(c *fiber.Ctx) error {
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

		data := struct {
			EngagementName string
			Assets []engagement.Asset
		} {
			EngagementName: name,
			Assets: e.GetAssets(),
		}
		body := BuildText("asset_list.html", data)

		return c.SendString(BuildPage("/ engagements / assets / ", name, getAssetView(name, body)))
	})

	app.Post("/asset/new/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}

		e, err := state.GetEngagement(name)
		if err != nil {
			return c.SendString(err.Error())
		}
		title := c.Get("HX-Prompt")
		key, _ := e.InsertAsset("", title, "")

		return c.Redirect("/asset/edit/" + name + "/" + strconv.Itoa(key))
	})

	app.Get("/asset/edit/:name/:key", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}

		key := c.Params("key")
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			return c.SendStatus(404)
		}
		
		e, err := state.GetEngagement(name)
		if err != nil {
			return c.SendString(err.Error())
		}
		a := e.GetAsset(keyInt)

		data := struct {
			EngagementName string
			Asset engagement.Asset
		} {
			EngagementName: name,
			Asset: a,
		}
		body := BuildText("asset_edit.html", data)

		return c.SendString(BuildPage("/ engagements / assets / ", name, getAssetView(name, body)))
	})

	app.Post("/asset/edit/:name/:key", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}

		key := c.Params("key")
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			return c.SendStatus(404)
		}
		
		e, err := state.GetEngagement(name)
		if err != nil {
			return c.SendString(err.Error())
		}

		newName := c.FormValue("name")

		e.UpdateAsset(keyInt, c.FormValue("parent"), newName, c.FormValue("type"))
		return c.Redirect("/asset/view/" + name + "/" + key)
	})

	app.Get("/asset/view/:name/:key", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}

		key := c.Params("key")
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			return c.SendStatus(404)
		}
		
		e, err := state.GetEngagement(name)
		if err != nil {
			return c.SendString(err.Error())
		}
		a := e.GetAsset(keyInt)
		data := struct {
			EngagementName string
			Asset engagement.Asset
			Findings []engagement.Finding
		} {
			EngagementName: name,
			Asset: a,
			Findings: e.GetFindingsWithAsset(a.Key),
		}
		body := BuildText("asset_view.html", data)

		return c.SendString(BuildPage("/ engagements / assets / " + name + " /", a.Name, getAssetView(name, body)))
	})
	
	app.Delete("asset/:name/:key", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		name := c.Params("name")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return c.SendStatus(404)
		}

		key := c.Params("key")
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			return c.SendStatus(404)
		}

		e, err := state.GetEngagement(name)
		if err != nil {
			return c.SendString(err.Error())
		}
		e.DeleteAsset(keyInt)
		e.DeleteFindingsWithAsset(keyInt)

		data := struct {
			EngagementName string
			Assets []engagement.Asset
		} {
			EngagementName: name,
			Assets: e.GetAssets(),
		}
		body := BuildText("asset_list.html", data)

		return c.SendString(BuildPage("/ engagements / assets / ", name, getAssetView(name, body)))
	})
}

func getAssetView(name string, body string) string {
	data := struct {
		Body string
		Name string
		Prompt string
		New string
	} {
		Body: body,
		New: "/asset/new/" + name,
		Prompt: "New Asset",
		Name: name ,
	}

	return BuildText("engagement_view.html", data)
}
