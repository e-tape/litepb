package t7

import (
	"encoding/binary"
	"testing"
)

var q = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
var index = 6
var r uint64

func BenchmarkNative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if len(q) > index+7 {
			r = binary.LittleEndian.Uint64(q[index:])
		}
	}
}
func BenchmarkNativeLen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if len(q) > index+7 {
			r = binary.LittleEndian.Uint64(q[index : index+8])
		}
	}
}

func BenchmarkSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if len(q) > index+7 {
			bb := q[index:]
			r = uint64(bb[0]) | uint64(bb[1])<<8 | uint64(bb[2])<<16 | uint64(bb[3])<<24 |
				uint64(bb[4])<<32 | uint64(bb[5])<<40 | uint64(bb[6])<<48 | uint64(bb[7])<<56
		}
	}
}

func BenchmarkSliceLen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if len(q) > index+7 {
			bb := q[index : index+8]
			r = uint64(bb[0]) | uint64(bb[1])<<8 | uint64(bb[2])<<16 | uint64(bb[3])<<24 |
				uint64(bb[4])<<32 | uint64(bb[5])<<40 | uint64(bb[6])<<48 | uint64(bb[7])<<56
		}
	}
}

func BenchmarkIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if len(q) > index+7 {
			_ = q[index+7]
			r = uint64(q[index]) | uint64(q[index+1])<<8 | uint64(q[index+2])<<16 | uint64(q[index+3])<<24 |
				uint64(q[index+4])<<32 | uint64(q[index+5])<<40 | uint64(q[index+6])<<48 | uint64(q[index+7])<<56
		}
	}
}
