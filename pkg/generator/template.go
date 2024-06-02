package generator

import (
	"embed"
	_ "embed"
	"strings"
	"text/template"
)

//go:embed templates/*.tmpl
var goTemplateFiles embed.FS

var goTemplate = template.Must(
	template.New("").
		Funcs(goTemplateFunc).
		ParseFS(goTemplateFiles, "templates/*.tmpl"),
)

const mainTemplate = "main"

var additionalImports = []string{"fmt"}

var goTemplateFunc = template.FuncMap{
	"lines": func(text string) []string {
		if text == "" {
			return nil
		}
		return strings.Split(text, "\n")
	},
}
