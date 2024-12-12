package net

import (
	"bytes"
	html "html/template"
	"os"
	"strconv"
	"strings"
	text "text/template"
)

func BuildHtml(fileName string, data interface{}) string {
	file, err := os.ReadFile("./tmpl/" + fileName)
	tmpl, err := html.New("html_template").Funcs(getFuncMap()).Parse(string(file))
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
	file, err := os.ReadFile("./tmpl/" + fileName)
	tmpl, err := text.New("html_template").Funcs(getFuncMap()).Parse(string(file))
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

func BuildPage(directory string, title string, body string) string {
	data := struct {
		Directory string
		Title string
		Body string
	} {
		Directory: directory,
		Title: title,
		Body: body,
	}

	return BuildText("build.html", data)
}

func getFuncMap() html.FuncMap {
	return html.FuncMap {
		"append": func(i int, j string) string {
		    return strings.ReplaceAll(strconv.Itoa(i)+j, " ", "")
		},
	}
}
