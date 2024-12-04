package net

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/eagledb14/paperlink/engagement"
	"github.com/eagledb14/paperlink/types"
	"github.com/gofiber/fiber/v2"
)

func Finding(state *types.State, app *fiber.App) {
	app.Get("/finding/list/:name", func(c *fiber.Ctx) error {
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

		findingData := struct {
			EngagementName string
			Findings []engagement.Finding
		} {
			EngagementName: name,
			Findings: e.GetFindings(),
		}

		body := BuildHtml("finding_list.html", findingData)

		return c.SendString(BuildPage("/ engagements / findings /", name, getFindingView(name, body)))
	})

	app.Post("/finding/new/:name", func(c *fiber.Ctx) error {
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
		key, _ := e.InsertFinding(0, time.Now(), title, "", 0, 0)

		return c.Redirect("/finding/edit/" + name + "/" + strconv.Itoa(key))
	})

	app.Get("/finding/edit/:name/:key", func(c *fiber.Ctx) error {
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
		f := e.GetFinding(keyInt)
		fmt.Println(keyInt, f)

		findingData := struct {
			Title string
			Assets []engagement.Asset
			EngagementName string
			Finding engagement.Finding
		} {
			Title: f.Title,
			Assets: e.GetAssets(),
			EngagementName: name,
			Finding: f,
		}
		
		body := BuildText("finding_edit.html", findingData)

		return c.SendString(BuildPage("/ engagements / narrative / ", name, getFindingView(name, body)))
	})

	app.Post("/finding/edit/:name/:key", func(c *fiber.Ctx) error {
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
		f := e.GetFinding(keyInt)
		_ = f

		assetKey, err := strconv.Atoi(c.FormValue("asset"))
		dictionaryKey, err := strconv.Atoi(c.FormValue("dictionary"))
		severity, err := strconv.Atoi(c.FormValue("severity"))

		e.UpdateFinding(keyInt, severity, f.TimeStamp, f.Title, c.FormValue("body"), dictionaryKey, assetKey)
		
		return c.Redirect("/finding/view/" + name + "/" + key)
	})

	app.Get("/finding/view/:name/:key", func(c *fiber.Ctx) error {
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
		f := e.GetFinding(keyInt)
		_ = f


		findingData := struct {
			EngagementName string
			Finding engagement.Finding
			Asset engagement.Asset
		} {
			EngagementName: name,
			Finding: f,
			Asset: e.GetAsset(f.AssetKey),
		}

		body := BuildHtml("finding_view.html", findingData)

		return c.SendString(BuildPage("/ engagements / findings / " + name, f.Title, getFindingView(name, body)))
	})

	app.Delete("/finding/:name/:key", func(c *fiber.Ctx) error {
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
		e.DeleteFinding(keyInt)
		return c.Redirect("/finding/list"+name)
	})

}

func getFindingView(name string, body string) string {

	data := struct {
		Body string
		Name string
		Prompt string
		New string
	} {
		Body: body,
		New: "/finding/new/" + name,
		Prompt: "New Finding",
		Name: name ,
	}

	return BuildText("engagement_view.html", data)
}