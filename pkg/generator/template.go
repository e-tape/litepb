package generator

import (
	"bytes"
	"embed"
	_ "embed"
	"text/template"

	"github.com/e-tape/litepb/pkg/plugin"
)

const mainTemplate = "file"

type Template struct {
	tmpl  *template.Template
	proto *plugin.File
}

var (
	//go:embed templates/*.gotmpl
	goTemplateFiles embed.FS
	tmpl            = &Template{}
)

func init() {
	tmpl.tmpl = template.Must(
		template.New("").
			Funcs(template.FuncMap{
				"import":     addImport,
				"arr":        arr,
				"kv":         kv,
				"lines":      lines,
				"replace":    replace,
				"is_msg":     isMsg,
				"is_map":     isMap,
				"get_result": getResult,
				"set_result": setResult,
				"render":     render,
				"sort":       sort,
			}).
			ParseFS(goTemplateFiles, "templates/*.gotmpl"),
	)
}

func (a *Template) Execute(proto *plugin.File) (string, error) {
	buf := bytes.NewBuffer(nil)
	a.proto = proto
	if err := a.tmpl.ExecuteTemplate(buf, mainTemplate, proto); err != nil {
		return "", err
	}
	return buf.String(), nil
}
