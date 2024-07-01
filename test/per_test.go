package main

import (
	"fmt"
	"sync"
	"testing"
)

type Qu struct {
	Int int
	err error
}

func (a *Qu) SetInt(value int) *Qu {
	if a == nil {
		a = &Qu{}
	}
	a.Int = value
	return a
}

func (a *Qu) ppp() {
	*a = Qu{}
}

func TestPointer(t *testing.T) {
	var q *Qu
	q.SetInt(55).SetInt(6547)
	fmt.Println(q)
	q = &Qu{Int: 555}
	fmt.Printf("%p\n", q)
	q.ppp()
	fmt.Printf("%p\n", q)
	q.ppp()
	fmt.Printf("%p\n", q)
}

type Q struct {
	f1 int
	f2 string
	f3 float32
	f4 []int32
	f5 map[int]int
	f6 int64
	f7 *Q
	f8 bool
	f9 []byte
	f0 uint32
}
type IQ interface {
	get_f1() int
	get_f2() string
	get_f3() float32
	get_f4() []int32
	get_f5() map[int]int
	get_f6() int64
	get_f7() *Q
	get_f8() bool
	get_f9() []byte
	get_f0() uint32
}

func (a *Q) get_f1() int {
	if a == nil {
		return 0
	}
	return a.f1
}
func (a *Q) get_f2() string {
	if a == nil {
		return ""
	}
	return a.f2
}
func (a *Q) get_f3() float32 {
	if a == nil {
		return 0
	}
	return a.f3
}
func (a *Q) get_f4() []int32 {
	if a == nil {
		return nil
	}
	return a.f4
}
func (a *Q) get_f5() map[int]int {
	if a == nil {
		return nil
	}
	return a.f5
}
func (a *Q) get_f6() int64 {
	if a == nil {
		return 0
	}
	return a.f6
}
func (a *Q) get_f7() *Q {
	if a == nil {
		return nil
	}
	return a.f7
}
func (a *Q) get_f8() bool {
	if a == nil {
		return false
	}
	return a.f8
}
func (a *Q) get_f9() []byte {
	if a == nil {
		return nil
	}
	return a.f9
}
func (a *Q) get_f0() uint32 {
	if a == nil {
		return 0
	}
	return a.f0
}

func BenchmarkField(b *testing.B) {
	q := Q{
		f1: 10,
		f2: "10",
		f3: 10,
		f4: []int32{10},
		f5: map[int]int{10: 10},
		f6: 10,
		f7: &Q{},
		f8: false,
		f9: []byte{10},
		f0: 10,
	}
	w := &Q{}
	for i := 0; i < b.N; i++ {
		w.f1 = q.f1
		w.f2 = q.f2
		w.f3 = q.f3
		w.f4 = q.f4
		w.f5 = q.f5
		w.f6 = q.f6
		w.f7 = q.f7
		w.f8 = q.f8
		w.f9 = q.f9
		w.f0 = q.f0
	}
}

func BenchmarkInterface(b *testing.B) {
	q := Q{
		f1: 10,
		f2: "10",
		f3: 10,
		f4: []int32{10},
		f5: map[int]int{10: 10},
		f6: 10,
		f7: &Q{},
		f8: false,
		f9: []byte{10},
		f0: 10,
	}
	w := &Q{}
	for i := 0; i < b.N; i++ {
		w.f1 = q.get_f1()
		w.f2 = q.get_f2()
		w.f3 = q.get_f3()
		w.f4 = q.get_f4()
		w.f5 = q.get_f5()
		w.f6 = q.get_f6()
		w.f7 = q.get_f7()
		w.f8 = q.get_f8()
		w.f9 = q.get_f9()
		w.f0 = q.get_f0()
	}
}

type EmptyStruct struct{}

func (a *EmptyStruct) U() int {
	return 10
}

func BenchmarkEmptyStructPool(b *testing.B) {
	pool := sync.Pool{
		New: func() any { return &EmptyStruct{} },
	}
	for i := 0; i < b.N; i++ {
		q := pool.Get().(*EmptyStruct)
		collect(q)
		pool.Put(q)
	}
}
func BenchmarkEmptyStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := &EmptyStruct{}
		collect(q)
	}
}
func BenchmarkEmptyStructConst(b *testing.B) {
	e := &EmptyStruct{}
	for i := 0; i < b.N; i++ {
		collect(e)
	}
}

var lastEmptyStruct *EmptyStruct

func collect(e *EmptyStruct) {
	lastEmptyStruct = e
}

type R struct{}

func (a *R) Reset() {}

func BenchmarkReset(b *testing.B) {
	r := &R{}
	for i := 0; i < b.N; i++ {
		r.Reset()
	}
}
func BenchmarkResetNil(b *testing.B) {
	var r *R = nil
	for i := 0; i < b.N; i++ {
		r.Reset()
	}
}

type Ii interface {
	Q()
}
type Q0 struct {
	i Ii
}
type Q1 struct{}
type Q2 struct{}

func (a *Q1) Q() {}
func (a *Q2) Q() {}

func TestCastType(t *testing.T) {
	q := &Q0{}
	if q.i != nil {
		q.i.Q()
	}
}
