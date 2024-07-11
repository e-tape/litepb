package t17

import (
	"unsafe"
)

type Bytes struct {
	B0 byte
	B1 byte
	B2 byte
	B3 byte
	B4 byte
}

func GotoModify() {
	var index int
	var value uint64
	var current int
	//arrPointer := uintptr(unsafe.Pointer(&data[0]))
	for len(data) > index {
		//arr := (*[5]byte)(unsafe.Pointer(arrPointer + uintptr(index)))
		//arr := (*[5]byte)(unsafe.Pointer(&data[index]))
		b := (*Bytes)(unsafe.Pointer(&data[index]))
		value = uint64(b.B0 & 127)
		if b.B0 < 128 {
			goto end
		}
		index++
		if len(data) > index {
			value |= uint64(b.B1&127) << 7
			if b.B1 < 128 {
				goto end
			}
		}
		index++
		if len(data) > index {
			value |= uint64(b.B2&127) << 14
			if b.B2 < 128 {
				goto end
			}
		}
		index++
		if len(data) > index {
			value |= uint64(b.B3&127) << 21
			if b.B3 < 128 {
				goto end
			}
		}
		index++
		if len(data) > index {
			value |= uint64(b.B4&127) << 28
			if b.B4 < 128 {
				goto end
			}
		}
		panic("eof")
	end:
		index++
		if value != dataValue[current] {
			panic("eq")
		}
		current++
	}
}
