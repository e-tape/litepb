package generator

import (
	"regexp"

	"github.com/e-tape/litepb/pkg/plugin"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// Generator of protobuf bindings
type (
	Generator struct {
		request    *pluginpb.CodeGeneratorRequest
		allFiles   map[Path]*generatorFile
		allTypes   map[Package]Type
		mapTypes   map[Package]*plugin.Message_Field_Type_Map
		aliasRegex *regexp.Regexp
	}
	Path    = string
	Package = string
	Type    struct {
		Name  string
		Alias string
	}
	// generatorFile of protobuf bindings for single file
	generatorFile struct {
		*Generator
		proto          *plugin.File
		sourceCodeInfo *descriptorpb.SourceCodeInfo
	}
)

func (a Type) reflect(alias string) *plugin.Message_Field_Type_Reflect {
	reflect := &plugin.Message_Field_Type_Reflect{
		Name: a.Name,
	}
	if a.Alias != alias {
		reflect.Path = a.Alias
	}
	return reflect
}
