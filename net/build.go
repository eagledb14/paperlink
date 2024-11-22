package net

import (
	"bytes"
	html "html/template"
	text "text/template"
)

func BuildHtml(fileName string, data interface{}) string {
	tmpl, err := html.ParseFiles("./templates/" + fileName)
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
	tmpl, err := text.ParseFiles("./templates/" + fileName)
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

func BuildPage(body string) string {
	data := struct {
		Body string
	} {
		Body: "this is the body",
	}

	return BuildText("build.html", data)
}
