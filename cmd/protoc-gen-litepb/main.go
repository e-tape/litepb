package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

//go:embed templates/go.tmpl
var goTemplateFile string

var goTemplate = template.Must(template.New("").Funcs(goTemplateFunc).Parse(goTemplateFile))

var goTemplateFunc = template.FuncMap{
	"lines": func(text string) []string {
		if text == "" {
			return nil
		}
		return strings.Split(text, "\n")
	},
}

func main() {
	if err := run(); err != nil {
		failf("%s: %s", filepath.Base(os.Args[0]), err)
	}
}

func run() error {
	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	request := &pluginpb.CodeGeneratorRequest{}
	if err = proto.Unmarshal(in, request); err != nil {
		return err
	}

	logf("COMPILER: %s", request.GetCompilerVersion())
	logf("FILES TO GENERATE: %s", strings.Join(request.GetFileToGenerate(), ", "))

	start := time.Now()
	response := generate(request)
	if err = goFmt(response); err != nil {
		return err
	}
	logf("GENERATED IN: %s", time.Since(start))

	out, err := proto.Marshal(response)
	if err != nil {
		return err
	}

	if _, err = os.Stdout.Write(out); err != nil {
		return err
	}

	return nil
}

func generate(request *pluginpb.CodeGeneratorRequest) *pluginpb.CodeGeneratorResponse {
	codeFiles := make([]*pluginpb.CodeGeneratorResponse_File, 0, len(request.GetProtoFile()))
	definedTypes := make(map[string]string)     // Proto package + type name => Go package + type name
	fileToImportPath := make(map[string]string) // Proto file => Go import path
	for _, protoFile := range request.ProtoFile {
		logf("FILE START")
		logf("\tNAME: %s", protoFile.GetName())
		logf("\tPACKAGE: %s", protoFile.GetPackage())
		logf("\tSYNTAX: %s", protoFile.GetSyntax())
		logf("\tDEPENDENCIES: %s", strings.Join(protoFile.GetDependency(), ", "))

		goPackage := protoFile.GetOptions().GetGoPackage()
		if goPackage == "" {
			return &pluginpb.CodeGeneratorResponse{
				Error: ptr(fmt.Sprintf("missing go_package option in %s", protoFile.GetName())),
			}
		}
		logf("\tGO PACKAGE: %s", goPackage)

		fileToImportPath[protoFile.GetName()] = goPackage

		protoFileName := path.Base(protoFile.GetName())
		protoFileNameExt := path.Ext(protoFileName)
		fileName := path.Join(goPackage, strings.TrimSuffix(protoFileName, protoFileNameExt)+".go")
		logf("\tGO FILE: %s", fileName)

		packagePrefix := "." + protoFile.GetPackage() + "."
		parseDefinedMessages(nil, protoFile.GetMessageType(), packagePrefix, goPackage, definedTypes)

		mapTypes := make(map[string][2]string) // Message name => [key type, value type]
		types := generateTypes(
			protoFile.GetMessageType(),
			goPackage,
			packagePrefix,
			definedTypes,
			mapTypes,
			protoFile.GetSourceCodeInfo(),
			[]int32{4},
		)

		imports := make([]string, 0, len(protoFile.GetDependency()))
		for _, dependency := range protoFile.GetDependency() {
			importPath, ok := fileToImportPath[dependency]
			if !ok {
				return &pluginpb.CodeGeneratorResponse{
					Error: ptr(fmt.Sprintf("missing Go dependency %s for %s", dependency, protoFile.GetName())),
				}
			}
			if importPath == goPackage {
				continue
			}
			imports = append(imports, importPath)
		}
		slices.Sort(imports)

		buf := bytes.NewBuffer(nil)
		err := goTemplate.Execute(buf, GoFile{
			Package: path.Base(goPackage),
			Source:  protoFile.GetName(),
			Imports: imports,
			Types:   types,
		})
		if err != nil {
			failf("execute template: %s", err)
		}

		codeFiles = append(codeFiles, &pluginpb.CodeGeneratorResponse_File{
			Name:    &fileName,
			Content: ptr(buf.String()),
		})

		logf("FILE END")
	}
	logf("FILE TO IMPORT PATH: %v", fileToImportPath)
	logf("DEFINED TYPES: %v", definedTypes)
	return &pluginpb.CodeGeneratorResponse{
		File: codeFiles,
	}
}

func parseDefinedMessages(
	parentMessage *descriptorpb.DescriptorProto, messages []*descriptorpb.DescriptorProto, packagePrefix string,
	goPackage string, definedTypes map[string]string,
) {
	for _, message := range messages {
		if parentMessage != nil {
			message.Name = ptr(parentMessage.GetName() + "." + message.GetName())
		}
		definedTypes[packagePrefix+message.GetName()] = goPackage + "." + strings.ReplaceAll(message.GetName(), ".", "_")
		parseDefinedMessages(message, message.GetNestedType(), packagePrefix, goPackage, definedTypes)
	}
}

func generateTypes(
	messages []*descriptorpb.DescriptorProto, goPackage, packagePrefix string, definedTypes map[string]string,
	mapTypes map[string][2]string, sourceCodeInfo *descriptorpb.SourceCodeInfo, sourceCodePath []int32,
) []GoType {
	types := make([]GoType, 0, len(messages))
	for i, message := range messages {
		if message.GetOptions().GetMapEntry() {
			mapTypes[packagePrefix+message.GetName()] = [2]string{
				fieldType(message.GetField()[0], definedTypes, mapTypes, goPackage),
				fieldType(message.GetField()[1], definedTypes, mapTypes, goPackage),
			}
			continue
		}

		nestedTypes := generateTypes(
			message.GetNestedType(), goPackage, packagePrefix, definedTypes, mapTypes, sourceCodeInfo,
			append(sourceCodePath, int32(i), 3),
		)

		fields := make([]GoTypeField, 0, len(message.GetField()))
		for j, field := range message.GetField() {
			typ := fieldType(field, definedTypes, mapTypes, goPackage)
			if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED && !strings.HasPrefix(typ, "map[") {
				typ = "[]" + typ
			}
			fields = append(fields, GoTypeField{
				Name:      snakeCaseToCamelCase(field.GetName()),
				Comments:  findMessageFieldComments(sourceCodeInfo, sourceCodePath, i, j),
				SnakeName: field.GetName(),
				Type:      typ,
			})
		}

		types = append(types, GoType{
			Name:     strings.ReplaceAll(message.GetName(), ".", "_"),
			Comments: findMessageComments(sourceCodeInfo, sourceCodePath, i),
			Fields:   fields,
		})

		types = append(types, nestedTypes...)
	}
	return types
}

type GoFile struct {
	Package string
	Source  string
	Imports []string
	Types   []GoType
}

type GoType struct {
	Name     string
	Comments string
	Fields   []GoTypeField
}

type GoTypeField struct {
	Name      string
	Comments  string
	SnakeName string
	Type      string
}

func findMessageComments(info *descriptorpb.SourceCodeInfo, sourceCodePath []int32, messageIndex int) string {
	return findComments(info, append(sourceCodePath, int32(messageIndex)))
}

func findMessageFieldComments(info *descriptorpb.SourceCodeInfo, sourceCodePath []int32, messageIndex, fieldIndex int) string {
	return findComments(info, append(sourceCodePath, int32(messageIndex), 2, int32(fieldIndex)))
}

func findComments(info *descriptorpb.SourceCodeInfo, ps []int32) string {
	for _, loc := range info.GetLocation() {
		if slices.Equal(loc.GetPath(), ps) {
			return strings.TrimSuffix(loc.GetLeadingComments()+loc.GetTrailingComments(), "\n")
		}
	}
	return ""
}

func fieldType(
	field *descriptorpb.FieldDescriptorProto,
	definedTypes map[string]string,
	mapTypes map[string][2]string,
	goPackage string,
) string {
	switch field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		return "float64"
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		return "float32"
	case descriptorpb.FieldDescriptorProto_TYPE_INT64:
		return "int64"
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64:
		return "uint64"
	case descriptorpb.FieldDescriptorProto_TYPE_INT32:
		return "int32"
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32:
		return "uint32"
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		return "uint64"
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		return "uint32"
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return "bool"
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return "string"
	case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
		panic("unimplemented")
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		goType, ok := definedTypes[field.GetTypeName()]
		if !ok {
			failf("unknown message type %s", field.GetTypeName())
			return ""
		}

		var kv [2]string
		if kv, ok = mapTypes[field.GetTypeName()]; ok {
			return "map[" + kv[0] + "]" + kv[1]
		}

		if strings.HasPrefix(goType, goPackage+".") {
			goType = goType[strings.LastIndex(goType, ".")+1:]
		} else {
			goType = goType[strings.LastIndex(goType, "/")+1:]
		}
		return "*" + goType
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return "[]byte"
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		panic("unimplemented")
	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
		return "int32"
	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
		return "int64"
	case descriptorpb.FieldDescriptorProto_TYPE_SINT32:
		return "int32"
	case descriptorpb.FieldDescriptorProto_TYPE_SINT64:
		return "int64"
	default:
		failf("unknown type %d", field.GetType())
		return ""
	}
}
