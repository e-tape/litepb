package generator

import (
	"cmp"
	"path"
	"slices"
	"strings"

	"github.com/e-tape/litepb/pkg/common"
	"github.com/e-tape/litepb/pkg/plugin"
	"github.com/e-tape/litepb/pkg/stderr"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/descriptorpb"
)

func (a *Generator) newFile(protoFile *descriptorpb.FileDescriptorProto) *generatorFile {
	goPackage := protoFile.GetOptions().GetGoPackage()
	if goPackage == "" {
		goPackage = strings.ReplaceAll(protoFile.GetPackage(), ".", "/")
	}
	stderr.Logf("\tGO PACKAGE: %s", goPackage)

	packagePath, pathAlias, ok := strings.Cut(goPackage, ";")
	packageName := strings.ToLower(path.Base(packagePath))
	if ok {
		packageName = pathAlias
	}
	alias := a.aliasRegex.ReplaceAllString(packagePath, `_`)
	return &generatorFile{
		Generator: a,
		proto: &plugin.File{
			Package: &plugin.Package{
				Dependence: &plugin.Dependence{
					Path:  packagePath,
					Alias: alias,
				},
				Name: packageName,
			},
			Source: protoFile.GetName(),
			Name: strings.ToLower(strings.TrimSuffix(
				path.Base(protoFile.GetName()),
				path.Ext(protoFile.GetName()),
			)) + ".lpb.go",
			Options: protoFile.GetOptions().ProtoReflect().GetUnknown(),
		},
		sourceCodeInfo: protoFile.GetSourceCodeInfo(),
	}
}

func (a *generatorFile) generateImports(dependencies []string) []*plugin.Dependence {
	result := make([]*plugin.Dependence, 0, len(dependencies))
	for _, dependency := range dependencies {
		if a.allFiles[dependency].proto.GetPackage().GetDependence().GetAlias() == a.proto.GetPackage().GetDependence().GetAlias() {
			continue
		}
		result = append(result, a.allFiles[dependency].proto.GetPackage().GetDependence())
	}
	result = append(result, &plugin.Dependence{
		Path: "fmt",
	})
	for _, message := range a.proto.Messages {
		if message.GetWithMemPool() {
			result = append(result, &plugin.Dependence{
				Path: "sync",
			})
			break
		}
	}
	slices.SortFunc(result, func(a, b *plugin.Dependence) int {
		return cmp.Compare(a.Alias, b.Alias)
	})
	return slices.CompactFunc(result, func(a, b *plugin.Dependence) bool {
		if a.Alias == "" && b.Alias == "" {
			return a.Path == b.Path
		}
		return a.Alias == b.Alias
	})
}

func (a *generatorFile) generatePackage(packages []string, item string) string {
	return a.generateJoin(packages, item, ".")
}
func (a *generatorFile) generateTypeName(names []string, name string) string {
	return a.generateJoin(names, name, "_")
}
func (a *generatorFile) generateJoin(items []string, item string, sep string) string {
	return strings.Join(slices.Concat(slices.Concat(items, []string{item})), sep)
}

func (a *generatorFile) collectTypes(
	enums []*descriptorpb.EnumDescriptorProto,
	messages []*descriptorpb.DescriptorProto,
	packages []string,
	names []string,
) {
	for _, enum := range enums {
		a.allTypes[a.generatePackage(packages, enum.GetName())] = Type{
			Name:  a.generateTypeName(names, enum.GetName()),
			Alias: a.proto.Package.Dependence.Alias,
		}
	}

	for _, message := range messages {
		a.allTypes[a.generatePackage(packages, message.GetName())] = Type{
			Name:  a.generateTypeName(names, message.GetName()),
			Alias: a.proto.Package.Dependence.Alias,
		}
		a.collectTypes(
			message.GetEnumType(),
			message.GetNestedType(),
			append(packages, message.GetName()),
			append(names, message.GetName()),
		)
	}
}

func (a *generatorFile) collectMapTypes(
	messages []*descriptorpb.DescriptorProto,
	packages []string,
) {
	for _, message := range messages {
		if message.GetOptions().GetMapEntry() {
			a.mapTypes[a.generatePackage(packages, message.GetName())] = &plugin.Message_Field_Type_Map{
				Key: &plugin.Message_Field_Type{
					InProto: plugin.Message_Field_Type_Proto(message.GetField()[0].GetType()),
					Reflect: a.generateReflect(message.GetField()[0], false),
				},
				Value: &plugin.Message_Field_Type{
					InProto: plugin.Message_Field_Type_Proto(message.GetField()[1].GetType()),
					Reflect: a.generateReflect(message.GetField()[1], false),
				},
			}
			continue
		}
		a.collectMapTypes(
			message.GetNestedType(),
			append(packages, message.GetName()),
		)
	}
}

func (a *generatorFile) generateEnums(
	enums []*descriptorpb.EnumDescriptorProto,
	messages []*descriptorpb.DescriptorProto,
	packages []string,
	messageSourceCodePath []int32,
	enumSourceCodePath []int32,
	names []string,
) []*plugin.Enum {
	result := make([]*plugin.Enum, 0, len(enums))
	for enumIndex, enum := range enums {
		values := make([]*plugin.Enum_Value, 0, len(enum.GetValue()))
		for valueIndex, value := range enum.GetValue() {
			values = append(values, &plugin.Enum_Value{
				Name:     value.GetName(),
				Comments: a.findEnumValueComments(enumSourceCodePath, enumIndex, valueIndex),
				Number:   value.GetNumber(),
				Options:  value.GetOptions().ProtoReflect().GetUnknown(),
			})
		}
		result = append(result, &plugin.Enum{
			Name:         a.generateTypeName(names, enum.GetName()),
			Comments:     a.findEnumComments(enumSourceCodePath, enumIndex),
			ValuesPrefix: enum.GetName(),
			Values:       values,
			Options:      enum.GetOptions().ProtoReflect().GetUnknown(),
		})
	}

	for i, message := range messages {
		nestedEnums := a.generateEnums(
			message.GetEnumType(),
			message.GetNestedType(),
			append(packages, message.GetName()),
			append(messageSourceCodePath, int32(i), 3),
			append(messageSourceCodePath, int32(i), 4),
			append(names, message.GetName()),
		)
		result = append(result, nestedEnums...)
	}

	return result
}

func (a *generatorFile) generateMessages(
	messages []*descriptorpb.DescriptorProto,
	packages []string,
	messageSourceCodePath []int32,
	names []string,
) []*plugin.Message {
	result := make([]*plugin.Message, 0, len(messages))
	for messageIndex, message := range messages {
		if message.GetOptions().GetMapEntry() {
			continue
		}

		result = append(result, a.generateMessages(
			message.GetNestedType(),
			append(packages, message.GetName()),
			append(messageSourceCodePath, int32(messageIndex), 3),
			append(names, message.GetName()),
		)...)

		oneOfs := make([]*plugin.Message_OneOf, 0, len(message.GetOneofDecl()))
		for _, oneof := range message.GetOneofDecl() {
			oneOfs = append(oneOfs, &plugin.Message_OneOf{
				Name:     a.generateOneofName(oneof.GetName()),
				Comments: "", // TODO rathil implement
				Options:  oneof.GetOptions().ProtoReflect().GetUnknown(),
				Fields:   make([]*plugin.Message_Field, 0, 1),
				Tags: map[string]string{
					"json": oneof.GetName(),
				},
				WithMemPool: true,
			})
		}
		properties := make([]*plugin.Message_Property, 0, len(message.GetField()))
		for fieldIndex, field := range message.GetField() {
			mapField, mapOk := a.mapTypes[field.GetTypeName()]
			msgField := &plugin.Message_Field{
				Number:   field.GetNumber(),
				Name:     common.SnakeCaseToPascalCase(field.GetName()),
				Comments: a.findMessageFieldComments(messageSourceCodePath, messageIndex, fieldIndex),
				Type: &plugin.Message_Field_Type{
					InProto:  plugin.Message_Field_Type_Proto(field.GetType()),
					Reflect:  a.generateReflect(field, mapOk),
					Repeated: !mapOk && field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED,
					Map:      mapField,
				},
				ZeroValue: a.fieldZeroValue(field),
				Options:   field.GetOptions().ProtoReflect().GetUnknown(),
				Tags: map[string]string{
					"json": field.GetJsonName(),
				},
			}
			property := &plugin.Message_Property{}
			if field.OneofIndex != nil {
				oneOfs[field.GetOneofIndex()].Fields = append(oneOfs[field.GetOneofIndex()].Fields, msgField)
				if len(oneOfs[field.GetOneofIndex()].Fields) > 1 {
					continue
				}
				property.Type = &plugin.Message_Property_Oneof{
					Oneof: oneOfs[field.GetOneofIndex()],
				}
			} else {
				property.Type = &plugin.Message_Property_Field{
					Field: msgField,
				}
			}
			properties = append(properties, property)
		}

		result = append(result, &plugin.Message{
			Name:        a.generateTypeName(names, message.GetName()),
			Comments:    a.findMessageComments(messageSourceCodePath, messageIndex),
			Properties:  properties,
			Options:     message.GetOptions().ProtoReflect().GetUnknown(),
			WithMemPool: len(message.GetField()) > 0, // TODO rathil add option to disable
		})
	}
	return result
}

func (a *generatorFile) generateOneofName(name string) string {
	return strings.ReplaceAll(
		cases.Title(language.English).
			String(
				strings.ReplaceAll(name, "_", " "),
			),
		" ", "",
	)
}

func (a *generatorFile) generateReflect(
	field *descriptorpb.FieldDescriptorProto,
	mapOk bool,
) *plugin.Message_Field_Type_Reflect {
	if mapOk {
		return nil
	}
	switch field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		return &plugin.Message_Field_Type_Reflect{Name: "float64"}
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		return &plugin.Message_Field_Type_Reflect{Name: "float32"}
	case descriptorpb.FieldDescriptorProto_TYPE_INT64,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64,
		descriptorpb.FieldDescriptorProto_TYPE_SINT64:
		return &plugin.Message_Field_Type_Reflect{Name: "int64"}
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		return &plugin.Message_Field_Type_Reflect{Name: "uint64"}
	case descriptorpb.FieldDescriptorProto_TYPE_INT32,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32:
		return &plugin.Message_Field_Type_Reflect{Name: "int32"}
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		return &plugin.Message_Field_Type_Reflect{Name: "uint32"}
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return &plugin.Message_Field_Type_Reflect{Name: "bool"}
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return &plugin.Message_Field_Type_Reflect{Name: "string"}
	case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
		stderr.Failf("groups are not supported")
		return nil
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
		descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		return a.allTypes[field.GetTypeName()].reflect(a.proto.Package.Dependence.Alias)
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return &plugin.Message_Field_Type_Reflect{Name: "[]byte"}
	default:
		stderr.Failf("unknown type %d", field.GetType())
		return nil
	}
}

func (a *generatorFile) fieldZeroValue(field *descriptorpb.FieldDescriptorProto) string {
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
