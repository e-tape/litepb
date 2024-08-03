package generator

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"text/template"

	litepb "github.com/e-tape/litepb/proto"
)

const mainTemplate = "file"

type Template struct {
	tmpl  *template.Template
	proto *litepb.File
}

var (
	//go:embed templates/*.gotmpl
	nativeTemplateFiles embed.FS
	tmpl                = &Template{}
)

func (a *Template) Execute(
	tmplFs fs.FS,
	proto *litepb.File,
) (string, error) {
	// TODO rathil once
	tmplFile, err := template.New("").
		Funcs(template.FuncMap{
			"import":               addImport,
			"arr":                  arr,
			"append":               arrAppend,
			"kv":                   kv,
			"lines":                lines,
			"replace":              replace,
			"is_msg":               isMsg,
			"is_map":               isMap,
			"is_generate":          isGenerate,
			"set":                  set,
			"get":                  get,
			"render":               render,
			"sort":                 sort,
			"add":                  add,
			"sub":                  sub,
			"mul":                  mul,
			"pack_field_num_bytes": packFieldNumBytes,
			"pack_field_num_int":   packFieldNumInt,
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
