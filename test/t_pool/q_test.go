package t_pool

import (
	"sync"
	"sync/atomic"
	"testing"
)

type Q struct {
	Int    int
	Int32  int32
	Int64  int64
	RInt64 []int64
}

var globalSyncPool = &sync.Pool{
	New: func() any { return &Q{} },
}

type ChPool struct {
	ch chan *Q
}

func (a *ChPool) Get() *Q {
	select {
	case obj := <-a.ch:
		return obj
	default:
		return &Q{}
	}
}

func (a *ChPool) Put(item *Q) {
	select {
	case a.ch <- item:
	default:
	}
}

const chPoolCount = 100

var chPool = &ChPool{
	ch: make(chan *Q, chPoolCount),
}

func init() {
	for i := 0; i < chPoolCount; i++ {
		chPool.Put(&Q{})
	}
}

const ringPoolCount = 100

type RingPool struct {
	buffer     [ringPoolCount]*Q
	readIndex  uint64
	writeIndex uint64
	size       uint64
}

func (a *RingPool) Put(item *Q) bool {
	for {
		r := atomic.LoadUint64(&a.readIndex)
		w := atomic.LoadUint64(&a.writeIndex)
		if (w+1)%a.size == r%a.size {
			return false
		}
		if atomic.CompareAndSwapUint64(&a.writeIndex, w, w+1) {
			a.buffer[w%a.size] = item
			return true
		}
	}
}

func (a *RingPool) Get() *Q {
	for {
		r := atomic.LoadUint64(&a.readIndex)
		w := atomic.LoadUint64(&a.writeIndex)
		if r == w {
			return &Q{}
		}
		if atomic.CompareAndSwapUint64(&a.readIndex, r, r+1) {
			return a.buffer[r%a.size]
		}
	}
}

var ringPool = &RingPool{
	size: ringPoolCount,
}

var threadPools sync.Map

func getLocalMapPool() *sync.Pool {
	pool, _ := threadPools.LoadOrStore(getThreadId(), &sync.Pool{
		New: func() any {
			return &Q{}
		},
	})
	return pool.(*sync.Pool)
}

/////////////

func BenchmarkSimpleGlobalSyncPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := globalSyncPool.Get().(*Q)
		q.Int = i
		globalSyncPool.Put(q)
	}
}

func BenchmarkSimpleLocalSyncPool(b *testing.B) {
	localSyncPool := &sync.Pool{
		New: func() any { return &Q{} },
	}
	for i := 0; i < b.N; i++ {
		q := localSyncPool.Get().(*Q)
		q.Int = i
		localSyncPool.Put(q)
	}
}

func BenchmarkSimpleChanPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := chPool.Get()
		q.Int = i
		chPool.Put(q)
	}
}

func BenchmarkSimpleRingPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := ringPool.Get()
		q.Int = i
		ringPool.Put(q)
	}
}

func BenchmarkSimpleLocalMapPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pool := getLocalMapPool()
		obj := pool.Get().(*Q)
		obj.Int = i
		pool.Put(obj)
	}
}

func BenchmarkAsyncGlobalSyncPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q := globalSyncPool.Get().(*Q)
			q.Int = 5
			globalSyncPool.Put(q)
		}
	})
}

func BenchmarkAsyncLocalSyncPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		localSyncPool := &sync.Pool{
			New: func() any { return &Q{} },
		}
		for pb.Next() {
			q := localSyncPool.Get().(*Q)
			q.Int = 5
			localSyncPool.Put(q)
		}
	})
}

func BenchmarkAsyncChanPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q := chPool.Get()
			q.Int = 5
			chPool.Put(q)
		}
	})
}

func BenchmarkAsyncRingPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q := ringPool.Get()
			q.Int = 5
			ringPool.Put(q)
		}
	})
}

func BenchmarkAsyncLocalMapPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		pool := getLocalMapPool()
		for pb.Next() {
			obj := pool.Get().(*Q)
			obj.Int = 5
			pool.Put(obj)
		}
	})
}
