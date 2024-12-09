package net

import (
	"strconv"

	"github.com/eagledb14/paperlink/dictionary"
	"github.com/eagledb14/paperlink/types"
	"github.com/gofiber/fiber/v2"
)


func Dictionary(state *types.State, app *fiber.App) {
	app.Get("/dictionary", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		data := struct {
			Words []dictionary.Word
		}{
			Words: state.Dictionary.GetWords(),
		}

		return c.SendString(BuildPage("/", "Dictionary", BuildText("dictionary_list.html", data)))
	})

	app.Get("/dictionary/new", func(c *fiber.Ctx) error {
		word := c.Get("HX-Prompt")

		data := struct {
			Word string
			Definition string
		}{
			Word: word,
			Definition: "",
		}

		return c.SendString(BuildPage("/", "Dictionary", BuildText("dictionary_new.html", data)))
	})

	app.Post("/dictionary/new", func(c *fiber.Ctx) error {
		word := c.FormValue("word")
		definition := c.FormValue("definition")
		state.Dictionary.InsertWord(word, definition)

		data := struct {
			Words []dictionary.Word
		}{
			Words: state.Dictionary.GetWords(),
		}

		return c.SendString(BuildPage("/", "Dictionary", BuildText("dictionary_list.html", data)))
	})

	app.Delete("/dictionary/:key", func(c *fiber.Ctx) error {
		key := c.Params("key")
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			return c.SendStatus(404)
		}
		state.Dictionary.Delete(keyInt)

		return c.SendString("")
	})
}
