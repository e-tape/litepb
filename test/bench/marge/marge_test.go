package main

import (
	"slices"
	"testing"

	gogo "bench/proto/gogo/bench"
	google "bench/proto/google/bench"
	litepbNoPool "bench/proto/litepb_no_pool/bench"
	"google.golang.org/protobuf/proto"
)

func TestMarge(t *testing.T) {
	d1, err := proto.Marshal(&google.Bench{Ifm: &google.Bench_InnerForMap{Uint64: 1}})
	if err != nil {
		panic(err)
	}
	d2, err := proto.Marshal(&google.Bench{Ifm: &google.Bench_InnerForMap{Uint32: 2}})
	if err != nil {
		panic(err)
	}
	{
		model := &google.Bench{}
		if err = proto.Unmarshal(d1, model); err != nil {
			panic(err)
		}
		if model.GetIfm().GetUint64() != 1 && model.GetIfm().GetUint32() != 0 {
			panic("eq")
		}
	}
	{
		model := &google.Bench{}
		if err = proto.Unmarshal(d2, model); err != nil {
			panic(err)
		}
		if model.GetIfm().GetUint64() != 0 && model.GetIfm().GetUint32() != 2 {
			panic("eq")
		}
	}
	{
		model := &google.Bench{}
		if err = proto.Unmarshal(slices.Concat(d2, d1), model); err != nil {
			panic(err)
		}
		if model.GetIfm().GetUint64() != 1 && model.GetIfm().GetUint32() != 2 {
			panic("eq")
		}
	}
	{
		model := &gogo.Bench{}
		if err = model.Unmarshal(d1); err != nil {
			panic(err)
		}
		if model.GetIfm().GetUint64() != 1 && model.GetIfm().GetUint32() != 0 {
			panic("eq")
		}
	}
	{
		model := &gogo.Bench{}
		if err = model.Unmarshal(d2); err != nil {
			panic(err)
		}
		if model.GetIfm().GetUint64() != 0 && model.GetIfm().GetUint32() != 2 {
			panic("eq")
		}
	}
	{
		model := &gogo.Bench{}
		if err = model.Unmarshal(slices.Concat(d2, d1)); err != nil {
			panic(err)
		}
		if model.GetIfm().GetUint64() != 1 && model.GetIfm().GetUint32() != 2 {
			panic("eq")
		}
	}
	{
		model := &litepbNoPool.Bench{}
		if err = model.UnmarshalProto(d1); err != nil {
			panic(err)
		}
		if model.GetIfm().GetUint64() != 1 && model.GetIfm().GetUint32() != 0 {
			panic("eq")
		}
	}
	{
		model := &litepbNoPool.Bench{}
		if err = model.UnmarshalProto(d2); err != nil {
			panic(err)
		}
		if model.GetIfm().GetUint64() != 0 && model.GetIfm().GetUint32() != 2 {
			panic("eq")
		}
	}
	{
		model := &litepbNoPool.Bench{}
		if err = model.UnmarshalProto(slices.Concat(d2, d1)); err != nil {
			panic(err)
		}
		if model.GetIfm().GetUint64() != 1 && model.GetIfm().GetUint32() != 2 {
			panic("eq")
		}
	}
}
