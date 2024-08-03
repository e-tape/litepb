package main

import (
	"fmt"
	"testing"

	"google.golang.org/protobuf/proto"
	"test/generated/google"
)

var data []byte

func init() {
	r := &google.R1{
		U1: 203,
		U2: 15,
	}
	d, err := proto.Marshal(r)
	if err != nil {
		panic(err)
	}
	data = d
}

func BenchmarkSelf(b *testing.B) {
	r := &google.R1{}
	for i := 0; i < b.N; i++ {
		var index int32
		{
			next, _, _ := decodeIndex(data[index:])
			index += next
		}
		{
			next, value := decodeUin64(data[index:])
			r.U1 = value
			index += next
		}
		{
			next, _, _ := decodeIndex(data[index:])
			index += next
		}
		{
			next, value := decodeUin32(data[index:])
			r.U2 = value
			index += next
		}
	}
}

func BenchmarkProto(b *testing.B) {
	r := &google.R1{}
	for i := 0; i < b.N; i++ {
		_ = proto.Unmarshal(data, r)
	}
}

func Test(*testing.T) {
	r := &google.R1{
		U1: 203,
		U2: 15,
	}
	data, err := proto.Marshal(r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%b\n", data)
	fmt.Printf("%d\n", data)

	r2 := &google.R1{}
	var index int32
	{
		next, number, wType := decodeIndex(data[index:])
		//_, _ = number, wType
		fmt.Println(index, number, wType)
		index += next
	}
	{
		next, value := decodeUin64(data[index:])
		fmt.Println(index, value)
		r2.U1 = value
		index += next
	}
	{
		next, number, wType := decodeIndex(data[index:])
		//_, _ = number, wType
		fmt.Println(index, number, wType)
		index += next
	}
	{
		next, value := decodeUin32(data[index:])
		r2.U2 = value
		fmt.Println(index, value)
		index += next
	}
	fmt.Println(r2)
}

// 127 - 0b01111111
// 128 - 0b10000000

func decodeIndex(data []byte) (index_ int32, fieldNum_ int32, wireType_ int32) {
	l := len(data)
	var fieldNum int32
	if l > 0 {
		index_++
		fieldNum |= int32(data[0] & 127)
		if data[0] < 128 {
			goto end
		}
	}
	if l > 1 {
		index_++
		fieldNum |= int32(data[1]&127) << 7
		if data[1] < 128 {
			goto end
		}
	}
	if l > 2 {
		index_++
		fieldNum |= int32(data[2]&127) << 14
		if data[2] < 128 {
			goto end
		}
	}
	if l > 3 {
		index_++
		fieldNum |= int32(data[3]&127) << 21
		if data[3] < 128 {
			goto end
		}
	}
	panic(`ErrIntOverflowProto`)
end:
	return index_, fieldNum >> 3, fieldNum & 0x7
}
func decodeUin32(data []byte) (index_ int32, result_ uint32) {
	l := len(data)
	if l > 0 {
		index_++
		result_ |= uint32(data[0] & 127)
		if data[0] < 128 {
			goto end
		}
	}
	if l > 1 {
		index_++
		result_ |= uint32(data[1]&127) << 7
		if data[1] < 128 {
			goto end
		}
	}
	if l > 2 {
		index_++
		result_ |= uint32(data[2]&127) << 14
		if data[2] < 128 {
			goto end
		}
	}
	if l > 3 {
		index_++
		result_ |= uint32(data[3]&127) << 21
		if data[3] < 128 {
			goto end
		}
	}
	panic(`ErrIntOverflowProto`)
end:
	return index_, result_
}
func decodeUin64(data []byte) (index_ int32, result_ uint64) {
	l := len(data)
	if l > 0 {
		index_++
		result_ |= uint64(data[0] & 127)
		if data[0] < 128 {
			goto end
		}
	}
	if l > 1 {
		index_++
		result_ |= uint64(data[1]&127) << 7
		if data[1] < 128 {
			goto end
		}
	}
	if l > 2 {
		index_++
		result_ |= uint64(data[2]&127) << 14
		if data[2] < 128 {
			goto end
		}
	}
	if l > 3 {
		index_++
		result_ |= uint64(data[3]&127) << 21
		if data[3] < 128 {
			goto end
		}
	}
	if l > 4 {
		index_++
		result_ |= uint64(data[4]&127) << 28
		if data[4] < 128 {
			goto end
		}
	}
	if l > 5 {
		index_++
		result_ |= uint64(data[5]&127) << 35
		if data[5] < 128 {
			goto end
		}
	}
	if l > 6 {
		index_++
		result_ |= uint64(data[6]&127) << 42
		if data[6] < 128 {
			goto end
		}
	}
	if l > 7 {
		index_++
		result_ |= uint64(data[7]&127) << 49
		if data[7] < 128 {
			goto end
		}
	}
	if l > 8 {
		index_++
		result_ |= uint64(data[8]&127) << 56
		if data[8] < 128 {
			goto end
		}
	}
	if l > 9 {
		index_++
		result_ |= uint64(data[9]&127) << 63
		if data[9] < 128 {
			goto end
		}
	}
	panic(`ErrIntOverflowProto`)
end:
	return index_, result_
}

//func benchIfGoto(bm bench) func(*testing.B) {
//	return func(bt *testing.B) {
//		for i := 0; i < bt.N; i++ {
//			l := len(bm.inData)
//			var i64 int64
//			if l > 0 {
//				i64 |= int64(bm.inData[0] & 0b01111111)
//				if bm.inData[0] < 0b10000000 {
//					goto end
//				}
//			}
//			if l > 1 {
//				i64 |= int64(bm.inData[1]&0b01111111) << 7
//				if bm.inData[1] < 0b10000000 {
//					goto end
//				}
//			}
//			if l > 2 {
//				i64 |= int64(bm.inData[2]&0b01111111) << 14
//				if bm.inData[2] < 0b10000000 {
//					goto end
//				}
//			}
//			if l > 3 {
//				i64 |= int64(bm.inData[3]&0b01111111) << 21
//				if bm.inData[3] < 0b10000000 {
//					goto end
//				}
//			}
//			if l > 4 {
//				i64 |= int64(bm.inData[4]&0b01111111) << 28
//				if bm.inData[4] < 0b10000000 {
//					goto end
//				}
//			}
//			if l > 5 {
//				i64 |= int64(bm.inData[5]&0b01111111) << 35
//				if bm.inData[5] < 0b10000000 {
//					goto end
//				}
//			}
//			if l > 6 {
//				i64 |= int64(bm.inData[6]&0b01111111) << 42
//				if bm.inData[6] < 0b10000000 {
//					goto end
//				}
//			}
//			if l > 7 {
//				i64 |= int64(bm.inData[7]&0b01111111) << 49
//				if bm.inData[7] < 0b10000000 {
//					goto end
//				}
//			}
//			if l > 8 {
//				i64 |= int64(bm.inData[8]&0b01111111) << 56
//				if bm.inData[8] < 0b10000000 {
//					goto end
//				}
//			}
//			if l > 9 {
//				i64 |= int64(bm.inData[9]&0b01111111) << 63
//				if bm.inData[9] < 0b10000000 {
//					goto end
//				}
//			}
//			panic(`ErrIntOverflowProto`)
//		end:
//			if i64 != bm.inI64 {
//				panic(`mismatch`)
//			}
//		}
//	}
//}
