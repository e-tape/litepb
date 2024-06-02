package generator

import (
	"bytes"
	"cmp"
	"fmt"
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
	request                    *pluginpb.CodeGeneratorRequest
	definedTypes               DefinedTypes
	protoFileToGoImport        ProtoFileToGoImport
	goImportPathToPackageAlias GoImportPathToPackageAlias
}

// NewGenerator creates new generator
func NewGenerator(request *pluginpb.CodeGeneratorRequest) *Generator {
	return &Generator{
		request:                    request,
		definedTypes:               make(DefinedTypes),
		protoFileToGoImport:        make(ProtoFileToGoImport),
		goImportPathToPackageAlias: make(GoImportPathToPackageAlias),
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
	GoImportPathToPackageAlias map[GoImportPath]PackageAlias
	GoImportPath               = string
	PackageAlias               = string
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

		goPackage := protoFile.GetOptions().GetGoPackage()
		if goPackage == "" {
			return &pluginpb.CodeGeneratorResponse{
				Error: common.Ptr(fmt.Sprintf("missing go_package option in %s", protoFile.GetName())),
			}
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
		g.protoFileToGoImport[protoFile.GetName()] = goImport
		g.goImportPathToPackageAlias[goImport.Path] = goImport.Alias

		protoFileName := path.Base(protoFile.GetName())
		protoFileNameExt := path.Ext(protoFileName)
		fileName := path.Join(importPath, strings.TrimSuffix(protoFileName, protoFileNameExt)+".pb.go")
		stderr.Logf("\tGO FILE: %s", fileName)

		packagePrefix := "." + protoFile.GetPackage() + "."
		g.parseDefinedTypes(
			nil,
			protoFile.GetMessageType(),
			protoFile.GetEnumType(),
			packagePrefix,
			importPath,
		)

		packageAliases := make(map[string]string) // Go package => package alias
		mapTypes := make(map[string][2]string)    // Message name => [key type, value type]

		types, enumTypes := g.generateTypes(
			protoFile.GetMessageType(),
			protoFile.GetEnumType(),
			goImport,
			packagePrefix,
			mapTypes,
			packageAliases,
			protoFile.GetSourceCodeInfo(),
			[]int32{4}, []int32{5},
		)

		imports := g.generateImports(packageAliases, goImport, protoFile.GetDependency())

		buf := bytes.NewBuffer(nil)
		err := goTemplate.Execute(buf, GoFile{
			Package:   goPackageName,
			Source:    protoFile.GetName(),
			Imports:   imports,
			Types:     types,
			EnumTypes: enumTypes,
		})
		if err != nil {
			stderr.Failf("execute template: %s", err)
		}

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

func (g *Generator) generateImports(
	packageAliases map[string]string, goImport GoImport, dependencies []string,
) []GoImport {
	imports := make([]GoImport, 0, len(dependencies))
	for _, dependency := range dependencies {
		depGoImport, ok := g.protoFileToGoImport[dependency]
		if !ok {
			stderr.Failf("missing dependency %s for %s", dependency, depGoImport.Path)
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

func (g *Generator) parseDefinedTypes(
	parentMessage *descriptorpb.DescriptorProto, messages []*descriptorpb.DescriptorProto,
	enums []*descriptorpb.EnumDescriptorProto, packagePrefix, importPath string,
) {
	for _, message := range messages {
		if parentMessage != nil {
			message.Name = common.Ptr(parentMessage.GetName() + "." + message.GetName())
		}
		g.definedTypes[packagePrefix+message.GetName()] = importPath + "." + strings.ReplaceAll(message.GetName(), ".", "_")
		g.parseDefinedTypes(message, message.GetNestedType(), message.GetEnumType(), packagePrefix, importPath)
	}

	for _, enum := range enums {
		if parentMessage != nil {
			enum.Name = common.Ptr(parentMessage.GetName() + "." + enum.GetName())
		}
		g.definedTypes[packagePrefix+enum.GetName()] = importPath + "." + strings.ReplaceAll(enum.GetName(), ".", "_")
	}
}

func (g *Generator) generateTypes(
	messages []*descriptorpb.DescriptorProto, enums []*descriptorpb.EnumDescriptorProto,
	goImport GoImport, packagePrefix string, mapTypes map[string][2]string,
	packageAliases map[string]string,
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
				g.fieldType(message.GetField()[0], mapTypes, packageAliases, goImport),
				g.fieldType(message.GetField()[1], mapTypes, packageAliases, goImport),
			}
			continue
		}

		nestedTypes, nestedEnumTypes := g.generateTypes(
			message.GetNestedType(),
			message.GetEnumType(),
			goImport,
			packagePrefix,
			mapTypes,
			packageAliases,
			sourceCodeInfo,
			append(msgSourceCodePath, int32(i), 3),
			append(msgSourceCodePath, int32(i), 4),
		)

		fields := make([]GoTypeField, 0, len(message.GetField()))
		for j, field := range message.GetField() {
			typ := g.fieldType(field, mapTypes, packageAliases, goImport)
			if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED && !strings.HasPrefix(typ, "map[") {
				typ = "[]" + typ
			}
			fields = append(fields, GoTypeField{
				Name:      common.SnakeCaseToPascalCase(field.GetName()),
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

func (g *Generator) fieldType(
	field *descriptorpb.FieldDescriptorProto,
	mapTypes map[string][2]string,
	packageAliases map[string]string,
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
		stderr.Failf("groups are not supported")
		return ""
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		if kv, ok := mapTypes[field.GetTypeName()]; ok {
			return "map[" + kv[0] + "]" + kv[1]
		}
		return g.fieldTypeMessageOrEnum(field, packageAliases, goImport)
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return "[]byte"
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		return g.fieldTypeMessageOrEnum(field, packageAliases, goImport)
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

func (g *Generator) fieldTypeMessageOrEnum(
	field *descriptorpb.FieldDescriptorProto,
	packageAliases map[string]string,
	goImport GoImport,
) string {
	goType, ok := g.definedTypes[field.GetTypeName()]
	if !ok {
		stderr.Failf("unknown type %s", field.GetTypeName())
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
		packageAlias = g.goImportPathToPackageAlias[importPath]

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
