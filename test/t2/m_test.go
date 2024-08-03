package main

import (
	"fmt"
	"testing"
	"unsafe"
)

func BenchmarkL(b *testing.B) {
	u := 7
	for i := 0; i < b.N; i++ {
		_ = (u) << 3
	}
}

func BenchmarkF(b *testing.B) {
	u := 7
	for i := 0; i < b.N; i++ {
		_ = (u & i) << 3
	}
}

func BenchmarkIf(b *testing.B) {
	u := 7
	c := 0
	for i := 0; i < b.N; i++ {
		if u == i {
			continue
		} else if u == i+1 {
			continue
		}
		c++
	}
}

func BenchmarkSwitch(b *testing.B) {
	u := 7
	c := 0
	for i := 0; i < b.N; i++ {
		switch u {
		case i:
			continue
		case i + 1:
			continue
		default:
			c++
		}
	}
}

type Q struct {
	Data uint64
}

func (a *Q) direct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		a.Data++
	}
}
func (a *Q) local(b *testing.B) {
	data := uint64(0)
	for i := 0; i < b.N; i++ {
		data++
	}
	a.Data = data
}

func BenchmarkDirect(b *testing.B) {
	(&Q{}).direct(b)
}

func BenchmarkLocal(b *testing.B) {
	(&Q{}).local(b)
}

var strData = []byte("Моё текст на русском языке 23 вот так")
var strLen = len(strData) - 3
var strString = string(strData[4:strLen])

func TestString(t *testing.T) {
	fmt.Println(strData)
	fmt.Println(string(strData))
	fmt.Println(unsafe.String(&strData[0], len(strData)-4))
}

func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if strString != string(strData[4:strLen]) {
			panic("eq")
		}
	}
}

var q string

func BenchmarkUnsafeString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q = unsafe.String(&strData[4], strLen-4)
		if strString != q {
			panic("eq")
		}
	}
}
