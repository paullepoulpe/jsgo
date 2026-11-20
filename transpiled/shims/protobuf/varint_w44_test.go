package protobuf

import "testing"

func putTag1VarintOld(dst []byte, pos int, tag byte, v uint64) int {
	if v < 0x80 {
		_ = dst[pos+1]
		dst[pos] = tag
		dst[pos+1] = byte(v)
		return pos + 2
	}
	dst[pos] = tag
	for pos++; v >= 0x80; pos++ {
		dst[pos] = byte(v) | 0x80
		v >>= 7
	}
	dst[pos] = byte(v)
	return pos + 1
}

func BenchmarkPutTag1Varint(b *testing.B) {
	dst := make([]byte, 8192)

	cases := []struct {
		name string
		gen  func(i int) uint64
	}{
		{"1Byte", func(i int) uint64 { return uint64(i%120 + 1) }},
		{"2Byte", func(i int) uint64 { return uint64(200 + (i*13)%16000) }},
		{"3Byte", func(i int) uint64 { return uint64(20000 + i*37) }},
		{"ChecksMix", func(i int) uint64 {
			switch i % 20 {
			case 0, 1, 2:
				return uint64(i%120 + 1)
			case 3, 4, 5, 6:
				return uint64(20000 + i*37)
			default:
				return uint64(200 + (i*13)%16000)
			}
		}},
	}

	for _, tc := range cases {
		vals := make([]uint64, 1024)
		for i := range vals {
			vals[i] = tc.gen(i)
		}
		b.Run(tc.name+"_Old", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				pos := 0
				for _, v := range vals {
					pos = putTag1VarintOld(dst, pos, 0x18, v)
				}
				if pos == 0 {
					b.Fatal(pos)
				}
			}
		})
		b.Run(tc.name+"_New", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				pos := 0
				for _, v := range vals {
					pos = putTag1Varint(dst, pos, 0x18, v)
				}
				if pos == 0 {
					b.Fatal(pos)
				}
			}
		})
	}
}
