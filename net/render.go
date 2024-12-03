package net

import (
	"strings"

	"github.com/eagledb14/paperlink/engagement"
)

func Render(sections []engagement.Section) string {
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

	return BuildText("render.html", data)
}
