package common

import (
	"bytes"
	"slices"
	"strings"
)

// SnakeCaseToPascalCase converts snake_case to PascalCase
func SnakeCaseToPascalCase(text string) string {
	data := []byte(strings.ToLower(text))
	data[0] -= 32

	i := bytes.IndexByte(data, '_')
	for ; i >= 0; i = bytes.IndexByte(data, '_') {
		if data[i+1] >= 'a' && data[i+1] <= 'z' || data[i+1] >= 'A' && data[i+1] <= 'Z' {
			data[i+1] -= 32
		}
		data = slices.Delete(data, i, i+1)
	}

	return string(data)
}
