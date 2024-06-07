package main

import (
	_ "embed"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/e-tape/litepb/pkg/generator"
	"github.com/e-tape/litepb/pkg/stderr"
)

func main() {
	if err := run(); err != nil {
		stderr.Failf("%s: %s", filepath.Base(os.Args[0]), err)
	}
}

func run() error {
	in, err := os.ReadFile(`bin/1.bin`)
	if err != nil {
		return err
	}

	request := &pluginpb.CodeGeneratorRequest{}
	if err = proto.Unmarshal(in, request); err != nil {
		return err
	}

	stderr.Logf("COMPILER: %s", request.GetCompilerVersion())
	stderr.Logf("FILES TO GENERATE: %s", strings.Join(request.GetFileToGenerate(), ", "))

	start := time.Now()
	response := generator.NewGenerator(request).Generate()
	stderr.Logf("GENERATED IN: %s", time.Since(start))

	for _, f := range response.File {
		err = os.MkdirAll(`test/`+path.Dir(f.GetName()), 0766)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(`test/`+f.GetName(), []byte(f.GetContent()), 0666)
		if err != nil {
			panic(err)
		}
	}
	return nil

	generator.GoFmt(response)

	out, err := proto.Marshal(response)
	if err != nil {
		return err
	}

	if _, err = os.Stdout.Write(out); err != nil {
		return err
	}

	return nil
}
