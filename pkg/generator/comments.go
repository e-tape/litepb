package generator

import (
	"slices"
	"strings"
)

func (a *generatorFile) findMessageComments(
	sourceCodePath []int32,
	messageIndex int,
) string {
	return a.findComments(append(sourceCodePath, int32(messageIndex)))
}

func (a *generatorFile) findMessageFieldComments(
	sourceCodePath []int32,
	messageIndex int,
	fieldIndex int,
) string {
	// TODO rathil move 2 to const
	return a.findComments(append(sourceCodePath, int32(messageIndex), 2, int32(fieldIndex)))
}

func (a *generatorFile) findEnumComments(
	sourceCodePath []int32,
	enumIndex int,
) string {
	return a.findComments(append(sourceCodePath, int32(enumIndex)))
}

func (a *generatorFile) findEnumValueComments(
	sourceCodePath []int32,
	enumIndex int,
	valueIndex int,
) string {
	// TODO rathil move 2 to const
	return a.findComments(append(sourceCodePath, int32(enumIndex), 2, int32(valueIndex)))
}

func (a *generatorFile) findComments(ps []int32) string {
	for _, loc := range a.sourceCodeInfo.GetLocation() {
		if slices.Equal(loc.GetPath(), ps) {
			return strings.Trim(loc.GetLeadingComments()+loc.GetTrailingComments(), "\n")
		}
	}
	return ""
}
