package t4

import "testing"

var found int64

func BenchmarkIf1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if i%3 == 0 {
			found++
		}
	}
}
func BenchmarkSwitch1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		switch i % 3 {
		case 0:
			found++
		}
	}
}

func BenchmarkIf2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if i%3 == 0 {
			found++
		} else if i%3 == 1 {
			found++
		}
	}
}

func BenchmarkSwitch2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		switch i % 3 {
		case 0:
			found++
		case 1:
			found++
		}
	}
}

var q int

func BenchmarkD(b *testing.B) {
	var v int
	for i := 0; i < b.N; i++ {
		v += i
	}
	v = v >> 3
	q = v
}

func BenchmarkL(b *testing.B) {
	var v int
	for i := 0; i < b.N; i++ {
		v += i
	}
	q = v >> 3
}
