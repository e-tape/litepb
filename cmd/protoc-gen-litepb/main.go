package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

//go:embed templates/go.tmpl
var goTemplateFile string

var goTemplate = template.Must(template.New("").Parse(goTemplateFile))

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
	for _, protoFile := range request.ProtoFile {
		logf("\tNAME: %s", protoFile.GetName())
		logf("\tPACKAGE: %s", protoFile.GetPackage())
		logf("\tSYNTAX: %s", protoFile.GetSyntax())
		logf("\tEDITION: %s", protoFile.GetEdition())

		goPackage := protoFile.GetOptions().GetGoPackage()
		if goPackage == "" {
			return &pluginpb.CodeGeneratorResponse{
				Error: ptr(fmt.Sprintf("missing go_package option in %s", protoFile.GetName())),
			}
		}
		logf("\tGO PACKAGE: %s", goPackage)

		protoFileName := path.Base(protoFile.GetName())
		protoFileNameExt := path.Ext(protoFileName)
		fileName := path.Join(path.Dir(goPackage), strings.TrimSuffix(protoFileName, protoFileNameExt)+".go")

		types := make([]GoType, 0, len(protoFile.GetMessageType()))
		for _, message := range protoFile.GetMessageType() {
			fields := make([]GoTypeField, 0, len(message.GetField()))
			for _, field := range message.GetField() {
				fields = append(fields, GoTypeField{
					Name: field.GetName(),
					Type: fieldType(field.GetType()),
				})
			}

			types = append(types, GoType{
				Name:   message.GetName(),
				Fields: fields,
			})
		}

		buf := bytes.NewBuffer(nil)
		err := goTemplate.Execute(buf, GoFile{
			Package: path.Base(goPackage),
			Source:  protoFile.GetName(),
			Types:   types,
		})
		if err != nil {
			failf("execute template: %s", err)
		}

		codeFiles = append(codeFiles, &pluginpb.CodeGeneratorResponse_File{
			Name:    &fileName,
			Content: ptr(buf.String()),
		})
	}
	return &pluginpb.CodeGeneratorResponse{
		File: codeFiles,
	}
}

type GoFile struct {
	Package string
	Source  string
	Types   []GoType
}

type GoType struct {
	Name   string
	Fields []GoTypeField
}

type GoTypeField struct {
	Name string
	Type string
}

func fieldType(typ descriptorpb.FieldDescriptorProto_Type) string {
	switch typ {
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
		panic("unimplemented")
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
		failf("unknown type %d", typ)
		return ""
	}
}
