package t14

import (
	"math/big"
	"sync"
	"testing"

	"github.com/google/uuid"
)

var gUuid uuid.UUID
var bb = []byte{0, 0, 18, 52, 86, 120, 21, 117, 69, 117, 101, 117, 135, 83, 19, 84, 5, 69, 0, 0}
var originUuid = uuid.Must(uuid.FromBytes(bb[2:18]))

func BenchmarkUuidParseString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		u := uuid.MustParse("12345678-1575-4575-6575-875313540545")
		gUuid = u
	}
}
func BenchmarkUuidFromBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		u, err := uuid.FromBytes(bb[2:18])
		if err != nil {
			panic(err)
		}
		if u != originUuid {
			panic("eq")
		}
		gUuid = u
	}
}
func BenchmarkUuidRefBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		u := uuid.UUID(bb[2:18])
		if u != originUuid {
			panic("eq")
		}
		gUuid = u
	}
}

var uuidPool = &sync.Pool{New: func() any { return new(uuid.UUID) }}

func BenchmarkUuidFromBytesPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		u := uuidPool.Get().(*uuid.UUID)
		err := u.UnmarshalBinary(bb[2:18])
		if err != nil {
			panic(err)
		}
		if *u != originUuid {
			panic("eq")
		}
		clear((*u)[:])
		uuidPool.Put(u)
		//gUuid = *u
	}
}

func BenchmarkUuidEq(b *testing.B) {
	for i := 0; i < b.N; i++ {
		u, err := uuid.FromBytes(bb[2:18])
		if err != nil {
			panic(err)
		}
		if u == [16]byte{} {
			panic(1)
		}
		if u != originUuid {
			panic("eq")
		}
	}
}

var gBig *big.Int

var bigPool = &sync.Pool{
	New: func() interface{} {
		return new(big.Int)
	},
}

type QBig struct {
	Value *big.Int
}

func initQBigPool(x uint64) QBig {
	u := bigPool.Get().(*big.Int)
	return QBig{
		Value: u.SetUint64(x),
	}
}

func initQBigNew(x uint64) QBig {
	return QBig{
		Value: big.NewInt(int64(x)),
	}
}

func BenchmarkBigIntPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		u := initQBigPool(1234567890)
		//u := bigPool.Get().(*big.Int)
		//u.SetInt64(1234567890)
		u.Value.SetUint64(0)
		uuidPool.Put(u.Value)
		gBig = u.Value
	}
}

func BenchmarkBigIntNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		u := initQBigNew(1234567890)
		u.Value.SetUint64(0)
		_ = u
		gBig = u.Value
	}
}

var parseBb = []byte{0, 1, 2, 3, 0, 2, 2, 2, 10, 0, 4, 4, 4, 0}

func BenchmarkParseLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		index := 0
		var key, val, pb byte
		var item *byte
		for len(parseBb) > index {
			pb = parseBb[index]
			index++
			switch pb {
			case 0:
				item = nil
			case 1:
				if item == nil {
					item = &key
				} else {
					*item += pb
				}
			case 2:
				if item == nil {
					item = &val
				} else {
					*item += pb
				}
			default:
				if item != nil {
					*item += pb
				}
			}
		}
		if key != 5 {
			panic("key")
		}
		if val != 14 {
			panic("val")
		}
	}
}
func BenchmarkParseLoopConst(b *testing.B) {
	for i := 0; i < b.N; i++ {
		index := 0
		var key, val, set, pb byte
		var item *byte
		for set != 2 {
			if len(parseBb) <= index {
				break
			}
			pb = parseBb[index]
			index++
			switch pb {
			case 0:
				if item != nil {
					set++
				}
				item = nil
			case 1:
				if item == nil {
					item = &key
				} else {
					*item += pb
				}
			case 2:
				if item == nil {
					item = &val
				} else {
					*item += pb
				}
			default:
				if item != nil {
					*item += pb
				}
			}
		}
		if key != 5 {
			panic("key")
		}
		if val != 14 {
			panic("val")
		}
	}
}

func BenchmarkParseGoto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		index := 0
		var key, val, set, pb byte
		var item *byte
	loop:
		if len(parseBb) <= index || set == 2 {
			goto check
		}
		pb = parseBb[index]
		index++
		switch pb {
		case 0:
			if item != nil {
				set++
			}
			item = nil
			goto loop
		case 1:
			if item == nil {
				item = &key
			} else {
				*item += pb
			}
			goto loop
		case 2:
			if item == nil {
				item = &val
			} else {
				*item += pb
			}
			goto loop
		default:
			if item != nil {
				*item += pb
			}
			goto loop
		}
	check:
		if key != 5 {
			panic("key")
		}
		if val != 14 {
			panic("val")
		}
	}
}
