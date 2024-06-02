//go:build debug

package stderr

import (
	"fmt"
	"os"
	"runtime"
)

// Logf logs into stderr
func Logf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
}

// Failf logs into stderr and exits with code 1
func Failf(format string, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf("%s:%d: ", file, line)+format+"\n", args...)
	os.Exit(1)
}
