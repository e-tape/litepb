package main

import (
	"testing"
	"unsafe"
)

func BenchmarkCast1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = data[:5]
	}
}
func BenchmarkCast2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = *(*Bytes)(unsafe.Pointer(&data[0]))
	}
}
