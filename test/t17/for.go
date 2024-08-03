package t17

func For() {
	var index int
	var value uint64
	var current int
	for len(data) > index {
		value = 0
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				panic("shift")
			}
			if index >= len(data) {
				panic("eof")
			}
			b := data[index]
			index++
			value |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		if value != dataValue[current] {
			panic("eq")
		}
		current++
	}
}
