package main

import (
	_ "embed"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/e-tape/litepb/config"
	"github.com/e-tape/litepb/pkg/generator"
	"github.com/e-tape/litepb/pkg/stderr"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	if err := run(); err != nil {
		stderr.Failf("%s: %s", filepath.Base(os.Args[0]), err)
	}
}

func run() error {
	in, err := os.ReadFile(`/tmp/1.bin`)
	if err != nil {
		return err
	}

	request := &pluginpb.CodeGeneratorRequest{}
	if err = proto.Unmarshal(in, request); err != nil {
		return err
	}
	cfg, err := config.Parse(request.Parameter)
	if err != nil {
		return err
	}

	stderr.Logf("COMPILER: %s", request.GetCompilerVersion())
	stderr.Logf("FILES TO GENERATE: %s", strings.Join(request.GetFileToGenerate(), ", "))

	start := time.Now()
	response := generator.NewGenerator(cfg, request).Generate()
	stderr.Logf("GENERATED IN: %s", time.Since(start))

	//generator.GoFmt(response)

	for _, f := range response.File {
		p := path.Join("test/bench/proto/qtest/", f.GetName())
		err = os.MkdirAll(path.Dir(p), 0766)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(p, []byte(f.GetContent()), 0666)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
