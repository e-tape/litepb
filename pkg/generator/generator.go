package generator

import (
	"path"
	"regexp"
	"strings"

	"github.com/e-tape/litepb/pkg/plugin"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/e-tape/litepb/pkg/stderr"
)

// NewGenerator creates new generator
func NewGenerator(request *pluginpb.CodeGeneratorRequest) *Generator {
	return &Generator{
		request:    request,
		allFiles:   make(map[Path]*generatorFile),
		allTypes:   make(map[Package]Type),
		mapTypes:   make(map[Package]*plugin.Message_Field_Type_Map),
		aliasRegex: regexp.MustCompile(`(?mi)[^a-z0-9]`),
	}
}

// Generate generates bindings
func (a *Generator) Generate() *pluginpb.CodeGeneratorResponse {
	codeFiles := make([]*pluginpb.CodeGeneratorResponse_File, 0, len(a.request.GetProtoFile()))
	for _, protoFile := range a.request.ProtoFile {
		stderr.Logf("FILE START")
		stderr.Logf("\tNAME: %s", protoFile.GetName())
		stderr.Logf("\tPACKAGE: %s", protoFile.GetPackage())
		stderr.Logf("\tSYNTAX: %s", protoFile.GetSyntax())
		stderr.Logf("\tDEPENDENCIES: %s", strings.Join(protoFile.GetDependency(), ", "))

		fg := a.newFile(protoFile)
		a.allFiles[protoFile.GetName()] = fg

		fg.collectTypes(
			protoFile.GetEnumType(),
			protoFile.GetMessageType(),
			[]string{"", protoFile.GetPackage()},
			nil,
		)
		fg.collectMapTypes(
			protoFile.GetMessageType(),
			[]string{"", protoFile.GetPackage()},
		)

		fg.proto.Enums = fg.generateEnums(
			protoFile.GetEnumType(),
			protoFile.GetMessageType(),
			[]string{"", protoFile.GetPackage()},
			[]int32{4},
			[]int32{5},
			nil,
		)

		fg.proto.Messages = fg.generateMessages(
			protoFile.GetMessageType(),
			[]string{"", protoFile.GetPackage()},
			[]int32{4},
			nil,
		)

		fg.proto.Imports = fg.generateImports(protoFile.GetDependency())

		// run plugin

		content, err := tmpl.Execute(fg.proto)
		if err != nil {
			stderr.Failf("generate go file for proto [%s]: %s", protoFile.GetName(), err)
		}

		filePath := path.Join(
			fg.proto.GetPackage().GetDependency().GetPath(),
			fg.proto.GetName(),
		)
		stderr.Logf("\tGO FILE: %s", filePath)

		codeFiles = append(codeFiles, &pluginpb.CodeGeneratorResponse_File{
			Name:    &filePath,
			Content: &content,
		})

		stderr.Logf("FILE END")
	}
	//stderr.Logf("FILE TO IMPORT PATH: %v", a.protoFileToGoImport)
	definedTypes := make([]string, 0, len(a.allTypes))
	for _, t := range a.allTypes {
		definedTypes = append(definedTypes, t.Name)
	}
	stderr.Logf("DEFINED TYPES: %v", definedTypes)
	return &pluginpb.CodeGeneratorResponse{
		File: codeFiles,
	}
}
