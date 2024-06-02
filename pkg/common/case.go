package common

import (
	"bytes"
	"slices"
	"strings"
)

// SnakeCaseToPascalCase converts snake_case to PascalCase
func SnakeCaseToPascalCase(text string) string {
	data := []byte(strings.ToLower(text))
	if isLowerCaseLetter(data[0]) {
		data[0] -= 32
	}

	i := bytes.IndexByte(data, '_')
	for ; i >= 0; i = bytes.IndexByte(data, '_') {
		if isLowerCaseLetter(data[i+1]) {
			data[i+1] -= 32
		}
		data = slices.Delete(data, i, i+1)
	}

	return string(data)
}

func isLowerCaseLetter(c byte) bool {
	return c >= 'a' && c <= 'z'
}
