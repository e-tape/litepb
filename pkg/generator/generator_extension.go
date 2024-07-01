package generator

import (
	"path"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

func (a *Generator) generateMessagesFromExtensions(protoFile *descriptorpb.FileDescriptorProto) {
	if len(protoFile.Extension) == 0 {
		return
	}
	msgNamePrefix := generateName(strings.TrimSuffix(
		path.Base(protoFile.GetName()),
		path.Ext(protoFile.GetName()),
	))

	msgMap := make(map[string]*descriptorpb.DescriptorProto)
	for _, ext := range protoFile.Extension {
		msgNameParts := strings.Split(ext.GetExtendee(), ".")
		msgName := msgNamePrefix + generateName(msgNameParts[len(msgNameParts)-1])
		if _, ok := msgMap[msgName]; !ok {
			msgMap[msgName] = &descriptorpb.DescriptorProto{
				Name: &msgName,
			}
		}
		msgMap[msgName].Field = append(msgMap[msgName].Field, ext)
	}
	for _, msg := range msgMap {
		protoFile.MessageType = append(protoFile.MessageType, msg)
	}
}
