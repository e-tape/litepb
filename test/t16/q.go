package main

import (
	"crypto/rand"
	"fmt"
	"time"
	"unsafe"
)

var (
	buf      = make([]byte, 2)
	bufAfter = make([]byte, 2)
)

func main() {
	var r time.Time
	var d time.Duration

	{
		// then we can call rand.Read.
		_, _ = rand.Read(buf)
		_, _ = rand.Read(bufAfter)
		fmt.Println("buf")
		fmt.Println(buf)
		fmt.Println((*Bytes)(unsafe.Pointer(&buf[0])))
		fmt.Println((*[5]byte)(unsafe.Pointer(&buf[0])))
		fmt.Println("bufAfter")
		fmt.Println(bufAfter)
		fmt.Println((*Bytes)(unsafe.Pointer(&bufAfter[0])))
		fmt.Println((*[5]byte)(unsafe.Pointer(&bufAfter[0])))
		fmt.Println("buf inc")
		r := uintptr(unsafe.Pointer(&buf[0]))
		r++
		fmt.Println((*Bytes)(unsafe.Pointer(r)))
		fmt.Println((*[5]byte)(unsafe.Pointer(r)))
		return
	}

	_ = data[11:]

	r = time.Now()
	u2 := (*Bytes)(unsafe.Pointer(&data[11]))
	d = time.Since(r)
	fmt.Println(d)

	r = time.Now()
	u1 := data[11:]
	d = time.Since(r)
	fmt.Println(d)

	fmt.Println(u1[0], u1[1], u1[2], u1[3])
	fmt.Println(u2.B1, u2.B2, u2.B3, u2.B4)
	fmt.Println(data, len(data))
}

var data = []byte{
	3,
	1<<7 | 3<<3, 2 << 4,
	1<<7 | 3<<3, 1<<7 | 3<<3, 2 << 4,
	1<<7 | 3<<3, 1<<7 | 3<<3, 1<<7 | 3<<3, 2 << 4,
	1<<7 | 3<<3, 1<<7 | 3<<3, 1<<7 | 3<<3, 1<<7 | 3<<3, 2 << 4,
}

type Bytes struct {
	B1 byte
	B2 byte
	B3 byte
	B4 byte
	B5 byte
}
type Q struct {
	Var1 int
	Var2 int
}

func (a *Q) DecodeCurrent(data []byte, ind int) {
	var index int
	for len(data) > index {
		switch ind {
		case 0:
			var value int
			if len(data) > index {
				value |= int(data[index] & 127)
				if data[index] < 128 {
					goto end1
				}
			}
			index++
			if len(data) > index {
				value |= int(data[index]&127) << 7
				if data[index] < 128 {
					goto end1
				}
			}
			index++
			if len(data) > index {
				value |= int(data[index]&127) << 14
				if data[index] < 128 {
					goto end1
				}
			}
			index++
			if len(data) > index {
				value |= int(data[index]&127) << 21
				if data[index] < 128 {
					goto end1
				}
			}
			index++
			if len(data) > index {
				value |= int(data[index]&127) << 28
				if data[index] < 128 {
					goto end1
				}
			}
			panic("---")
		end1:
			index++
			a.Var1 = value
		case 1:
			var value int
			if len(data) > index {
				value |= int(data[index] & 127)
				if data[index] < 128 {
					goto end2
				}
			}
			index++
			if len(data) > index {
				value |= int(data[index]&127) << 7
				if data[index] < 128 {
					goto end2
				}
			}
			index++
			if len(data) > index {
				value |= int(data[index]&127) << 14
				if data[index] < 128 {
					goto end2
				}
			}
			index++
			if len(data) > index {
				value |= int(data[index]&127) << 21
				if data[index] < 128 {
					goto end2
				}
			}
			index++
			if len(data) > index {
				value |= int(data[index]&127) << 28
				if data[index] < 128 {
					goto end2
				}
			}
			panic("---")
		end2:
			index++
			a.Var2 = value
		}
	}
}

func (a *Q) DecodeCurrentModify(data []byte, ind int) {
	var index int
	for len(data) > index {
		switch ind {
		case 0, 1:
			var value int
			if len(data) > index {
				value |= int(data[index] & 127)
				if data[index] < 128 {
					goto end2
				}
			}
			index++
			if len(data) > index {
				value |= int(data[index]&127) << 7
				if data[index] < 128 {
					goto end2
				}
			}
			index++
			if len(data) > index {
				value |= int(data[index]&127) << 14
				if data[index] < 128 {
					goto end2
				}
			}
			index++
			if len(data) > index {
				value |= int(data[index]&127) << 21
				if data[index] < 128 {
					goto end2
				}
			}
			index++
			if len(data) > index {
				value |= int(data[index]&127) << 28
				if data[index] < 128 {
					goto end2
				}
			}
			panic("---")
		end2:
			index++
			switch ind {
			case 0:
				a.Var1 = value
			case 1:
				a.Var2 = value
			}
		}
	}
}

func (a *Q) DecodeModify(data []byte, ind int) {
	var index int
	for len(data) > index {
		switch ind {
		case 0, 1:
			var value int
			switch {
			case len(data) > index && data[index] < 128:
				value = int(data[index] & 127)
				index += 1
			case len(data) > index+1 && data[index+1] < 128:
				value = int(data[index]&127) | int(data[index+1]&127)<<7
				index += 2
			case len(data) > index+2 && data[index+2] < 128:
				value = int(data[index]&127) | int(data[index+1]&127)<<7 | int(data[index+2]&127)<<14
				index += 3
			case len(data) > index+3 && data[index+3] < 128:
				value = int(data[index]&127) | int(data[index+1]&127)<<7 | int(data[index+2]&127)<<14 | int(data[index+3]&127)<<21
				index += 4
			case len(data) > index+4 && data[index+4] < 128:
				value = int(data[index]&127) | int(data[index+1]&127)<<7 | int(data[index+2]&127)<<14 | int(data[index+3]&127)<<21 | int(data[index+4]&127)<<28
				index += 5
			}
			switch ind {
			case 0:
				a.Var1 = value
			case 1:
				a.Var2 = value
			}
		}
	}
}

func (a *Q) DecodeModify21(data []byte, ind int) {
	var index int
	for len(data) > index {
		switch ind {
		case 0, 1:
			var value int
			bd := *(*Bytes)(unsafe.Pointer(&data[index]))
			switch {
			case len(data) > index && bd.B1 < 128:
				value = int(bd.B1 & 127)
				index += 1
			case len(data) > index+1 && bd.B2 < 128:
				value = int(bd.B1&127) | int(bd.B2&127)<<7
				index += 2
			case len(data) > index+2 && bd.B3 < 128:
				value = int(bd.B1&127) | int(bd.B2&127)<<7 | int(bd.B3&127)<<14
				index += 3
			case len(data) > index+3 && bd.B4 < 128:
				value = int(bd.B1&127) | int(bd.B2&127)<<7 | int(bd.B3&127)<<14 | int(bd.B4&127)<<21
				index += 4
			case len(data) > index+4 && bd.B5 < 128:
				value = int(bd.B1&127) | int(bd.B2&127)<<7 | int(bd.B3&127)<<14 | int(bd.B4&127)<<21 | int(bd.B5&127)<<28
				index += 5
			}
			switch ind {
			case 0:
				a.Var1 = value
			case 1:
				a.Var2 = value
			}
		}
	}
}
func (a *Q) DecodeModify2(data []byte, ind int) {
	var index int
	for len(data) > index {
		bd := (*Bytes)(unsafe.Pointer(&data[index]))
		switch ind {
		case 0:
			var value int
			if len(data) > index {
				value |= int(bd.B1 & 127)
				if bd.B1 < 128 {
					goto end1
				}
			}
			index++
			if len(data) > index {
				value |= int(bd.B2&127) << 7
				if bd.B2 < 128 {
					goto end1
				}
			}
			index++
			if len(data) > index {
				value |= int(bd.B3&127) << 14
				if bd.B3 < 128 {
					goto end1
				}
			}
			index++
			if len(data) > index {
				value |= int(bd.B4&127) << 21
				if bd.B4 < 128 {
					goto end1
				}
			}
			index++
			if len(data) > index {
				value |= int(bd.B5&127) << 28
				if bd.B5 < 128 {
					goto end1
				}
			}
			panic("---")
		end1:
			index++
			a.Var1 = value
		case 1:
			var value int
			if len(data) > index {
				value |= int(bd.B1 & 127)
				if bd.B1 < 128 {
					goto end2
				}
			}
			index++
			if len(data) > index {
				value |= int(bd.B2&127) << 7
				if bd.B2 < 128 {
					goto end2
				}
			}
			index++
			if len(data) > index {
				value |= int(bd.B3&127) << 14
				if bd.B3 < 128 {
					goto end2
				}
			}
			index++
			if len(data) > index {
				value |= int(bd.B4&127) << 21
				if bd.B4 < 128 {
					goto end2
				}
			}
			index++
			if len(data) > index {
				value |= int(bd.B5&127) << 28
				if bd.B5 < 128 {
					goto end2
				}
			}
			panic("---")
		end2:
			index++
			a.Var2 = value
		}
	}
}
func (a *Q) DecodeModify3(data []byte, ind int) {
	var index int
	for len(data) > index {
		bd := (*[5]byte)(unsafe.Pointer(&data[index]))
		switch ind {
		case 0:
			var value int
			if len(data) > index {
				value |= int(bd[0] & 127)
				if bd[0] < 128 {
					goto end1
				}
			}
			index++
			if len(data) > index {
				value |= int(bd[1]&127) << 7
				if bd[1] < 128 {
					goto end1
				}
			}
			index++
			if len(data) > index {
				value |= int(bd[2]&127) << 14
				if bd[2] < 128 {
					goto end1
				}
			}
			index++
			if len(data) > index {
				value |= int(bd[3]&127) << 21
				if bd[3] < 128 {
					goto end1
				}
			}
			index++
			if len(data) > index {
				value |= int(bd[4]&127) << 28
				if bd[4] < 128 {
					goto end1
				}
			}
			panic("---")
		end1:
			index++
			a.Var1 = value
		case 1:
			var value int
			if len(data) > index {
				value |= int(bd[0] & 127)
				if bd[0] < 128 {
					goto end2
				}
			}
			index++
			if len(data) > index {
				value |= int(bd[1]&127) << 7
				if bd[1] < 128 {
					goto end2
				}
			}
			index++
			if len(data) > index {
				value |= int(bd[2]&127) << 14
				if bd[2] < 128 {
					goto end2
				}
			}
			index++
			if len(data) > index {
				value |= int(bd[3]&127) << 21
				if bd[3] < 128 {
					goto end2
				}
			}
			index++
			if len(data) > index {
				value |= int(bd[4]&127) << 28
				if bd[4] < 128 {
					goto end2
				}
			}
			panic("---")
		end2:
			index++
			a.Var2 = value
		}
	}
}
