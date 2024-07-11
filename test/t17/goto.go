package t17

func Goto() {
	var index int
	var value uint64
	var current int
	for len(data) > index {
		value = uint64(data[index] & 127)
		if data[index] < 128 {
			goto end
		}
		index++
		if len(data) > index {
			value |= uint64(data[index]&127) << 7
			if data[index] < 128 {
				goto end
			}
		}
		index++
		if len(data) > index {
			value |= uint64(data[index]&127) << 14
			if data[index] < 128 {
				goto end
			}
		}
		index++
		if len(data) > index {
			value |= uint64(data[index]&127) << 21
			if data[index] < 128 {
				goto end
			}
		}
		index++
		if len(data) > index {
			value |= uint64(data[index]&127) << 28
			if data[index] < 128 {
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
