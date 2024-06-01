package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"runtime"
	"slices"
	"strings"

	"google.golang.org/protobuf/types/pluginpb"
)

func logf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func failf(format string, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	logf(fmt.Sprintf("%s:%d: ", file, line)+format, args...)
	os.Exit(1)
}

func ptr[T any](value T) *T {
	return &value
}

func goFmt(resp *pluginpb.CodeGeneratorResponse) error {
	for i := 0; i < len(resp.File); i++ {
		formatted, err := format.Source([]byte(resp.File[i].GetContent()))
		if err != nil {
			return fmt.Errorf("go fmt: %w", err)
		}
		resp.File[i].Content = ptr(string(formatted))
	}
	return nil
}

func snakeCaseToCamelCase(text string) string {
	data := []byte(strings.ToLower(text))
	data[0] -= 32

	i := bytes.IndexByte(data, '_')
	for ; i >= 0; i = bytes.IndexByte(data, '_') {
		data[i+1] -= 32
		data = slices.Delete(data, i, i+1)
	}

	return string(data)
}
