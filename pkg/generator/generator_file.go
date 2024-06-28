package generator

import (
	"path"
	"slices"
	"strings"

	"github.com/e-tape/litepb/pkg/stderr"
	litepb "github.com/e-tape/litepb/proto"
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
	name := strings.ToLower(strings.TrimSuffix(
		path.Base(protoFile.GetName()),
		path.Ext(protoFile.GetName()),
	)) + ".lpb.go"
	namePath := path.Dir(packagePath)
	if a.cfg.SourceRelative {
		namePath = path.Dir(protoFile.GetName())
	}
	return &generatorFile{
		Generator: a,
		proto: &litepb.File{
			Package: &litepb.Package{
				Dependency: &litepb.Dependency{
					Path:  packagePath,
					Alias: alias,
				},
				Name: packageName,
			},
			Source:  protoFile.GetName(),
			Name:    path.Join(namePath, name),
			Options: protoFile.GetOptions().ProtoReflect().GetUnknown(),
		},
		sourceCodeInfo: protoFile.GetSourceCodeInfo(),
	}
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
		a.allTypes[a.generatePackage(packages, enum.GetName())] = &litepb.Message_Field_Type_Reflect{
			Name:       a.generateTypeName(names, generateName(enum.GetName())),
			Dependency: a.proto.Package.Dependency,
		}
	}

	for _, message := range messages {
		a.allTypes[a.generatePackage(packages, message.GetName())] = &litepb.Message_Field_Type_Reflect{
			Name:       a.generateTypeName(names, generateName(message.GetName())),
			Dependency: a.proto.Package.Dependency,
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
			a.mapTypes[a.generatePackage(packages, message.GetName())] = &litepb.Message_Field_Type_Map{
				Key: &litepb.Message_Field_Type{
					InProto: litepb.Message_Field_Type_Proto(message.GetField()[0].GetType()),
					Reflect: a.generateReflect(message.GetField()[0], false),
				},
				Value: &litepb.Message_Field_Type{
					InProto: litepb.Message_Field_Type_Proto(message.GetField()[1].GetType()),
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
) []*litepb.Enum {
	result := make([]*litepb.Enum, 0, len(enums))
	for enumIndex, enum := range enums {
		values := make([]*litepb.Enum_Value, 0, len(enum.GetValue()))
		for valueIndex, value := range enum.GetValue() {
			values = append(values, &litepb.Enum_Value{
				Name:     value.GetName(),
				Comments: a.findEnumValueComments(enumSourceCodePath, enumIndex, valueIndex),
				Number:   value.GetNumber(),
				Options:  value.GetOptions().ProtoReflect().GetUnknown(),
			})
		}
		result = append(result, &litepb.Enum{
			Name:         a.generateTypeName(names, generateName(enum.GetName())),
			Comments:     a.findEnumComments(enumSourceCodePath, enumIndex),
			ValuesPrefix: generateName(enum.GetName()),
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
) {
	for messageIndex, message := range messages {
		if message.GetOptions().GetMapEntry() {
			continue
		}

		msg := &litepb.Message{
			Name:       a.generateTypeName(names, generateName(message.GetName())),
			Comments:   a.findMessageComments(messageSourceCodePath, messageIndex),
			Properties: make([]*litepb.Message_Property, 0, len(message.GetField())),
			Options:    message.GetOptions().ProtoReflect().GetUnknown(),
			//MemPoolMessage:     len(message.GetField()) > 0, // TODO rathil add option to disable
			//MemPoolList: true,                        // TODO rathil from options
			//MemPoolMap:  true,                        // TODO rathil from options
		}
		msg.MemPoolMessage = a.cfg.MemPoolMessageAll == litepb.Activity_Active
		msg.MemPoolList = a.cfg.MemPoolListAll == litepb.Activity_Active
		msg.MemPoolMap = a.cfg.MemPoolMapAll == litepb.Activity_Active
		// TODO rathil get option to active/inactive all MemPool
		a.messages = append(a.messages, msg)

		oneOfs := make([]*litepb.Message_OneOf, 0, len(message.GetOneofDecl()))
		for fieldIndex, field := range message.GetField() {
			mapField, mapOk := a.mapTypes[field.GetTypeName()]
			msgField := &litepb.Message_Field{
				Number:   field.GetNumber(),
				Name:     a.generateFieldName(msg, field.GetName()),
				Comments: a.findMessageFieldComments(messageSourceCodePath, messageIndex, fieldIndex),
				Type: &litepb.Message_Field_Type{
					InProto:  litepb.Message_Field_Type_Proto(field.GetType()),
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
			property := &litepb.Message_Property{}
			if field.OneofIndex != nil {
				if len(oneOfs) <= int(field.GetOneofIndex()) {
					oneofDecl := message.GetOneofDecl()[field.GetOneofIndex()]
					oneof := &litepb.Message_OneOf{
						Name:     a.generateFieldName(msg, oneofDecl.GetName()),
						Comments: "", // TODO rathil implement
						Options:  oneofDecl.GetOptions().ProtoReflect().GetUnknown(),
						Fields:   make([]*litepb.Message_Field, 0, 1),
						Tags: map[string]string{
							"json": oneofDecl.GetName(),
						},
					}
					switch {
					case a.cfg.MemPoolOneofAll == litepb.Activity_Inactive:
						oneof.MemPool = false
					case a.cfg.MemPoolOneofAll == litepb.Activity_Active:
						oneof.MemPool = true
					}
					// TODO rathil get option to active/inactive MemPoolMessage
					oneOfs = append(oneOfs, oneof)
				}
				oneOfs[field.GetOneofIndex()].Fields = append(oneOfs[field.GetOneofIndex()].Fields, msgField)
				if len(oneOfs[field.GetOneofIndex()].Fields) > 1 {
					continue
				}
				property.Type = &litepb.Message_Property_Oneof{
					Oneof: oneOfs[field.GetOneofIndex()],
				}
			} else {
				property.Type = &litepb.Message_Property_Field{
					Field: msgField,
				}
			}
			msg.Properties = append(msg.Properties, property)
		}
		a.generateMessages(
			message.GetNestedType(),
			append(packages, message.GetName()),
			append(messageSourceCodePath, int32(messageIndex), 3),
			append(names, message.GetName()),
		)
	}
}

func (a *generatorFile) generateFieldName(msg *litepb.Message, name string) string {
	pName := generateName(name)
	switch pName {
	case "String":
		pName += "_"
	}
	for _, p := range msg.Properties {
		switch pt := p.Type.(type) {
		case *litepb.Message_Property_Field:
			if pt.Field.GetName() == pName {
				pName += "_"
			}
		case *litepb.Message_Property_Oneof:
			if pt.Oneof.GetName() == pName {
				pName += "_"
			}
		}
	}
	return pName
}

func (a *generatorFile) generateReflect(
	field *descriptorpb.FieldDescriptorProto,
	mapOk bool,
) *litepb.Message_Field_Type_Reflect {
	if mapOk {
		return nil
	}
	switch field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		return &litepb.Message_Field_Type_Reflect{Name: "float64"}
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		return &litepb.Message_Field_Type_Reflect{Name: "float32"}
	case descriptorpb.FieldDescriptorProto_TYPE_INT64,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64,
		descriptorpb.FieldDescriptorProto_TYPE_SINT64:
		return &litepb.Message_Field_Type_Reflect{Name: "int64"}
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		return &litepb.Message_Field_Type_Reflect{Name: "uint64"}
	case descriptorpb.FieldDescriptorProto_TYPE_INT32,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32:
		return &litepb.Message_Field_Type_Reflect{Name: "int32"}
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		return &litepb.Message_Field_Type_Reflect{Name: "uint32"}
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return &litepb.Message_Field_Type_Reflect{Name: "bool"}
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return &litepb.Message_Field_Type_Reflect{Name: "string"}
	case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
		stderr.Failf("groups are not supported")
		return nil
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
		descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		return a.allTypes[field.GetTypeName()] //.reflect(a.proto.Package.Dependency.Alias)
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return &litepb.Message_Field_Type_Reflect{Name: "[]byte"}
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
