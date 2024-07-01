package main

import (
	_ "embed"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/e-tape/litepb/config"
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
	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	//_ = os.WriteFile(`litepb.git/bin/1.bin`, in, 0666)

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
