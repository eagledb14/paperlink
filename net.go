package main

import (
	"os"
	"encoding/csv"
	"github.com/eagledb14/paperlink/net"
	"github.com/eagledb14/paperlink/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)


func Run() {
	EnsureCSVHasHeader()
	port := ":8080"
	state := types.NewState()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost" + port,
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.Redirect("/engagement")
	})

	app.Static("/style.css", "./tmpl/styles.css")
	app.Static("/tinymce.min.js", "./node_modules/tinymce/tinymce.min.js")
	app.Static("/themes/silver/theme.min.js", "./node_modules/tinymce/themes/silver/theme.min.js")
	app.Static("/plugins", "./node_modules/tinymce/plugins")
	app.Static("/models", "./node_modules/tinymce/models")
	app.Static("/icons", "./node_modules/tinymce/icons")
	app.Static("/skins", "./node_modules/tinymce/skins")

	net.Auth(state, app)
	net.Engagement(state, app)
	net.Section(state, app)
	net.Finding(state, app)
	net.Asset(state, app)
	net.Code(state, app)
	net.Dictionary(state, app)
	net.Template(state, app)
	net.Profile(state, app)

	app.Listen(port)
}

func EnsureCSVHasHeader() {
	// Check if the file exists
	_, err := os.Stat("access_logs.csv")
	if os.IsNotExist(err) {
		// If file doesn't exist, create it and write the header
		file, _ := os.Create("access_logs.csv")
		defer file.Close()

		// Create a CSV writer and write the header row
		writer := csv.NewWriter(file)
		defer writer.Flush()

		header := []string{"Username","Method", "Endpoint", "Timestamp", "HTTP Body"}
		writer.Write(header)
	}
}
