package generator

import "testing"

func BenchmarkReset(b *testing.B) {
	type Test struct {
		Int1    int32
		Int2    int32
		String1 string
		String2 string
		Map1    map[string]string
		Map2    map[string]string
	}

	b.Run("struct-set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			t := &Test{
				Int1:    1,
				Int2:    2,
				String1: "1",
				String2: "2",
				Map1:    make(map[string]string),
				Map2:    make(map[string]string),
			}

			*t = Test{}
			_ = t
		}
	})

	b.Run("field-set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			t := &Test{
				Int1:    1,
				Int2:    2,
				String1: "1",
				String2: "2",
				Map1:    make(map[string]string),
				Map2:    make(map[string]string),
			}

			t.Int1 = 0
			t.Int2 = 0
			t.String1 = ""
			t.String2 = ""
			t.Map1 = nil
			t.Map2 = nil
			_ = t
		}
	})
}
