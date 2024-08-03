package main

import (
	"math/big"
	"testing"
)

const count = 10

var bigInt, _ = new(big.Int).SetString("565432413251654654654625434542454264546235446254462544625344", 10)
var bigFloat, _ = new(big.Float).SetString("56543241325165465465462543454245426454.6235446254462544625344")

func BenchmarkSliceDynamic(b *testing.B) {
	q := make([]int32, 0, count)
	for i := 0; i < b.N; i++ {
		q = q[:0]
		for j := int32(0); j < count; j++ {
			q = append(q, j)
		}
	}
}

func BenchmarkSliceFixed(b *testing.B) {
	q := make([]int32, count)
	for i := 0; i < b.N; i++ {
		for j := int32(0); j < count; j++ {
			q[j] = j
		}
	}
}

func BenchmarkArray(b *testing.B) {
	q := [count]int32{}
	for i := 0; i < b.N; i++ {
		for j := int32(0); j < count; j++ {
			q[j] = j
		}
	}
}

func BenchmarkBigIntString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = bigInt.String()
	}
}

func BenchmarkBigIntBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = bigInt.Bytes()
	}
}

func BenchmarkBigIntGob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = bigInt.GobEncode()
	}
}

func BenchmarkBigFloatString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = bigFloat.String()
	}
}

func BenchmarkBigFloatGob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = bigFloat.GobEncode()
	}
}
