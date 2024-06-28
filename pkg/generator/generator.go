package generator

import (
	"regexp"
	"strings"

	"github.com/e-tape/litepb/config"
	litepb "github.com/e-tape/litepb/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/e-tape/litepb/pkg/stderr"
)

// NewGenerator creates new generator
func NewGenerator(
	cfg config.Config,
	request *pluginpb.CodeGeneratorRequest,
) *Generator {
	return &Generator{
		cfg:        cfg,
		request:    request,
		allFiles:   make(map[Path]*generatorFile),
		allTypes:   make(map[Package]*litepb.Message_Field_Type_Reflect),
		mapTypes:   make(map[Package]*litepb.Message_Field_Type_Map),
		aliasRegex: regexp.MustCompile(`(?mi)[^a-z0-9]`),
	}
}

// Generate generates bindings
func (a *Generator) Generate() *pluginpb.CodeGeneratorResponse {
	codeFiles := make([]*pluginpb.CodeGeneratorResponse_File, 0, len(a.request.GetProtoFile()))

	pluginData := &litepb.Plugin{}
	for _, protoFile := range a.request.ProtoFile {
		stderr.Logf("FILE START")
		stderr.Logf("\tNAME: %s", protoFile.GetName())
		stderr.Logf("\tPACKAGE: %s", protoFile.GetPackage())
		stderr.Logf("\tSYNTAX: %s", protoFile.GetSyntax())
		stderr.Logf("\tDEPENDENCIES: %s", strings.Join(protoFile.GetDependency(), ", "))

		a.generateMessagesFromExtensions(protoFile)
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

		fg.generateMessages(
			protoFile.GetMessageType(),
			[]string{"", protoFile.GetPackage()},
			[]int32{4},
			nil,
		)
		fg.proto.Messages = fg.messages

		// TODO rathil del!!!
		//for _, msg := range fg.proto.Messages {
		//	msg.MemPoolMessage = false
		//}

		//fg.proto.Imports = fg.generateImports(protoFile.GetDependency())

		fg.proto.Generates = []litepb.File_Generate{
			litepb.File_STRUCT,
			litepb.File_INTERFACE,
			litepb.File_POOL,
			litepb.File_ENUM,
			litepb.File_LIST,
			litepb.File_MAP,
			litepb.File_NEW,
			litepb.File_RETURN_TO_POOL,
			litepb.File_PROTO_MESSAGE,
			litepb.File_CONVERT_TO,
			litepb.File_STRING,
			litepb.File_RESET,
			litepb.File_CLONE,
			litepb.File_GETTER,
			litepb.File_SETTER,
			litepb.File_SIZE,
		}

		marshal := *fg.proto // TODO rathil refactoring to litepb clone
		marshal.Name = strings.ReplaceAll(marshal.Name, ".lpb.go", "_marshal.lpb.go")
		marshal.Generates = []litepb.File_Generate{
			litepb.File_MARSHAL,
		}

		unmarshal := *fg.proto // TODO rathil refactoring to litepb clone
		unmarshal.Name = strings.ReplaceAll(marshal.Name, ".lpb.go", "_unmarshal.lpb.go")
		unmarshal.Generates = []litepb.File_Generate{
			litepb.File_UNMARSHAL,
		}

		pluginData.Files = append(pluginData.Files, fg.proto, &marshal, &unmarshal)
	}

	// ENCODE plugin
	// TODO rathil run plugin
	// DECODE plugin

	//pluginData.Templates = append(pluginData.Templates, &plugin.Template{
	//	Name: "templates/func_name_message_new.gotmpl",
	//	Content: []byte(`
	//{{- define "func_name_message_new" -}}
	//   SuperNew{{ .GetName }}
	//{{- end -}}
	//`),
	//})

	for _, file := range pluginData.Files {
		content, err := tmpl.Execute(TemplateFs(pluginData.Templates), file)
		if err != nil {
			stderr.Failf("generate go file for proto [%s]: %s", file.GetName(), err)
		}

		stderr.Logf("\tGO FILE: %s", file.GetName())

		codeFiles = append(codeFiles, &pluginpb.CodeGeneratorResponse_File{
			Name:    &file.Name,
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
