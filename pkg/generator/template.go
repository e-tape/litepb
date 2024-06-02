package generator

import (
	_ "embed"
	"strings"
	"text/template"
)

//go:embed templates/go.tmpl
var goTemplateFile string

var goTemplate = template.Must(template.New("").Funcs(goTemplateFunc).Parse(goTemplateFile))

var goTemplateFunc = template.FuncMap{
	"lines": func(text string) []string {
		if text == "" {
			return nil
		}
		return strings.Split(text, "\n")
	},
}
