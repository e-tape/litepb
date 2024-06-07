package generator

import (
	"embed"
	_ "embed"
	"text/template"
)

const mainTemplate = "file"

var (
	//go:embed templates/*.gotmpl
	goTemplateFiles embed.FS
	goTemplate      = template.Must(
		template.New("").
			Funcs(template.FuncMap{
				"arr":     arr,
				"dict":    dict,
				"lines":   lines,
				"replace": replace,
				"is_msg":  isMsg,
				"is_map":  isMap,
			}).
			ParseFS(goTemplateFiles, "templates/*.gotmpl"),
	)
)
