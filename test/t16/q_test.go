package main

import (
	"testing"
)

func BenchmarkSwitch(b *testing.B) {
	q := make(map[int]int)
	for i := 0; i < b.N; i++ {
		u := i % 30
		switch u {
		case 3:
			q[1]++
		case 4:
			q[1]++
		case 22:
			q[1]++
		case 17:
			q[1]++
		case 8:
			q[2]++
		case 10:
			q[2]++
		case 1:
			q[3]++
		case 5:
			q[3]++
		case 18:
			q[3]++
		default:
			q[5]++
		}
	}
}

func BenchmarkSwitch2(b *testing.B) {
	q := make(map[int]int)
	for i := 0; i < b.N; i++ {
		u := i % 30
		switch u {
		case 3, 4, 22, 17:
			switch u {
			case 3:
				q[1]++
			case 4:
				q[1]++
			case 22:
				q[1]++
			case 17:
				q[1]++
			}
		case 8, 10:
			switch u {
			case 8:
				q[2]++
			case 10:
				q[2]++
			}
		case 1, 5:
			switch u {
			case 1:
				q[3]++
			case 5:
				q[3]++
			}
		case 18:
			q[3]++
		default:
			q[5]++
		}
	}
}

func BenchmarkDecodeCurrent(b *testing.B) {
	q := &Q{}
	for i := 0; i < b.N; i++ {
		q.DecodeCurrent(data, i%2)
	}
}

func BenchmarkDecodeCurrentModify(b *testing.B) {
	q := &Q{}
	for i := 0; i < b.N; i++ {
		q.DecodeCurrentModify(data, i%2)
	}
}

func BenchmarkDecodeModify(b *testing.B) {
	q := &Q{}
	for i := 0; i < b.N; i++ {
		q.DecodeModify(data, i%2)
		if i%2 == 0 {
			if q.Var1 != 8640662552 {
				panic("eq")
			}
		} else {
			if q.Var2 != 8640662552 {
				panic("eq")
			}
		}
	}
}

func BenchmarkDecodeModify2(b *testing.B) {
	q := &Q{}
	for i := 0; i < b.N; i++ {
		q.DecodeModify2(data, i%2)
		if i%2 == 0 {
			if q.Var1 != 8640662552 {
				panic("eq")
			}
		} else {
			if q.Var2 != 8640662552 {
				panic("eq")
			}
		}
	}
}

func BenchmarkDecodeModify3(b *testing.B) {
	q := &Q{}
	for i := 0; i < b.N; i++ {
		q.DecodeModify3(data, i%2)
		if i%2 == 0 {
			if q.Var1 != 8640662552 {
				panic("eq")
			}
		} else {
			if q.Var2 != 8640662552 {
				panic("eq")
			}
		}
	}
}
