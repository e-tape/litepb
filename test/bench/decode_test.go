package main

import (
	"testing"

	"google.golang.org/protobuf/proto"
	gogo "test/bench/proto/gogo/bench"
	google "test/bench/proto/google/bench"
	litepb "test/bench/proto/litepb/bench"
)

var decodeData []byte
var decodeModel = &google.Bench{
	Uint64: 2065657434543,
	Uint32: 156547,
	//String_: "123456",
	Smap: map[int32]*google.Bench_InnerForMap{
		1: {
			Uint64: 1,
			Uint32: 20,
		},
		2: {
			Uint64: 1000,
			Uint32: 2000,
		},
		9: {
			Uint64: 99,
			Uint32: 99,
		},
		14: {
			Uint64: 14,
			Uint32: 14,
		},
		72: {
			Uint64: 72,
			Uint32: 72,
		},
		602: {
			Uint64: 602,
			Uint32: 602,
		},
	},
	Iarr: []*google.Bench_InnerForMap{
		{
			Uint64: 1001,
			Uint32: 1001,
		},
		{
			Uint64: 2002,
			Uint32: 2002,
		},
		{
			Uint64: 62002,
			Uint32: 62002,
		},
	},
}

func init() {
	d, err := proto.Marshal(decodeModel)
	if err != nil {
		panic(err)
	}
	decodeData = d
}

func BenchmarkSimpleGoogle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		model := &google.Bench{}
		if err := proto.Unmarshal(decodeData, model); err != nil {
			panic(err)
		}
		if decodeModel.Uint32 != model.Uint32 ||
			decodeModel.Uint64 != model.Uint64 ||
			decodeModel.String_ != model.String_ ||
			len(decodeModel.Smap) != len(model.Smap) ||
			len(decodeModel.Iarr) != len(model.Iarr) {
			panic(`eq`)
		}
	}
}

func BenchmarkSimpleGogo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		model := &gogo.Bench{}
		if err := model.Unmarshal(decodeData); err != nil {
			panic(err)
		}
		if decodeModel.Uint32 != model.Uint32 ||
			decodeModel.Uint64 != model.Uint64 ||
			decodeModel.String_ != model.String_ ||
			len(decodeModel.Smap) != len(model.Smap) ||
			len(decodeModel.Iarr) != len(model.Iarr) {
			panic(`eq`)
		}
	}
}

func BenchmarkSimpleLitePb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		model := &litepb.Bench{}
		if err := model.Unmarshal(decodeData); err != nil {
			panic(err)
		}
		if decodeModel.Uint32 != model.Uint32 ||
			decodeModel.Uint64 != model.Uint64 ||
			decodeModel.String_ != model.String_ ||
			len(decodeModel.Smap) != len(model.Smap) ||
			len(decodeModel.Iarr) != len(model.Iarr) {
			panic(`eq`)
		}
	}
}

func BenchmarkSimpleLitePbReturnToPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		model := litepb.NewBench()
		if err := model.Unmarshal(decodeData); err != nil {
			panic(err)
		}
		if decodeModel.Uint32 != model.Uint32 ||
			decodeModel.Uint64 != model.Uint64 ||
			decodeModel.String_ != model.String_ ||
			len(decodeModel.Smap) != len(model.Smap) ||
			len(decodeModel.Iarr) != len(model.Iarr) {
			panic(`eq`)
		}
		model.ReturnToPool()
	}
}

func BenchmarkParallelGoogle(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			model := &google.Bench{}
			if err := proto.Unmarshal(decodeData, model); err != nil {
				panic(err)
			}
			if decodeModel.Uint32 != model.Uint32 ||
				decodeModel.Uint64 != model.Uint64 ||
				decodeModel.String_ != model.String_ ||
				len(decodeModel.Smap) != len(model.Smap) ||
				len(decodeModel.Iarr) != len(model.Iarr) {
				panic(`eq`)
			}
		}
	})
}

func BenchmarkParallelGogo(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			model := &gogo.Bench{}
			if err := model.Unmarshal(decodeData); err != nil {
				panic(err)
			}
			if decodeModel.Uint32 != model.Uint32 ||
				decodeModel.Uint64 != model.Uint64 ||
				decodeModel.String_ != model.String_ ||
				len(decodeModel.Smap) != len(model.Smap) ||
				len(decodeModel.Iarr) != len(model.Iarr) {
				panic(`eq`)
			}
		}
	})
}

func BenchmarkParallelLitePb(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			model := &litepb.Bench{}
			if err := model.Unmarshal(decodeData); err != nil {
				panic(err)
			}
			if decodeModel.Uint32 != model.Uint32 ||
				decodeModel.Uint64 != model.Uint64 ||
				decodeModel.String_ != model.String_ ||
				len(decodeModel.Smap) != len(model.Smap) ||
				len(decodeModel.Iarr) != len(model.Iarr) {
				panic(`eq`)
			}
		}
	})
}

func BenchmarkParallelLitePbReturnToPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			model := litepb.NewBench()
			if err := model.Unmarshal(decodeData); err != nil {
				panic(err)
			}
			if decodeModel.Uint32 != model.Uint32 ||
				decodeModel.Uint64 != model.Uint64 ||
				decodeModel.String_ != model.String_ ||
				len(decodeModel.Smap) != len(model.Smap) ||
				len(decodeModel.Iarr) != len(model.Iarr) {
				panic(`eq`)
			}
			model.ReturnToPool()
		}
	})
}
