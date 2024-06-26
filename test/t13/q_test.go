package t13

import (
	"fmt"
	"sync"
	"testing"
)

type Q struct {
	Int int
}
type ListQ []*Q

func Benchmark1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list := get()
		list = append(list, nil)
		put(list)
	}
}

func Benchmark2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list := *getP()
		list = append(list, nil)
		putP(&list)
	}
}

func Benchmark3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list := getP()
		*list = append(*list, nil)
		putP(list)
	}
}

func Benchmark4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pt := getP()
		list := *pt
		list = append(list, nil)
		*pt = list
		putP(pt)
	}
}

var pool = sync.Pool{New: func() any { return make(ListQ, 0) }}

func get() ListQ {
	return pool.Get().(ListQ)
}
func put(list ListQ) {
	list = list[:0]
	pool.Put(list)
}

var poolP = sync.Pool{New: func() any { list := make(ListQ, 0); return &list }}

func getP() *ListQ {
	return poolP.Get().(*ListQ)
}
func putP(list *ListQ) {
	*list = (*list)[:0]
	poolP.Put(list)
}

type U struct {
	Int int
}

func Test22(t *testing.T) {
	q := make([]*U, 0)
	q = append(q, nil, &U{1}, &U{2}, nil)
	fmt.Println(q)
	q = q[:0]
	fmt.Println(q)
	q = q[:3]
	fmt.Println(q)
	clear(q)
	q = q[:0]
	fmt.Println(q)
	q = q[:3]
	fmt.Println(q)
}
