package generator

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
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
	nativeTemplateFiles embed.FS
	tmpl                = &Template{}
)

func (a *Template) Execute(
	tmplFs fs.FS,
	proto *plugin.File,
) (string, error) {
	// TODO rathil once
	tmplFile, err := template.New("").
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
		ParseFS(tmplFs, "templates/*.gotmpl")
	if err != nil {
		return "", fmt.Errorf("create template, err: %w", err)
	}
	a.tmpl = tmplFile
	buf := bytes.NewBuffer(nil)
	a.proto = proto
	if err = a.tmpl.ExecuteTemplate(buf, mainTemplate, proto); err != nil {
		return "", err
	}
	return buf.String(), nil
}
