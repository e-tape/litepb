package generator

import (
	"slices"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

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
