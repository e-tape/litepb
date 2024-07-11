package main

import (
	"testing"

	gogo "bench/proto/gogo/bench"
	litepbNoPool "bench/proto/litepb_no_pool/bench"
	litepbPool "bench/proto/litepb_pool/bench"
)

func BenchmarkParallelGogoField(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			model := &gogo.Bench{}
			model.Iarr = []*gogo.Bench_InnerForMap{
				{},
				{},
				{},
			}
		}
	})
}

func BenchmarkParallelField(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			model := &litepbNoPool.Bench{}
			model.Iarr = []*litepbNoPool.Bench_InnerForMap{
				{},
				{},
				{},
			}
		}
	})
}

func BenchmarkParallelSetter(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			model := litepbPool.NewBench()
			model.SetIarr([]litepbPool.IBench_InnerForMapGet{
				litepbPool.NewBench_InnerForMap().SetUint32(2).SetUint64(5).SetUint64(654),
				litepbPool.NewBench_InnerForMap().SetUint64(6),
				litepbPool.NewBench_InnerForMap(),
			})
			model.ReturnToPool()
		}
	})
}
