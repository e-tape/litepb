package t10

import (
	"sync"
	"testing"
)

func TestQ(t *testing.T) {
	m := &sync.Map{}
	getData[uint](m)
	getData[int](m)
	getData[uint32](m)
	getData[int32](m)
	getData[float32](m)
	getData[float64](m)
	getData[string](m)
	getData[byte](m)
	getData[[3]byte](m)
	getData[[4]byte](m)
}

func getData[K comparable](m *sync.Map) {
	var key K
	_, loaded := m.LoadOrStore(key, key)
	if loaded {
		panic("loaded")
	}
}

var mapV1 sync.Map
var mapV2 sync.Map

type Q struct {
	Int  int
	Int1 int
	Int2 int
	Int3 int
	Int4 int
	Int5 int
	Int6 int
	Int7 int
}

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		poolRaw, _ := mapV1.LoadOrStore(i%5, &sync.Pool{New: func() any { return make(map[int]*Q) }})
		pool := poolRaw.(*sync.Pool)
		qm := pool.Get().(map[int]*Q)
		for j := 0; j < 300; j++ {
			qm[j%3] = &Q{
				Int:  i,
				Int2: j,
				Int7: i,
			}
		}
		clear(qm)
		pool.Put(qm)
	}
}

func BenchmarkMapV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		poolRaw, ok := mapV2.Load(i % 5)
		if !ok {
			poolRaw, _ = mapV2.LoadOrStore(i%5, &sync.Pool{})
		}
		pool := poolRaw.(*sync.Pool)
		qm, ok := pool.Get().(map[int]*Q)
		if !ok {
			qm = make(map[int]*Q, 100)
		}
		for j := 0; j < 300; j++ {
			qm[j%3] = &Q{
				Int:  i,
				Int2: j,
				Int7: i,
			}
		}
		clear(qm)
		pool.Put(qm)
	}
}

func BenchmarkMake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		qm := make(map[int]*Q)
		for j := 0; j < 300; j++ {
			qm[j%3] = &Q{
				Int:  i,
				Int2: j,
				Int7: i,
			}
		}
	}
}

func BenchmarkMakeCap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		qm := make(map[int]*Q, 100)
		for j := 0; j < 300; j++ {
			qm[j%3] = &Q{
				Int:  i,
				Int2: j,
				Int7: i,
			}
		}
	}
}
