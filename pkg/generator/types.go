package generator

import (
	"regexp"

	"github.com/e-tape/litepb/config"
	litepb "github.com/e-tape/litepb/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// Generator of protobuf bindings
type (
	Generator struct {
		cfg        config.Config
		request    *pluginpb.CodeGeneratorRequest
		allFiles   map[Path]*generatorFile
		allTypes   map[Package]*litepb.Message_Field_Type_Reflect
		mapTypes   map[Package]*litepb.Message_Field_Type_Map
		aliasRegex *regexp.Regexp
	}
	Path    = string
	Package = string
	// generatorFile of protobuf bindings for single file
	generatorFile struct {
		*Generator
		proto          *litepb.File
		messages       []*litepb.Message
		sourceCodeInfo *descriptorpb.SourceCodeInfo
	}
)
