package t6

import "testing"

func BenchmarkFor(b *testing.B) {
	var q1, q2 int
	for i := 0; i < b.N; i++ {
		for j := 0; j < 2; j++ {
			switch j {
			case 0:
				q1++
			case 1:
				q2++
			}
		}
	}
}

func BenchmarkGoto(b *testing.B) {
	var q1, q2 int
	for i := 0; i < b.N; i++ {
		j := 0
	l:
		j++
		switch j {
		case 1:
			q1++
			goto l
		case 2:
			q2++
			goto l
		}
	}
}
