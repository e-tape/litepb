package t12

import (
	"fmt"
	"sync"
	"testing"
)

func TestQ(t *testing.T) {
	pool := sync.Pool{}
	q := make([]int, 2, 10)
	fmt.Println(len(q), cap(q))
	pool.Put(q)
	q = pool.Get().([]int)
	fmt.Println(len(q), cap(q))
}

func BenchmarkQ(b *testing.B) {
	//countGet := 0
	//countNot := 0
	for i := 0; i < b.N; i++ {
		pool := sync.Pool{}
		q := make([]int, 2, 10)
		pool.Put(q)
		if _, ok := pool.Get().([]int); ok {
			//countGet++
			//} else {
			//	countNot++
		}
	}
	//fmt.Println("count get", countGet, "count not", countNot)
}

func BenchmarkParallelSliceNew(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q := make(ListQ, 0, 20)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
		}
	})
}

func BenchmarkParallelSlice(b *testing.B) {
	var pool = &Pool[*Q]{&sync.Pool{}}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q := pool.Get(20)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			q = append(q, nil)
			pool.Put(q)
		}
	})
}

func BenchmarkParallelSlicePointer(b *testing.B) {
	var pool = &PoolPointer[*Q]{&sync.Pool{}}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q := pool.Get(20)
			*q = append(*q, nil)
			*q = append(*q, nil)
			*q = append(*q, nil)
			*q = append(*q, nil)
			*q = append(*q, nil)
			*q = append(*q, nil)
			*q = append(*q, nil)
			*q = append(*q, nil)
			*q = append(*q, nil)
			*q = append(*q, nil)
			*q = append(*q, nil)
			*q = append(*q, nil)
			*q = append(*q, nil)
			pool.Put(q)
		}
	})
}

type PoolPointer[T any] struct {
	*sync.Pool
}
type Pool[T any] struct {
	*sync.Pool
}
type Q struct {
	Int int
}
type ListQ []*Q

func (a *PoolPointer[T]) Get(cap int) *[]T {
	if q, ok := a.Pool.Get().(*[]T); ok {
		return q
	}
	q := make([]T, 0, cap)
	return &q
}

func (a *PoolPointer[T]) Put(q *[]T) {
	clear(*q)
	*q = (*q)[:0]
	a.Pool.Put(q)
}

func (a *Pool[T]) Get(cap int) []T {
	if q, ok := a.Pool.Get().(*[]T); ok {
		return *q
	}
	return make([]T, 0, cap)
}

func (a *Pool[T]) Put(q []T) {
	clear(q)
	q = q[:0]
	a.Pool.Put(&q)
}
