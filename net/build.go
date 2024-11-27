package net

import (
	"bytes"
	html "html/template"
	text "text/template"
)

func BuildHtml(fileName string, data interface{}) string {
	tmpl, err := html.ParseFiles("./tmpl/" + fileName)
	if err != nil {
		return err.Error()
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, data)

	if err != nil {
		return err.Error()
	}

	return b.String()
}

func BuildText(fileName string, data interface{}) string {
	tmpl, err := text.ParseFiles("./tmpl/" + fileName)
	if err != nil {
		return err.Error()
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, data)

	if err != nil {
		return err.Error()
	}

	return b.String()
}

func BuildPage(title string, body string) string {
	data := struct {
		Title string
		Body string
	} {
		Title: title,
		Body: body,
	}

	return BuildText("build.html", data)
}
