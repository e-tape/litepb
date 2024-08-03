package zerocast

import (
	"strings"
	"testing"
	"unsafe"
)

const size = 2000

var s = strings.Repeat("текст", size)
var bb = []byte(strings.Repeat("текст", size))
var s1 string
var bb1 []byte
var s1Init = string(bb)

//func BenchmarkStringToBytesStandard(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		bb1 = []byte(s)
//	}
//}

func BenchmarkBytesToStringStandard(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s1 = string(bb)
	}
	if string(bb) != s1Init {
		panic("eq")
	}
}

//func BenchmarkStringToBytes(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		bb1 = StringToBytes(s)
//	}
//}

func BenchmarkBytesToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//s1 = unsafe.String(&bb[10], len(bb)-10)
		s1 = BytesToString(bb)
		//if unsafe.String(&bb[10], len(bb)-10) != s1Init {
	}
	if s1 != s1Init {
		panic("eq")
	}
}

//	func StringToBytes(s string) []byte {
//		p := unsafe.StringData(s)
//		b := unsafe.Slice(p, len(s))
//		return b
//	}
func BytesToString(b []byte) string {
	return unsafe.String(&b[0], len(b))
}
