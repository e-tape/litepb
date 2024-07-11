package main

import (
	"bytes"
	"testing"

	gogo "bench/proto/gogo/bench"
	google "bench/proto/google/bench"
	litepbNoPool "bench/proto/litepb_no_pool/bench"
	litepbPool "bench/proto/litepb_pool/bench"
	litepb2 "github.com/e-tape/litepb/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

// RUN `make test-compile-for-bench` for init test proto
// go test -bench . -benchmem -benchtime 2s -cpuprofile cpu.prof
// go tool pprof -http :7000 ./cpu.prof

var _ = litepb2.UUID{}
var decodeData []byte
var decodeModel = &google.Bench{
	Uuid: &litepb2.UUID{
		Value: []byte{18, 52, 86, 120, 21, 117, 69, 117, 101, 117, 135, 83, 19, 84, 21, 64},
	},
	Uint64:  2065657434543,
	Uint32:  156547,
	String_: "123456",
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
	Ifm: &google.Bench_InnerForMap{
		Uint64: 64,
		Uint32: 32,
	},
	Fixed64: 2065657434543,
	Fixed32: 156547,
	RInt32:  []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
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
		if !bytes.Equal(decodeModel.GetUuid().GetValue(), model.GetUuid().GetValue()) ||
			decodeModel.Uint32 != model.Uint32 ||
			decodeModel.Uint64 != model.Uint64 ||
			decodeModel.String_ != model.String_ ||
			len(decodeModel.Smap) != len(model.Smap) ||
			len(decodeModel.Iarr) != len(model.Iarr) ||
			decodeModel.GetIfm().GetUint32() != model.GetIfm().GetUint32() ||
			decodeModel.GetIfm().GetUint64() != model.GetIfm().GetUint64() ||
			len(decodeModel.RInt32) != len(model.RInt32) ||
			decodeModel.Fixed32 != model.Fixed32 {
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
		if !bytes.Equal(decodeModel.GetUuid().GetValue(), model.GetUuid().GetValue()) ||
			decodeModel.Uint32 != model.Uint32 ||
			decodeModel.Uint64 != model.Uint64 ||
			decodeModel.String_ != model.String_ ||
			len(decodeModel.Smap) != len(model.Smap) ||
			len(decodeModel.Iarr) != len(model.Iarr) ||
			decodeModel.GetIfm().GetUint32() != model.GetIfm().GetUint32() ||
			decodeModel.GetIfm().GetUint64() != model.GetIfm().GetUint64() ||
			len(decodeModel.RInt32) != len(model.RInt32) ||
			decodeModel.Fixed32 != model.Fixed32 {
			panic(`eq`)
		}
	}
}

func BenchmarkSimpleLitePb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		model := litepbPool.NewBench()
		if err := model.UnmarshalProto(decodeData); err != nil {
			panic(err)
		}
		if !(bytes.Equal(decodeModel.GetUuid().GetValue(), model.Uuid[:]) || model.Uuid == uuid.Nil) ||
			decodeModel.Uint32 != model.Uint32 ||
			decodeModel.Uint64 != model.Uint64 ||
			decodeModel.String_ != model.String_ ||
			len(decodeModel.Smap) != len(model.Smap) ||
			len(decodeModel.Iarr) != len(model.Iarr) ||
			decodeModel.GetIfm().GetUint32() != model.GetIfm().GetUint32() ||
			decodeModel.GetIfm().GetUint64() != model.GetIfm().GetUint64() ||
			len(decodeModel.RInt32) != len(model.RInt32) ||
			decodeModel.Fixed32 != model.Fixed32 {
			panic(`eq`)
		}
		model.ReturnToPool()
	}
}

func BenchmarkSimpleLitePbNoPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		model := litepbNoPool.NewBench()
		if err := model.UnmarshalProto(decodeData); err != nil {
			panic(err)
		}
		if !(bytes.Equal(decodeModel.GetUuid().GetValue(), model.Uuid[:]) || model.Uuid == uuid.Nil) ||
			decodeModel.Uint32 != model.Uint32 ||
			decodeModel.Uint64 != model.Uint64 ||
			decodeModel.String_ != model.String_ ||
			len(decodeModel.Smap) != len(model.Smap) ||
			len(decodeModel.Iarr) != len(model.Iarr) ||
			decodeModel.GetIfm().GetUint32() != model.GetIfm().GetUint32() ||
			decodeModel.GetIfm().GetUint64() != model.GetIfm().GetUint64() ||
			len(decodeModel.RInt32) != len(model.RInt32) ||
			decodeModel.Fixed32 != model.Fixed32 {
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
			if !bytes.Equal(decodeModel.GetUuid().GetValue(), model.GetUuid().GetValue()) ||
				decodeModel.Uint32 != model.Uint32 ||
				decodeModel.Uint64 != model.Uint64 ||
				decodeModel.String_ != model.String_ ||
				len(decodeModel.Smap) != len(model.Smap) ||
				len(decodeModel.Iarr) != len(model.Iarr) ||
				decodeModel.GetIfm().GetUint32() != model.GetIfm().GetUint32() ||
				decodeModel.GetIfm().GetUint64() != model.GetIfm().GetUint64() ||
				len(decodeModel.RInt32) != len(model.RInt32) ||
				decodeModel.Fixed32 != model.Fixed32 {
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
			if !bytes.Equal(decodeModel.GetUuid().GetValue(), model.GetUuid().GetValue()) ||
				decodeModel.Uint32 != model.Uint32 ||
				decodeModel.Uint64 != model.Uint64 ||
				decodeModel.String_ != model.String_ ||
				len(decodeModel.Smap) != len(model.Smap) ||
				len(decodeModel.Iarr) != len(model.Iarr) ||
				decodeModel.GetIfm().GetUint32() != model.GetIfm().GetUint32() ||
				decodeModel.GetIfm().GetUint64() != model.GetIfm().GetUint64() ||
				len(decodeModel.RInt32) != len(model.RInt32) ||
				decodeModel.Fixed32 != model.Fixed32 {
				panic(`eq`)
			}
		}
	})
}

func BenchmarkParallelLitePb(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			model := litepbPool.NewBench()
			if err := model.UnmarshalProto(decodeData); err != nil {
				panic(err)
			}
			if !(bytes.Equal(decodeModel.GetUuid().GetValue(), model.Uuid[:]) || model.Uuid == uuid.Nil) ||
				decodeModel.Uint32 != model.Uint32 ||
				decodeModel.Uint64 != model.Uint64 ||
				decodeModel.String_ != model.String_ ||
				len(decodeModel.Smap) != len(model.Smap) ||
				len(decodeModel.Iarr) != len(model.Iarr) ||
				decodeModel.GetIfm().GetUint32() != model.GetIfm().GetUint32() ||
				decodeModel.GetIfm().GetUint64() != model.GetIfm().GetUint64() ||
				len(decodeModel.RInt32) != len(model.RInt32) ||
				decodeModel.Fixed32 != model.Fixed32 {
				panic(`eq`)
			}
			model.ReturnToPool()
		}
	})
}

func BenchmarkParallelLitePbNoPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			model := litepbNoPool.NewBench()
			if err := model.UnmarshalProto(decodeData); err != nil {
				panic(err)
			}
			if !(bytes.Equal(decodeModel.GetUuid().GetValue(), model.Uuid[:]) || model.Uuid == uuid.Nil) ||
				decodeModel.Uint32 != model.Uint32 ||
				decodeModel.Uint64 != model.Uint64 ||
				decodeModel.String_ != model.String_ ||
				len(decodeModel.Smap) != len(model.Smap) ||
				len(decodeModel.Iarr) != len(model.Iarr) ||
				decodeModel.GetIfm().GetUint32() != model.GetIfm().GetUint32() ||
				decodeModel.GetIfm().GetUint64() != model.GetIfm().GetUint64() ||
				len(decodeModel.RInt32) != len(model.RInt32) ||
				decodeModel.Fixed32 != model.Fixed32 {
				panic(`eq`)
			}
		}
	})
}
