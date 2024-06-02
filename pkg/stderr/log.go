//go:build !debug

package stderr

import (
	"fmt"
	"os"
)

// Logf noop
func Logf(_ string, _ ...any) {}

// Failf logs into stderr and exits with code 1
func Failf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
