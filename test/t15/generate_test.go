package t15

import (
	"fmt"
	"go/format"
	"os"
	"strings"
	"testing"
)

var (
	fieldNum = [][]int{
		{1, 20},
		//{30, 50},
		//{120, 150},
		//{300, 320},
		//{1150, 1200},
		//{1950, 2000},
		//{59950, 60000},
		//{24999990, 25000000},
		//{25000000, 25000000},
	}
	countWireType = 2
)

func getFieldNum(fieldNumber int, wireType int) []byte {
	fieldNumber <<= 3
	fieldNumber |= wireType
	var result []byte
	for fieldNumber >= 1<<7 {
		result = append(result, byte(fieldNumber&127|128))
		fieldNumber >>= 7
	}
	result = append(result, byte(fieldNumber))
	return result
}

func TestGenerateData(t *testing.T) {
	data := `package t15

import (
	"slices"
)

var data = slices.Concat(
	{ITEM}
)
`
	for _, r := range fieldNum {
		for i := r[0]; i <= r[1]; i++ {
			for j := 1; j < countWireType+1; j++ {
				val := getFieldNum(i, j)

				var valInt any
				switch len(val) {
				case 1:
					valInt = val[0]
				case 2:
					valInt = uint16(val[0]) | uint16(val[1])<<8
				case 3:
					valInt = uint32(val[0]) | uint32(val[1])<<8 | uint32(val[2])<<16
				case 4:
					valInt = uint32(val[0]) | uint32(val[1])<<8 | uint32(val[2])<<16 | uint32(val[3])<<24
				}
				data = strings.Replace(
					data,
					"{ITEM}",
					fmt.Sprintf(
						`getFieldNum(%d, %d), // %d - %d
	{ITEM}`, i, j, valInt, val,
					),
					1,
				)
			}
		}
	}
	data = strings.Replace(data, "{ITEM}", "", 1)
	if err := os.WriteFile(`q_data_test.go`, []byte(data), 0666); err != nil {
		panic(err)
	}
}

func TestGenerateDecode(t *testing.T) {
	bench := `package t15

import (
	"testing"
)

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var c int
		var index int
		for len(data) > index {
			fieldNum := int32(data[index] & 127)
			if data[index] < 128 {
				goto end
			}
			index++
			if len(data) > index {
				fieldNum |= int32(data[index]&127) << 7
				if data[index] < 128 {
					goto end
				}
			}
			index++
			if len(data) > index {
				fieldNum |= int32(data[index]&127) << 14
				if data[index] < 128 {
					goto end
				}
			}
			index++
			if len(data) > index {
				fieldNum |= int32(data[index]&127) << 21
				if data[index] < 128 {
					goto end
				}
			}
			panic("litepb: unexpected end of field number")
		end:
			index++
			switch fieldNum >> 3 {
			{FIELD_NUMBER}
			default:
				panic("litepb: unknown field number")
			}
		}
	}
}
`
	for _, r := range fieldNum {
		for i := r[0]; i <= r[1]; i++ {
			bench = strings.Replace(
				bench,
				"{FIELD_NUMBER}",
				fmt.Sprintf(
					`case %d:
				switch fieldNum & 7 {
				{WIRE_TYPE}
				default:
					panic("litepb: incorrect wire type")
				}
				{FIELD_NUMBER}`, i,
				),
				1,
			)
			for j := 1; j < countWireType+1; j++ {
				bench = strings.Replace(
					bench,
					"{WIRE_TYPE}",
					fmt.Sprintf(
						`case %d:
					c++
					{WIRE_TYPE}`, j,
					),
					1,
				)
			}
			bench = strings.Replace(bench, "{WIRE_TYPE}", "", 1)
		}
	}
	bench = strings.Replace(bench, "{FIELD_NUMBER}", "", 1)
	source, _ := format.Source([]byte(bench))
	if err := os.WriteFile(`q_decode_test.go`, source, 0666); err != nil {
		panic(err)
	}
}

func TestGenerateDecodeIf(t *testing.T) {
	bench := `package t15

import (
	"testing"
)

func BenchmarkDecodeIf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var c int
		var index int
		for len(data) > index {
			fieldNum := int32(data[index] & 127)
			if data[index] > 127 {
				index++
				if len(data) > index {
					fieldNum |= int32(data[index]&127) << 7
					if data[index] > 127 {
						index++
						if len(data) > index {
							fieldNum |= int32(data[index]&127) << 14
							if data[index] > 127 {
								index++
								if len(data) > index {
									fieldNum |= int32(data[index]&127) << 21
									if data[index] > 127 {
										panic("litepb: unexpected end of field number")
									}
								}
							}
						}
					}
				}
			}
			index++
			switch fieldNum >> 3 {
			{FIELD_NUMBER}
			default:
				panic("litepb: unknown field number")
			}
		}
	}
}
`
	for _, r := range fieldNum {
		for i := r[0]; i <= r[1]; i++ {
			bench = strings.Replace(
				bench,
				"{FIELD_NUMBER}",
				fmt.Sprintf(
					`case %d:
				switch fieldNum & 7 {
				{WIRE_TYPE}
				default:
					panic("litepb: incorrect wire type")
				}
				{FIELD_NUMBER}`, i,
				),
				1,
			)
			for j := 1; j < countWireType+1; j++ {
				bench = strings.Replace(
					bench,
					"{WIRE_TYPE}",
					fmt.Sprintf(
						`case %d:
					c++
					{WIRE_TYPE}`, j,
					),
					1,
				)
			}
			bench = strings.Replace(bench, "{WIRE_TYPE}", "", 1)
		}
	}
	bench = strings.Replace(bench, "{FIELD_NUMBER}", "", 1)
	source, _ := format.Source([]byte(bench))
	if err := os.WriteFile(`q_decode_if_test.go`, source, 0666); err != nil {
		panic(err)
	}
}

func TestGenerateDecodeStatic(t *testing.T) {
	bench := `package t15

import (
	"testing"
)

func BenchmarkDecodeStatic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var c int
		var index int
		for len(data) > index {{SWITCH}
			panic("litepb: unknown field number or incorrect wire type")
		}
	}
}
`
	caseList := make([]string, 4)
	for _, r := range fieldNum {
		for i := r[0]; i <= r[1]; i++ {
			for j := 1; j < countWireType+1; j++ {
				item := getFieldNum(i, j)
				var itemText []string
				for _, itemOne := range item {
					itemText = append(itemText, fmt.Sprintf("%d", itemOne))
				}
				var val any
				switch len(item) {
				case 1:
					val = item[0]
				case 2:
					val = uint16(item[0]) | uint16(item[1])<<8
				case 3:
					val = uint32(item[0]) | uint32(item[1])<<8 | uint32(item[2])<<16
				case 4:
					val = uint32(item[0]) | uint32(item[1])<<8 | uint32(item[2])<<16 | uint32(item[3])<<24
				}
				caseList[len(item)-1] += fmt.Sprintf(
					`
				case %d: // %[3]d,%[4]d
					c++`,
					val,
					len(item),
					i, j,
				)
			}
		}
	}
	switchCode := ""
	for i := 0; i < len(caseList); i++ {
		if caseList[i] == "" {
			continue
		}
		var val string
		switch i {
		case 0:
			val = "data[index]"
		case 1:
			val = "uint16(data[index]) | uint16(data[index+1])<<8"
		case 2:
			val = "uint32(data[index]) | uint32(data[index+1])<<8 | uint32(data[index+2])<<16"
		case 3:
			val = "uint32(data[index]) | uint32(data[index+1])<<8 | uint32(data[index+2])<<16 | uint32(data[index+3])<<24"
		}
		switchCode += fmt.Sprintf(`
			if len(data) > index+%d && data[index+%[1]d] < 128 {
				switch %s {%s
				}
				index += %d
				continue
			}`, i, val, caseList[i], i+1,
		)
	}
	bench = strings.Replace(bench, "{SWITCH}", switchCode, 1)
	if err := os.WriteFile(`q_decode_static_test.go`, []byte(bench), 0666); err != nil {
		panic(err)
	}
}
