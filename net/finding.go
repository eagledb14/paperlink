package net

import (
	"net/url"
	"strconv"
	"time"

	"github.com/eagledb14/paperlink/dictionary"
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

		body := BuildText("finding_list.html", findingData)

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

		findingData := struct {
			Title string
			Assets []engagement.Asset
			EngagementName string
			Finding engagement.Finding
			Words []dictionary.Word
		} {
			Title: f.Title,
			Assets: e.GetAssets(),
			EngagementName: name,
			Finding: f,
			Words: state.Dictionary.GetWords(),
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

		assetKey := 0
		newAsset := c.FormValue("newAsset")
		if newAsset != "" {
			assetKey, err = e.InsertAsset("", newAsset,"")
		} else {
			assetKey, err = strconv.Atoi(c.FormValue("asset"))
		}
		newName := c.FormValue("name")

		dictionaryKey, err := strconv.Atoi(c.FormValue("dictionary"))
		severity, err := strconv.Atoi(c.FormValue("severity"))

		e.UpdateFinding(keyInt, severity, f.TimeStamp, newName, c.FormValue("body"), dictionaryKey, assetKey)
		
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

		findingData := struct {
			EngagementName string
			Finding engagement.Finding
			Asset engagement.Asset
			Word dictionary.Word
		} {
			EngagementName: name,
			Finding: f,
			Asset: e.GetAsset(f.AssetKey),
			Word: state.Dictionary.GetWord(f.DictionaryKey),
		}

		body := BuildText("finding_view.html", findingData)

		return c.SendString(BuildPage("/ engagements / findings / " + name + " /", f.Title, getFindingView(name, body)))
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

		findingData := struct {
			EngagementName string
			Findings []engagement.Finding
		} {
			EngagementName: name,
			Findings: e.GetFindings(),
		}

		body := BuildText("finding_list.html", findingData)

		return c.SendString(BuildPage("/ engagements / findings /", name, getFindingView(name, body)))
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
