package t5

import (
	"bytes"
	"testing"
)

const size = 200

var bb = bytes.Repeat([]byte{1, 2, 3, 4, 5}, size)

func BenchmarkCap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r(bb[20:200:200])
	}
}

func BenchmarkLen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r(bb[20:200])
	}
}

func r(data []byte) {
	_ = len(data)
}
