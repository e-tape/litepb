package generator

import (
	"bytes"
	"cmp"
	"regexp"
	"slices"
	"strings"

	litepb "github.com/e-tape/litepb/proto"
)

var aliasRegex = regexp.MustCompile("(?mi)[^a-z0-9_]")

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
		impAlias = aliasRegex.ReplaceAllString(path, "_")
	}
	if path == impAlias {
		impAlias = ""
	}
	tmpl.proto.Imports = append(tmpl.proto.Imports, &litepb.Dependency{
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

func arrAppend(arr []any, values ...any) []any {
	return append(arr, values...)
}

func kv(values ...any) map[any]any {
	result := make(map[any]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		result[values[i]] = values[i+1]
	}
	return result
}

func isMsg(fieldType *litepb.Message_Field_Type) bool {
	return fieldType.GetInProto() == litepb.Message_Field_Type_MESSAGE_OR_MAP && fieldType.GetMap() == nil
}

func isMap(fieldType *litepb.Message_Field_Type) bool {
	return fieldType.GetInProto() == litepb.Message_Field_Type_MESSAGE_OR_MAP && fieldType.GetMap() != nil
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

func set(value any, k any, v any) any {
	switch vt := value.(type) {
	case map[any]any:
		vt[k] = v
	case []any:
		kt := k.(int)
		if len(vt) > kt {
			vt[kt] = v
		}
	}
	return ""
}

func get(value any, k any, defaultValue ...any) any {
	switch vt := value.(type) {
	case map[any]any:
		if v, ok := vt[k]; ok {
			return v
		}
	case []any:
		kt := k.(int)
		if len(vt) > kt {
			return vt[kt]
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
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
	case []*litepb.Dependency:
		slices.SortFunc(tItems, func(a, b *litepb.Dependency) int {
			return cmp.Compare(a.Alias, b.Alias)
		})
		return slices.CompactFunc(tItems, func(a, b *litepb.Dependency) bool {
			if a.Alias == "" && b.Alias == "" {
				return a.Path == b.Path
			}
			return a.Alias == b.Alias
		})
	}
	return items
}

func add(x, y int) int { return x + y }
func sub(x, y int) int { return x - y }
func mul(x, y int) int { return x * y }

func packFieldNumBytes(
	number int32,
	inProtoAny any,
	packed bool,
) []byte {
	var wireType int32
	var inProto litepb.Message_Field_Type_Proto
	switch it := inProtoAny.(type) {
	case litepb.Message_Field_Type_Proto:
		inProto = it
	case int:
		inProto = litepb.Message_Field_Type_Proto(it)
	}
	switch litepb.Message_Field_Type_Proto(inProto) {
	case litepb.Message_Field_Type_INT32,
		litepb.Message_Field_Type_INT64,
		litepb.Message_Field_Type_UINT32,
		litepb.Message_Field_Type_UINT64,
		litepb.Message_Field_Type_SINT32,
		litepb.Message_Field_Type_SINT64,
		litepb.Message_Field_Type_BOOL:
		if !packed {
			wireType = 0
		} else {
			wireType = 2
		}
	case litepb.Message_Field_Type_ENUM:
		wireType = 0
	case litepb.Message_Field_Type_FIXED64,
		litepb.Message_Field_Type_SFIXED64,
		litepb.Message_Field_Type_DOUBLE:
		if !packed {
			wireType = 1
		} else {
			wireType = 2
		}
	case litepb.Message_Field_Type_STRING,
		litepb.Message_Field_Type_BYTES,
		litepb.Message_Field_Type_MESSAGE_OR_MAP:
		wireType = 2
	case litepb.Message_Field_Type_FIXED32,
		litepb.Message_Field_Type_SFIXED32,
		litepb.Message_Field_Type_FLOAT:
		if !packed {
			wireType = 5
		} else {
			wireType = 2
		}
	}

	num := number<<3 | wireType
	result := make([]byte, 0, 4)
	for num >= 1<<7 {
		result = append(result, byte(num&127|128))
		num >>= 7
	}
	return append(result, byte(num))
}

func packFieldNumInt(
	number int32,
	inProto any,
	packed bool,
) any {
	data := packFieldNumBytes(number, inProto, packed)
	switch len(data) {
	case 1:
		return data[0]
	case 2:
		return uint16(data[0]) | uint16(data[1])<<8
	case 3:
		return uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16
	default:
		return uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
	}
}
