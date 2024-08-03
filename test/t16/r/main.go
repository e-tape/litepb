package main

/*
go build -gcflags -S main.go 2> asm.txt
go build -gcflags "-d=ssa/check_bce/debug=1" main.go
*/

var data = []byte{
	3,
	1<<7 | 3<<3, 2 << 4,
	1<<7 | 3<<3, 1<<7 | 3<<3, 2 << 4,
	1<<7 | 3<<3, 1<<7 | 3<<3, 1<<7 | 3<<3, 2 << 4,
	1<<7 | 3<<3, 1<<7 | 3<<3, 1<<7 | 3<<3, 1<<7 | 3<<3, 2 << 4,
}

func main() {
	//q := &Q{}
	//r := time.Now()
	//q.DecodeModify2(data, 1)
	//fmt.Println(time.Since(r))
}

type Q struct {
	Var1   int
	Var164 int64
	Var2   int
}

//func (a *Q) DecodeModify2(data []byte, ind int) byte {
//	_ = data[ind]
//	if len(data) > ind+10 {
//		return data[ind]
//	}
//	return 0
//}
//
//func (a *Q) Decode(data []byte, ind int) {
//	for index := 0; index < len(data); index++ {
//		switch ind {
//		case 0:
//			var value int
//			if len(data) > index {
//				value |= int(data[index] & 127)
//				if data[index] < 128 {
//					goto end1
//				}
//			}
//			index++
//			if len(data) > index {
//				value |= int(data[index]&127) << 7
//				if data[index] < 128 {
//					goto end1
//				}
//			}
//			index++
//			if len(data) > index {
//				value |= int(data[index]&127) << 14
//				if data[index] < 128 {
//					goto end1
//				}
//			}
//			index++
//			if len(data) > index {
//				value |= int(data[index]&127) << 21
//				if data[index] < 128 {
//					goto end1
//				}
//			}
//			index++
//			if len(data) > index {
//				value |= int(data[index]&127) << 28
//				if data[index] < 128 {
//					goto end1
//				}
//			}
//			panic("---")
//		end1:
//			index++
//			a.Var1 = value
//		case 1:
//			var value int
//			if len(data) > index {
//				value |= int(data[index] & 127)
//				if data[index] < 128 {
//					goto end2
//				}
//			}
//			index++
//			if len(data) > index {
//				value |= int(data[index]&127) << 7
//				if data[index] < 128 {
//					goto end2
//				}
//			}
//			index++
//			if len(data) > index {
//				value |= int(data[index]&127) << 14
//				if data[index] < 128 {
//					goto end2
//				}
//			}
//			index++
//			if len(data) > index {
//				value |= int(data[index]&127) << 21
//				if data[index] < 128 {
//					goto end2
//				}
//			}
//			index++
//			if len(data) > index {
//				value |= int(data[index]&127) << 28
//				if data[index] < 128 {
//					goto end2
//				}
//			}
//			panic("---")
//		end2:
//			index++
//			a.Var2 = value
//		}
//	}
//}

func (a *Q) Decode2(data []byte, ind int) {
	for index := 0; index < len(data); index++ {
		switch {
		//case len(data) > index && data[index] < 128:
		//	a.Var164 = int64(data[index] & 127)
		//	index += 1
		//case len(data) > index+1 && data[index+1] < 128:
		//	a.Var164 = int64(data[index]&127) | int64(data[index+1]&127)<<7
		//	index += 2
		//case len(data) > index+2 && data[index+2] < 128:
		//	a.Var164 = int64(data[index]&127) | int64(data[index+1]&127)<<7 | int64(data[index+2]&127)<<14
		//	index += 3
		//case len(data) > index+3 && data[index+3] < 128:
		//	a.Var164 = int64(data[index]&127) | int64(data[index+1]&127)<<7 | int64(data[index+2]&127)<<14 | int64(data[index+3]&127)<<21
		//	index += 4
		//case len(data) > index+4 && data[index+4] < 128:
		//	a.Var164 = int64(data[index]&127) | int64(data[index+1]&127)<<7 | int64(data[index+2]&127)<<14 | int64(data[index+3]&127)<<21 | int64(data[index+4]&127)<<28
		//	index += 5
		//case len(data) > index+5 && data[index+5] < 128:
		//	a.Var164 = int64(data[index]&127) | int64(data[index+1]&127)<<7 | int64(data[index+2]&127)<<14 | int64(data[index+3]&127)<<21 | int64(data[index+4]&127)<<28 | int64(data[index+5]&127)<<35
		//	index += 6
		//case len(data) > index+6 && data[index+6] < 128:
		//	a.Var164 = int64(data[index]&127) | int64(data[index+1]&127)<<7 | int64(data[index+2]&127)<<14 | int64(data[index+3]&127)<<21 | int64(data[index+4]&127)<<28 | int64(data[index+5]&127)<<35 | int64(data[index+6]&127)<<42
		//	index += 7
		//case len(data) > index+7 && data[index+7] < 128:
		//	a.Var164 = int64(data[index]&127) | int64(data[index+1]&127)<<7 | int64(data[index+2]&127)<<14 | int64(data[index+3]&127)<<21 | int64(data[index+4]&127)<<28 | int64(data[index+5]&127)<<35 | int64(data[index+6]&127)<<42 | int64(data[index+7]&127)<<49
		//	index += 8
		case len(data) > index+8 && data[index+8] < 128:
			q := ([9]byte)(data[index : index+8])
			//a.Var164 = int64(data[index]&127) | int64(data[index+1]&127)<<7 | int64(data[index+2]&127)<<14 | int64(data[index+3]&127)<<21 | int64(data[index+4]&127)<<28 | int64(data[index+5]&127)<<35 | int64(data[index+6]&127)<<42 | int64(data[index+7]&127)<<49 | int64(data[index+8]&127)<<56
			a.Var164 = int64(q[0]&127) |
				int64(q[1]&127)<<7 |
				int64(q[2]&127)<<14 |
				int64(q[3]&127)<<21 |
				int64(q[4]&127)<<28 |
				int64(q[5]&127)<<35 |
				int64(q[6]&127)<<42 |
				int64(q[7]&127)<<49 |
				int64(q[8]&127)<<56
			index += 9
			//case len(data) > index+9 && data[index+9] < 128:
			//	a.Var164 = int64(data[index]&127) | int64(data[index+1]&127)<<7 | int64(data[index+2]&127)<<14 | int64(data[index+3]&127)<<21 | int64(data[index+4]&127)<<28 | int64(data[index+5]&127)<<35 | int64(data[index+6]&127)<<42 | int64(data[index+7]&127)<<49 | int64(data[index+8]&127)<<56 | int64(data[index+9]&127)<<63
			//	index += 10
		}
	}
}
