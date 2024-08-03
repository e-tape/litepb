package t9

import "testing"

type Q struct {
	int32 int32
	int64 int64
}

var qList = []*Q{
	{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
}

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := make(map[int]*Q)
		for j := 0; j < 100; j++ {
			q[j%10] = qList[j%10]
		}
	}
}

func BenchmarkMapCap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := make(map[int]*Q, 20)
		for j := 0; j < 100; j++ {
			q[j%10] = qList[j%10]
		}
	}
}
