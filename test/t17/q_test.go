package t17

import "testing"

func BenchmarkFor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		For()
	}
}

func BenchmarkGoto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Goto()
	}
}

func BenchmarkGotoModify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GotoModify()
	}
}
