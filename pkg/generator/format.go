package generator

import (
	"go/format"

	"google.golang.org/protobuf/types/pluginpb"

	"github.com/e-tape/litepb/pkg/common"
	"github.com/e-tape/litepb/pkg/stderr"
)

// GoFmt formats all generated files
func GoFmt(resp *pluginpb.CodeGeneratorResponse) {
	for i := 0; i < len(resp.File); i++ {
		formatted, err := format.Source([]byte(resp.File[i].GetContent()))
		if err != nil {
			stderr.Failf("go fmt: %s: %s", resp.File[i].GetName(), err)
		}
		resp.File[i].Content = common.Ptr(string(formatted))
	}
}
