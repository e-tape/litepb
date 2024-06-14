package main

import (
	"math/bits"
	"strconv"
	"testing"
)

var mathBits = []uint64{
	1, 10, 100, 1_000, 10_000, 100_000, 1_000_000, 10_000_000, 100_000_000, 1_000_000_000, 10_000_000_000, 100_000_000_000, 1_000_000_000_000,
}

func BenchmarkMathBits(b *testing.B) {
	bLen := 0
	for _, bm := range mathBits {
		b.Run(strconv.FormatUint(bm, 10), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bLen = (bits.Len64(bm+uint64(i)|1) + 6) / 7
			}
		})
	}
	_ = bLen
}

func BenchmarkMathBitsIf(b *testing.B) {
	bLen := 0
	for _, bm := range mathBits {
		b.Run(strconv.FormatUint(bm, 10), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bLen = mathBitsIf(bm + uint64(i))
			}
		})
	}
	_ = bLen
}

func mathBitsIf(value uint64) int {
	if value <= 127 {
		return 1
	}
	if value <= 16383 {
		return 2
	}
	if value <= 2097151 {
		return 3
	}
	if value <= 268435455 {
		return 4
	}
	if value <= 34359738367 {
		return 5
	}
	if value <= 4398046511103 {
		return 6
	}
	if value <= 562949953421311 {
		return 7
	}
	if value <= 72057594037927935 {
		return 8
	}
	if value <= 9223372036854775807 {
		return 9
	}
	return 10
}

func BenchmarkMathBitsCase(b *testing.B) {
	bLen := 0
	for _, bm := range mathBits {
		b.Run(strconv.FormatUint(bm, 10), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bLen = mathBitsCase(bm + uint64(i))
			}
		})
	}
	_ = bLen
}

func mathBitsCase(value uint64) int {
	switch {
	case value <= 127:
		return 1
	case value <= 16383:
		return 2
	case value <= 2097151:
		return 3
	case value <= 268435455:
		return 4
	case value <= 34359738367:
		return 5
	case value <= 4398046511103:
		return 6
	case value <= 562949953421311:
		return 7
	case value <= 72057594037927935:
		return 8
	case value <= 9223372036854775807:
		return 9
	default:
		return 10
	}
}
