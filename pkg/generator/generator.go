package generator

import (
	"bytes"
	"cmp"
	"path"
	"slices"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/e-tape/litepb/pkg/common"
	"github.com/e-tape/litepb/pkg/stderr"
)

// Generator of protobuf bindings
type Generator struct {
	request             *pluginpb.CodeGeneratorRequest
	definedTypes        DefinedTypes
	protoFileToGoImport ProtoFileToGoImport
	goImportPathToAlias GoImportPathToAlias
}

// NewGenerator creates new generator
func NewGenerator(request *pluginpb.CodeGeneratorRequest) *Generator {
	return &Generator{
		request:             request,
		definedTypes:        make(DefinedTypes),
		protoFileToGoImport: make(ProtoFileToGoImport),
		goImportPathToAlias: make(GoImportPathToAlias),
	}
}

type (
	DefinedTypes  map[ProtoFullType]GoFullType
	ProtoFullType = string // Proto package + type name
	GoFullType    = string // Go package + type name
)

type (
	ProtoFileToGoImport map[ProtoFile]GoImport
	ProtoFile           = string
)

type (
	GoImportPathToAlias map[GoImportPath]GoImportAlias
	GoImportPath        = string
	GoImportAlias       = string
)

// FileGenerator of protobuf bindings for single file
type FileGenerator struct {
	*Generator
	importAliases      ImportAliases
	mapTypes           MapTypes
	goImport           GoImport
	protoPackagePrefix string
	sourceCodeInfo     *descriptorpb.SourceCodeInfo
}

func NewFileGenerator(g *Generator, protoFile *descriptorpb.FileDescriptorProto) *FileGenerator {
	goPackage := protoFile.GetOptions().GetGoPackage()
	if goPackage == "" {
		stderr.Failf("missing go_package option in %s", protoFile.GetName())
	}
	stderr.Logf("\tGO PACKAGE: %s", goPackage)

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

	importAliases := make(ImportAliases, len(additionalImports))
	for _, additionalImport := range additionalImports {
		importAliases[additionalImport] = strings.ReplaceAll(additionalImport, "/", "_")
	}

	return &FileGenerator{
		Generator:          g,
		importAliases:      importAliases,
		mapTypes:           make(MapTypes),
		goImport:           goImport,
		protoPackagePrefix: "." + protoFile.GetPackage() + ".",
		sourceCodeInfo:     protoFile.GetSourceCodeInfo(),
	}
}

type ImportAliases map[GoImportPath]GoImportAlias

type (
	MapTypes      map[MessageName]KeyValueTypes
	MessageName   = string
	KeyValueTypes = [2]string
)

// Generate generates bindings
func (g *Generator) Generate() *pluginpb.CodeGeneratorResponse {
	codeFiles := make([]*pluginpb.CodeGeneratorResponse_File, 0, len(g.request.GetProtoFile()))
	for _, protoFile := range g.request.ProtoFile {
		stderr.Logf("FILE START")
		stderr.Logf("\tNAME: %s", protoFile.GetName())
		stderr.Logf("\tPACKAGE: %s", protoFile.GetPackage())
		stderr.Logf("\tSYNTAX: %s", protoFile.GetSyntax())
		stderr.Logf("\tDEPENDENCIES: %s", strings.Join(protoFile.GetDependency(), ", "))

		fg := NewFileGenerator(g, protoFile)

		g.protoFileToGoImport[protoFile.GetName()] = fg.goImport
		g.goImportPathToAlias[fg.goImport.Path] = fg.goImport.Alias

		fg.defineTypes(
			nil,
			protoFile.GetMessageType(),
			protoFile.GetEnumType(),
		)

		types, enumTypes := fg.generateTypes(
			protoFile.GetMessageType(),
			protoFile.GetEnumType(),
			[]int32{4}, []int32{5},
		)

		imports := fg.generateImports(protoFile.GetDependency())

		buf := bytes.NewBuffer(nil)
		err := goTemplate.ExecuteTemplate(buf, mainTemplate, GoFile{
			Package:   fg.goImport.Alias,
			Source:    protoFile.GetName(),
			Imports:   imports,
			Types:     types,
			EnumTypes: enumTypes,
		})
		if err != nil {
			stderr.Failf("execute template: %s", err)
		}

		protoFileName := path.Base(protoFile.GetName())
		protoFileNameExt := path.Ext(protoFileName)
		fileName := path.Join(fg.goImport.Path, strings.TrimSuffix(protoFileName, protoFileNameExt)+".pb.go")
		stderr.Logf("\tGO FILE: %s", fileName)

		codeFiles = append(codeFiles, &pluginpb.CodeGeneratorResponse_File{
			Name:    &fileName,
			Content: common.Ptr(buf.String()),
		})

		stderr.Logf("FILE END")
	}
	stderr.Logf("FILE TO IMPORT PATH: %v", g.protoFileToGoImport)
	stderr.Logf("DEFINED TYPES: %v", g.definedTypes)
	return &pluginpb.CodeGeneratorResponse{
		File: codeFiles,
	}
}

func (g *FileGenerator) generateImports(dependencies []string) []GoImport {
	slices.Sort(dependencies)
	slices.Compact(dependencies)

	imports := make([]GoImport, 0, len(dependencies)+len(additionalImports))
	for _, dependency := range dependencies {
		depGoImport, ok := g.protoFileToGoImport[dependency]
		if !ok {
			stderr.Failf("missing dependency %s for %s", dependency, depGoImport.Path)
		}
		if depGoImport.Path == g.goImport.Path {
			continue
		}

		importAlias, ok := g.importAliases[depGoImport.Path]
		if !ok {
			importAlias = "_"
		}

		imports = append(imports, GoImport{
			Path:  depGoImport.Path,
			Alias: importAlias,
		})
	}

	for _, additionalImport := range additionalImports {
		imports = append(imports, GoImport{
			Path:  additionalImport,
			Alias: strings.ReplaceAll(additionalImport, "/", "_"),
		})
	}

	slices.SortFunc(imports, func(a, b GoImport) int {
		return cmp.Compare(a.Path, b.Path)
	})

	return imports
}

func (g *FileGenerator) defineTypes(
	parentMessage *descriptorpb.DescriptorProto,
	messages []*descriptorpb.DescriptorProto, enums []*descriptorpb.EnumDescriptorProto,
) {
	for _, message := range messages {
		if parentMessage != nil {
			message.Name = common.Ptr(parentMessage.GetName() + "." + message.GetName())
		}
		g.definedTypes[g.protoPackagePrefix+message.GetName()] = g.goImport.Path + "." +
			strings.ReplaceAll(message.GetName(), ".", "_")
		g.defineTypes(message, message.GetNestedType(), message.GetEnumType())
	}

	for _, enum := range enums {
		if parentMessage != nil {
			enum.Name = common.Ptr(parentMessage.GetName() + "." + enum.GetName())
		}
		g.definedTypes[g.protoPackagePrefix+enum.GetName()] = g.goImport.Path + "." +
			strings.ReplaceAll(enum.GetName(), ".", "_")
	}
}

func (g *FileGenerator) generateTypes(
	messages []*descriptorpb.DescriptorProto, enums []*descriptorpb.EnumDescriptorProto,
	msgSourceCodePath, enumSourceCodePath []int32,
) ([]GoType, []GoEnumType) {
	types := make([]GoType, 0, len(messages))
	enumTypes := make([]GoEnumType, 0, len(enums))

	for i, enum := range enums {
		values := make([]GoEnumTypeValue, 0, len(enum.GetValue()))
		for j, value := range enum.GetValue() {
			values = append(values, GoEnumTypeValue{
				Name:     value.GetName(),
				Comments: g.findEnumValueComments(enumSourceCodePath, i, j),
				Number:   value.GetNumber(),
			})
		}

		valuesPrefix := enum.GetName()
		if ix := strings.LastIndex(enum.GetName(), "."); ix >= 0 {
			valuesPrefix = enum.GetName()[:ix]
		}

		enumTypes = append(enumTypes, GoEnumType{
			Name:         strings.ReplaceAll(enum.GetName(), ".", "_"),
			Comments:     g.findEnumComments(enumSourceCodePath, i),
			ValuesPrefix: strings.ReplaceAll(valuesPrefix, ".", "_"),
			Values:       values,
		})
	}

	for i, message := range messages {
		if message.GetOptions().GetMapEntry() {
			g.mapTypes[g.protoPackagePrefix+message.GetName()] = [2]string{
				g.fieldType(message.GetField()[0]),
				g.fieldType(message.GetField()[1]),
			}
			continue
		}

		nestedTypes, nestedEnumTypes := g.generateTypes(
			message.GetNestedType(),
			message.GetEnumType(),
			append(msgSourceCodePath, int32(i), 3),
			append(msgSourceCodePath, int32(i), 4),
		)

		fields := make([]GoTypeField, 0, len(message.GetField()))
		for j, field := range message.GetField() {
			typ := g.fieldType(field)
			if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED &&
				!strings.HasPrefix(typ, "map[") {
				typ = "[]" + typ
			}
			fields = append(fields, GoTypeField{
				Name:      common.SnakeCaseToPascalCase(field.GetName()),
				Comments:  g.findMessageFieldComments(msgSourceCodePath, i, j),
				SnakeName: field.GetName(),
				Type:      typ,
				ZeroValue: g.fieldZeroValue(field),
			})
		}

		types = append(types, GoType{
			Name:     strings.ReplaceAll(message.GetName(), ".", "_"),
			Comments: g.findMessageComments(msgSourceCodePath, i),
			Fields:   fields,
		})

		types = append(types, nestedTypes...)
		enumTypes = append(enumTypes, nestedEnumTypes...)
	}

	return types, enumTypes
}

func (g *FileGenerator) fieldType(field *descriptorpb.FieldDescriptorProto) string {
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
		stderr.Failf("groups are not supported")
		return ""
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		if kv, ok := g.mapTypes[field.GetTypeName()]; ok {
			return "map[" + kv[0] + "]" + kv[1]
		}
		return "*" + g.fieldTypeMessageOrEnum(field)
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return "[]byte"
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		return g.fieldTypeMessageOrEnum(field)
	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
		return "int32"
	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
		return "int64"
	case descriptorpb.FieldDescriptorProto_TYPE_SINT32:
		return "int32"
	case descriptorpb.FieldDescriptorProto_TYPE_SINT64:
		return "int64"
	default:
		stderr.Failf("unknown type %d", field.GetType())
		return ""
	}
}

func (g *FileGenerator) fieldTypeMessageOrEnum(field *descriptorpb.FieldDescriptorProto) string {
	goType, ok := g.definedTypes[field.GetTypeName()]
	if !ok {
		stderr.Failf("unknown type %s", field.GetTypeName())
		return ""
	}

	ld := strings.LastIndex(goType, ".")
	typePackage := goType[:ld]
	typeName := goType[ld+1:]

	if typePackage == g.goImport.Path {
		return typeName
	}

	var importAlias string
	if importAlias, ok = g.importAliases[typePackage]; !ok {
		importPath := goType[:strings.LastIndex(goType, ".")]
		importAlias = g.goImportPathToAlias[importPath]

		aliases := make([]string, 0, len(g.importAliases))
		for _, alias := range g.importAliases {
			aliases = append(aliases, alias)
		}

		if slices.Contains(aliases, importAlias) {
			importPath = path.Dir(importPath)
		}
		for slices.Contains(aliases, importAlias) {
			if importPath == "" || importPath == "." {
				panic("unreachable")
			}
			importAlias = path.Base(importPath) + "_" + importAlias
			importPath = path.Dir(importPath)
		}

		g.importAliases[typePackage] = importAlias
	}

	return importAlias + "." + typeName
}

func (g *FileGenerator) fieldZeroValue(field *descriptorpb.FieldDescriptorProto) string {
	if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
		return "nil"
	}

	switch field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE, descriptorpb.FieldDescriptorProto_TYPE_FLOAT,
		descriptorpb.FieldDescriptorProto_TYPE_INT64, descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		descriptorpb.FieldDescriptorProto_TYPE_INT32, descriptorpb.FieldDescriptorProto_TYPE_UINT32,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64, descriptorpb.FieldDescriptorProto_TYPE_FIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32, descriptorpb.FieldDescriptorProto_TYPE_SFIXED64,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32, descriptorpb.FieldDescriptorProto_TYPE_SINT64:
		return "0"
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return "false"
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return `""`
	case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
		stderr.Failf("groups are not supported")
		return ""
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return "nil"
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		return "0"
	default:
		stderr.Failf("unknown type %d", field.GetType())
		return ""
	}
}
