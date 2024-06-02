package main

import (
	"bytes"
	"cmp"
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
		// return err // TODO: Enable
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
	definedTypes := make(map[string]string)               // Proto package + type name => Go package + type name
	protoFileToGoImport := make(map[string]GoImport)      // Proto file => Go import
	goImportPathToPackageAlias := make(map[string]string) // Go import path => package alias
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

		importPath := goPackage
		goPackage, goPackageName, ok := strings.Cut(goPackage, ";")
		if ok {
			importPath = goPackage
			goPackage = path.Dir(goPackage) + "/" + goPackageName
		} else {
			goPackageName = path.Base(goPackage)
		}

		goImport := GoImport{
			Path:  importPath,
			Alias: goPackageName,
		}
		protoFileToGoImport[protoFile.GetName()] = goImport
		goImportPathToPackageAlias[goImport.Path] = goImport.Alias

		protoFileName := path.Base(protoFile.GetName())
		protoFileNameExt := path.Ext(protoFileName)
		fileName := path.Join(importPath, strings.TrimSuffix(protoFileName, protoFileNameExt)+".pb.go")
		logf("\tGO FILE: %s", fileName)

		packagePrefix := "." + protoFile.GetPackage() + "."
		parseDefinedTypes(
			nil,
			protoFile.GetMessageType(),
			protoFile.GetEnumType(),
			packagePrefix,
			importPath,
			definedTypes,
		)

		packageAliases := make(map[string]string) // Go package => package alias
		mapTypes := make(map[string][2]string)    // Message name => [key type, value type]

		types, enumTypes := generateTypes(
			protoFile.GetMessageType(),
			protoFile.GetEnumType(),
			goImport,
			packagePrefix,
			definedTypes,
			mapTypes,
			packageAliases,
			goImportPathToPackageAlias,
			protoFile.GetSourceCodeInfo(),
			[]int32{4}, []int32{5},
		)

		imports := generateImports(protoFileToGoImport, packageAliases, goImport, protoFile.GetDependency())

		buf := bytes.NewBuffer(nil)
		err := goTemplate.Execute(buf, GoFile{
			Package:   goPackageName,
			Source:    protoFile.GetName(),
			Imports:   imports,
			Types:     types,
			EnumTypes: enumTypes,
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
	logf("FILE TO IMPORT PATH: %v", protoFileToGoImport)
	logf("DEFINED TYPES: %v", definedTypes)
	return &pluginpb.CodeGeneratorResponse{
		File: codeFiles,
	}
}

func generateImports(
	protoFileToGoImport map[string]GoImport, packageAliases map[string]string, goImport GoImport, dependencies []string,
) []GoImport {
	imports := make([]GoImport, 0, len(dependencies))
	for _, dependency := range dependencies {
		depGoImport, ok := protoFileToGoImport[dependency]
		if !ok {
			failf("missing dependency %s for %s", dependency, depGoImport.Path)
		}
		if depGoImport.Path == goImport.Path {
			continue
		}

		packageAlias, ok := packageAliases[depGoImport.Path]
		if !ok {
			packageAlias = "_"
		}

		imports = append(imports, GoImport{
			Path:  depGoImport.Path,
			Alias: packageAlias,
		})
	}
	slices.SortFunc(imports, func(a, b GoImport) int {
		return cmp.Compare(a.Path, b.Path)
	})
	return imports
}

func parseDefinedTypes(
	parentMessage *descriptorpb.DescriptorProto, messages []*descriptorpb.DescriptorProto,
	enums []*descriptorpb.EnumDescriptorProto, packagePrefix, importPath string, definedTypes map[string]string,
) {
	for _, message := range messages {
		if parentMessage != nil {
			message.Name = ptr(parentMessage.GetName() + "." + message.GetName())
		}
		definedTypes[packagePrefix+message.GetName()] = importPath + "." + strings.ReplaceAll(message.GetName(), ".", "_")
		parseDefinedTypes(message, message.GetNestedType(), message.GetEnumType(), packagePrefix, importPath, definedTypes)
	}

	for _, enum := range enums {
		if parentMessage != nil {
			enum.Name = ptr(parentMessage.GetName() + "." + enum.GetName())
		}
		definedTypes[packagePrefix+enum.GetName()] = importPath + "." + strings.ReplaceAll(enum.GetName(), ".", "_")
	}
}

func generateTypes(
	messages []*descriptorpb.DescriptorProto, enums []*descriptorpb.EnumDescriptorProto,
	goImport GoImport, packagePrefix string, definedTypes map[string]string, mapTypes map[string][2]string,
	packageAliases map[string]string, goImportPathToPackageAlias map[string]string,
	sourceCodeInfo *descriptorpb.SourceCodeInfo, msgSourceCodePath, enumSourceCodePath []int32,
) ([]GoType, []GoEnumType) {
	types := make([]GoType, 0, len(messages))
	enumTypes := make([]GoEnumType, 0, len(enums))

	for i, enum := range enums {
		values := make([]GoEnumTypeValue, 0, len(enum.GetValue()))
		for j, value := range enum.GetValue() {
			values = append(values, GoEnumTypeValue{
				Name:     value.GetName(),
				Comments: findEnumValueComments(sourceCodeInfo, enumSourceCodePath, i, j),
				Number:   value.GetNumber(),
			})
		}

		valuesPrefix := enum.GetName()
		if ix := strings.LastIndex(enum.GetName(), "."); ix >= 0 {
			valuesPrefix = enum.GetName()[:ix]
		}

		enumTypes = append(enumTypes, GoEnumType{
			Name:         strings.ReplaceAll(enum.GetName(), ".", "_"),
			Comments:     findEnumComments(sourceCodeInfo, enumSourceCodePath, i),
			ValuesPrefix: strings.ReplaceAll(valuesPrefix, ".", "_"),
			Values:       values,
		})
	}

	for i, message := range messages {
		if message.GetOptions().GetMapEntry() {
			mapTypes[packagePrefix+message.GetName()] = [2]string{
				fieldType(message.GetField()[0], definedTypes, mapTypes, packageAliases, goImportPathToPackageAlias, goImport),
				fieldType(message.GetField()[1], definedTypes, mapTypes, packageAliases, goImportPathToPackageAlias, goImport),
			}
			continue
		}

		nestedTypes, nestedEnumTypes := generateTypes(
			message.GetNestedType(),
			message.GetEnumType(),
			goImport,
			packagePrefix,
			definedTypes,
			mapTypes,
			packageAliases,
			goImportPathToPackageAlias,
			sourceCodeInfo,
			append(msgSourceCodePath, int32(i), 3),
			append(msgSourceCodePath, int32(i), 4),
		)

		fields := make([]GoTypeField, 0, len(message.GetField()))
		for j, field := range message.GetField() {
			typ := fieldType(field, definedTypes, mapTypes, packageAliases, goImportPathToPackageAlias, goImport)
			if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED && !strings.HasPrefix(typ, "map[") {
				typ = "[]" + typ
			}
			fields = append(fields, GoTypeField{
				Name:      snakeCaseToCamelCase(field.GetName()),
				Comments:  findMessageFieldComments(sourceCodeInfo, msgSourceCodePath, i, j),
				SnakeName: field.GetName(),
				Type:      typ,
			})
		}

		types = append(types, GoType{
			Name:     strings.ReplaceAll(message.GetName(), ".", "_"),
			Comments: findMessageComments(sourceCodeInfo, msgSourceCodePath, i),
			Fields:   fields,
		})

		types = append(types, nestedTypes...)
		enumTypes = append(enumTypes, nestedEnumTypes...)
	}

	return types, enumTypes
}

type GoFile struct {
	Package   string
	Source    string
	Imports   []GoImport
	Types     []GoType
	EnumTypes []GoEnumType
}

type GoImport struct {
	Path  string
	Alias string
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

type GoEnumType struct {
	Name         string
	Comments     string
	ValuesPrefix string
	Values       []GoEnumTypeValue
}

type GoEnumTypeValue struct {
	Name     string
	Comments string
	Number   int32
}

func findMessageComments(info *descriptorpb.SourceCodeInfo, sourceCodePath []int32, messageIndex int) string {
	return findComments(info, append(sourceCodePath, int32(messageIndex)))
}

func findMessageFieldComments(info *descriptorpb.SourceCodeInfo, sourceCodePath []int32, messageIndex, fieldIndex int) string {
	return findComments(info, append(sourceCodePath, int32(messageIndex), 2, int32(fieldIndex)))
}

func findEnumComments(info *descriptorpb.SourceCodeInfo, sourceCodePath []int32, enumIndex int) string {
	return findComments(info, append(sourceCodePath, int32(enumIndex)))
}

func findEnumValueComments(info *descriptorpb.SourceCodeInfo, sourceCodePath []int32, enumIndex, valueIndex int) string {
	return findComments(info, append(sourceCodePath, int32(enumIndex), 2, int32(valueIndex)))
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
	packageAliases map[string]string,
	goImportPathToPackageAlias map[string]string,
	goImport GoImport,
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
		failf("groups are not supported")
		return ""
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		if kv, ok := mapTypes[field.GetTypeName()]; ok {
			return "map[" + kv[0] + "]" + kv[1]
		}
		return fieldTypeMessageOrEnum(field, definedTypes, packageAliases, goImportPathToPackageAlias, goImport)
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return "[]byte"
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		return fieldTypeMessageOrEnum(field, definedTypes, packageAliases, goImportPathToPackageAlias, goImport)
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

func fieldTypeMessageOrEnum(
	field *descriptorpb.FieldDescriptorProto,
	definedTypes map[string]string,
	packageAliases map[string]string,
	goImportPathToPackageAlias map[string]string,
	goImport GoImport,
) string {
	goType, ok := definedTypes[field.GetTypeName()]
	if !ok {
		failf("unknown type %s", field.GetTypeName())
		return ""
	}

	ld := strings.LastIndex(goType, ".")
	typePackage := goType[:ld]
	typeName := goType[ld+1:]

	if typePackage == goImport.Path {
		return "*" + typeName
	}

	var packageAlias string
	if packageAlias, ok = packageAliases[typePackage]; !ok {
		importPath := goType[:strings.LastIndex(goType, ".")]
		packageAlias = goImportPathToPackageAlias[importPath]

		aliases := make([]string, 0, len(packageAliases))
		for _, alias := range packageAliases {
			aliases = append(aliases, alias)
		}

		for slices.Contains(aliases, packageAlias) {
			if importPath == "" || importPath == "." {
				panic("unreachable")
			}
			packageAlias = path.Base(importPath) + "_" + packageAlias
			importPath = path.Dir(importPath)
		}

		packageAliases[typePackage] = packageAlias
	}

	return "*" + packageAlias + "." + typeName
}
