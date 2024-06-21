package t8

import "testing"

type Q struct {
	Arr []int
}

var q = &Q{
	Arr: make([]int, 0),
}
var count int

func BenchmarkFree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, el := range q.Arr {
			count += el
		}
	}
}

func BenchmarkIf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if len(q.Arr) > 0 {
			for _, el := range q.Arr {
				count += el
			}
		}
	}
}

func BenchmarkEmpty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if q != nil && len(q.Arr) > 0 {
			for _, el := range q.Arr {
				count += el
			}
		}
	}
}
