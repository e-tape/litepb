package generator

import (
	_ "embed"
	"strings"

	"github.com/e-tape/litepb/pkg/plugin"
)

func arr(values ...any) []any {
	return values
}

func dict(values ...any) map[any]any {
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

func lines(text string) []string {
	if text == "" {
		return nil
	}
	return strings.Split(text, "\n")
}

func replace(input string, values ...string) string {
	return strings.NewReplacer(values...).Replace(input)
}
