package main

import (
	"bytes"
	"testing"

	gogo "bench/proto/gogo/bench"
	google "bench/proto/google/bench"
	litepb "bench/proto/litepb_old/bench"
	litepb2 "github.com/e-tape/litepb/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

// RUN `make test-compile-for-bench` for init test proto

var decodeData []byte
var decodeModel = &google.Bench{
	Uuid: &litepb2.UUID{
		Value: []byte{18, 52, 86, 120, 21, 117, 69, 117, 101, 117, 135, 83, 19, 84, 21, 64},
	},
	Uint64:  2065657434543,
	Uint32:  156547,
	String_: "123456",
	//Smap: map[int32]*google.Bench_InnerForMap{
	//	1: {
	//		Uint64: 1,
	//		Uint32: 20,
	//	},
	//	2: {
	//		Uint64: 1000,
	//		Uint32: 2000,
	//	},
	//	9: {
	//		Uint64: 99,
	//		Uint32: 99,
	//	},
	//	14: {
	//		Uint64: 14,
	//		Uint32: 14,
	//	},
	//	72: {
	//		Uint64: 72,
	//		Uint32: 72,
	//	},
	//	602: {
	//		Uint64: 602,
	//		Uint32: 602,
	//	},
	//},
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
	Fixed32: 32,
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
			decodeModel.Fixed32 != model.Fixed32 {
			panic(`eq`)
		}
	}
}

func BenchmarkSimpleLiteOldPb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		model := &litepb.Bench{}
		if err := model.UnmarshalProto(decodeData); err != nil {
			panic(err)
		}
		if !(bytes.Equal(decodeModel.GetUuid().GetValue(), model.Uuid[:]) || model.Uuid == uuid.Nil) ||
			decodeModel.Uint32 != model.Uint32 ||
			decodeModel.Uint64 != model.Uint64 ||
			decodeModel.String_ != model.String_ ||
			len(decodeModel.Smap) != len(model.Smap) ||
			len(decodeModel.Iarr) != len(model.Iarr) ||
			decodeModel.Fixed32 != model.Fixed32 {
			panic(`eq`)
		}
	}
}

func BenchmarkSimpleLiteOldPbReturnToPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		model := litepb.NewBench()
		if err := model.UnmarshalProto(decodeData); err != nil {
			panic(err)
		}
		if !(bytes.Equal(decodeModel.GetUuid().GetValue(), model.Uuid[:]) || model.Uuid == uuid.Nil) ||
			decodeModel.Uint32 != model.Uint32 ||
			decodeModel.Uint64 != model.Uint64 ||
			decodeModel.String_ != model.String_ ||
			len(decodeModel.Smap) != len(model.Smap) ||
			len(decodeModel.Iarr) != len(model.Iarr) ||
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
				decodeModel.Fixed32 != model.Fixed32 {
				panic(`eq`)
			}
		}
	})
}

func BenchmarkParallelLiteOldPb(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			model := &litepb.Bench{}
			if err := model.UnmarshalProto(decodeData); err != nil {
				panic(err)
			}
			if !(bytes.Equal(decodeModel.GetUuid().GetValue(), model.Uuid[:]) || model.Uuid == uuid.Nil) ||
				decodeModel.Uint32 != model.Uint32 ||
				decodeModel.Uint64 != model.Uint64 ||
				decodeModel.String_ != model.String_ ||
				len(decodeModel.Smap) != len(model.Smap) ||
				len(decodeModel.Iarr) != len(model.Iarr) ||
				decodeModel.Fixed32 != model.Fixed32 {
				panic(`eq`)
			}
		}
	})
}

func BenchmarkParallelLiteOldPbReturnToPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			model := litepb.NewBench()
			if err := model.UnmarshalProto(decodeData); err != nil {
				panic(err)
			}
			if !(bytes.Equal(decodeModel.GetUuid().GetValue(), model.Uuid[:]) || model.Uuid == uuid.Nil) ||
				decodeModel.Uint32 != model.Uint32 ||
				decodeModel.Uint64 != model.Uint64 ||
				decodeModel.String_ != model.String_ ||
				len(decodeModel.Smap) != len(model.Smap) ||
				len(decodeModel.Iarr) != len(model.Iarr) ||
				decodeModel.Fixed32 != model.Fixed32 {
				panic(`eq`)
			}
			model.ReturnToPool()
		}
	})
}

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
			model := &litepb.Bench{}
			model.Iarr = []*litepb.Bench_InnerForMap{
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
			model := litepb.NewBench()
			model.SetIarr([]litepb.IBench_InnerForMapGet{
				litepb.NewBench_InnerForMap().SetUint32(2).SetUint64(5).SetUint64(654),
				litepb.NewBench_InnerForMap().SetUint64(6),
				litepb.NewBench_InnerForMap(),
			})
			model.ReturnToPool()
		}
	})
}
