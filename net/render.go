package net

import (
	"strings"

	"github.com/eagledb14/paperlink/engagement"
)

func Render(sections []engagement.Section, codes []engagement.Code) string {
	builder := strings.Builder{}

	for _, section := range sections {
		builder.WriteString("\n<br>\n")
		builder.WriteString("<h2>"+ section.Title +"</h2>")
		builder.WriteString("\n<br>\n")
		builder.WriteString(section.Body)
	}

	data := struct {
		Document string
	} {
		Document: builder.String(),
	}

	document := BuildText("render.html", data)

	for _, code := range codes {
		document = strings.ReplaceAll(document, "%%" + code.Code + "%%", code.Paste)
	}

	return document
}
