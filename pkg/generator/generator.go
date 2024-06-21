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
		allTypes:   make(map[Package]*plugin.Message_Field_Type_Reflect),
		mapTypes:   make(map[Package]*plugin.Message_Field_Type_Map),
		aliasRegex: regexp.MustCompile(`(?mi)[^a-z0-9]`),
	}
}

// Generate generates bindings
func (a *Generator) Generate() *pluginpb.CodeGeneratorResponse {
	codeFiles := make([]*pluginpb.CodeGeneratorResponse_File, 0, len(a.request.GetProtoFile()))

	pluginData := &plugin.Plugin{}
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

		fg.generateMessages(
			protoFile.GetMessageType(),
			[]string{"", protoFile.GetPackage()},
			[]int32{4},
			nil,
		)
		fg.proto.Messages = fg.messages

		// TODO rathil del!!!
		//for _, msg := range fg.proto.Messages {
		//	msg.WithMemPool = false
		//}

		//fg.proto.Imports = fg.generateImports(protoFile.GetDependency())

		fg.proto.Generates = []plugin.File_Generate{
			plugin.File_STRUCT,
			plugin.File_INTERFACE,
			plugin.File_POOL,
			plugin.File_ENUM,
			plugin.File_NEW,
			plugin.File_RETURN_TO_POOL,
			plugin.File_PROTO_MESSAGE,
			plugin.File_CONVERT_TO,
			plugin.File_STRING,
			plugin.File_RESET,
			plugin.File_CLONE,
			plugin.File_GETTER,
			plugin.File_SETTER,
			plugin.File_SIZE,
		}

		marshal := *fg.proto // TODO rathil refactoring to litepb clone
		marshal.Name = strings.ToLower(strings.TrimSuffix(
			path.Base(protoFile.GetName()),
			path.Ext(protoFile.GetName()),
		)) + "_marshal.lpb.go"
		marshal.Generates = []plugin.File_Generate{
			plugin.File_MARSHAL,
		}

		unmarshal := *fg.proto // TODO rathil refactoring to litepb clone
		unmarshal.Name = strings.ToLower(strings.TrimSuffix(
			path.Base(protoFile.GetName()),
			path.Ext(protoFile.GetName()),
		)) + "_unmarshal.lpb.go"
		unmarshal.Generates = []plugin.File_Generate{
			plugin.File_UNMARSHAL,
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

		filePath := path.Join(
			file.GetPackage().GetDependency().GetPath(),
			file.GetName(),
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
