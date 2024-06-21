package generator

import (
	"bytes"
	"cmp"
	"slices"
	"strings"

	"github.com/e-tape/litepb/pkg/plugin"
)

func addImport(path string, alias ...string) string {
	if path == "" || path == tmpl.proto.GetPackage().GetDependency().GetPath() {
		return ""
	}
	for _, imp := range tmpl.proto.Imports {
		if imp.Path == path {
			if imp.Alias != "" {
				return imp.Alias
			}
			parts := strings.Split(imp.Path, "/")
			return parts[len(parts)-1]
		}
	}
	var impAlias string
	if len(alias) > 0 {
		impAlias = alias[0]
	} else {
		impAlias = strings.ReplaceAll(path, "/", "_")
	}
	if path == impAlias {
		impAlias = ""
	}
	tmpl.proto.Imports = append(tmpl.proto.Imports, &plugin.Dependency{
		Path:  path,
		Alias: impAlias,
	})
	if impAlias != "" {
		return impAlias
	}
	return path
}

func arr(values ...any) []any {
	return values
}

func kv(values ...any) map[any]any {
	result := make(map[any]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		result[values[i]] = values[i+1]
	}
	return result
}

func isMsg(fieldType *plugin.Message_Field_Type) bool {
	return fieldType.GetInProto() == plugin.Message_Field_Type_MESSAGE_OR_MAP && fieldType.GetMap() == nil
}

func isMap(fieldType *plugin.Message_Field_Type) bool {
	return fieldType.GetInProto() == plugin.Message_Field_Type_MESSAGE_OR_MAP && fieldType.GetMap() != nil
}

func isGenerate(generate string) bool {
	for _, f := range tmpl.proto.GetGenerates() {
		if f.String() == generate {
			return true
		}
	}
	return false
}

func lines(text string) []string {
	if text == "" {
		return nil
	}
	return strings.Split(text, "\n")
}

func replace(input string, values ...string) string {
	return strings.NewReplacer(values...).Replace(input)
}

func getResult(value map[any]any) any {
	return value["result"]
}

func setResult(value map[any]any, result any) any {
	value["result"] = result
	return ""
}

func render(name string, data any) (string, error) {
	var result bytes.Buffer
	if err := tmpl.tmpl.ExecuteTemplate(&result, name, data); err != nil {
		return "", err
	}
	return result.String(), nil
}

func sort(items any) any {
	switch tItems := items.(type) {
	case []string:
		slices.SortFunc(tItems, cmp.Compare[string])
		return tItems
	case []*plugin.Dependency:
		slices.SortFunc(tItems, func(a, b *plugin.Dependency) int {
			return cmp.Compare(a.Alias, b.Alias)
		})
		return slices.CompactFunc(tItems, func(a, b *plugin.Dependency) bool {
			if a.Alias == "" && b.Alias == "" {
				return a.Path == b.Path
			}
			return a.Alias == b.Alias
		})
	}
	return items
}
