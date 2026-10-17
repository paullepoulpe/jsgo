package jutil

import (
	"sync"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
)

// Go-only: language-adapter — EUC-JP and ISO-2022-JP. The JIS X 0208 and
// JIS X 0212 tables come from golang.org/x/text/encoding/japanese
// (LICENSES/BSD-3-Clause-golang-x-text.txt). Black-box comparison with the
// JDK (as in charset_dbcs.go) found the same JIS X 0212 table and these
// differences in JIS X 0208, applied as overrides: the JDK maps seven
// symbols as JIS does where x/text follows Windows (eucJPOverrides), and
// leaves out the NEC row 13 and the IBM extensions in rows 89 to 92.
//
// Observed behaviour of the JDK's EUC-JP decoder (black-box, REPORT mode):
//   - 0x8E starts half-width katakana: 0x8E 0xA1..0xDF are U+FF61..U+FF9F;
//     0x8E with any other byte is unmappable, length 2;
//   - 0x8F starts a JIS X 0212 char in the next two bytes; any three bytes
//     that are not a mapped char are unmappable, length 3;
//   - any other byte from 0x80 starts a JIS X 0208 char; any two bytes that
//     are not a mapped char (whatever the second byte) are unmappable,
//     length 2;
//   - a sequence cut short by the end of the input is malformed, with the
//     length of what is left.
//
// The encoder prefers JIS X 0208 to JIS X 0212, and also encodes the yen
// sign and the overline to the ASCII bytes 0x5C and 0x7E (one way).

var eucJPOverrides = []dbcsRun{
	{0xA1BD, 0xA1BD, 0x2014}, {0xA1C1, 0xA1C1, 0x301C}, {0xA1C2, 0xA1C2, 0x2016},
	{0xA1DD, 0xA1DD, 0x2212}, {0xA1F1, 0xA1F2, 0x00A2}, {0xA2CC, 0xA2CC, 0x00AC},
	{0xADA1, 0xADFE, -1}, {0xF9A1, 0xF9FE, -1}, {0xFAA1, 0xFAFE, -1},
	{0xFBA1, 0xFBFE, -1}, {0xFCA1, 0xFCFE, -1},
}

// jisTables holds JIS X 0208 and JIS X 0212 as the JDK maps them, indexed
// by row and cell (0 to 93 each; 0: unmapped), and the encodings of every
// char they cover.
type jisTables struct {
	jis0208, jis0212 [94][94]uint16
	// fromChar packs the EUC-JP bytes of each char.
	fromChar charTable
}

var (
	jisTablesOnce sync.Once
	jisTablesData *jisTables
)

func loadJISTables() *jisTables {
	jisTablesOnce.Do(func() {
		t := &jisTables{}
		d := japanese.EUCJP.NewDecoder()
		var dst [16]byte
		decode := func(src []byte) uint16 {
			d.Reset()
			n, used, err := d.Transform(dst[:], src, true)
			if err != nil || used != len(src) {
				return 0
			}
			r, size := utf8.DecodeRune(dst[:n])
			if size != n || r == utf8.RuneError || r >= 0x10000 {
				return 0
			}
			return uint16(r)
		}
		for row := 0; row < 94; row++ {
			for cell := 0; cell < 94; cell++ {
				t.jis0208[row][cell] = decode([]byte{byte(0xA1 + row), byte(0xA1 + cell)})
				t.jis0212[row][cell] = decode([]byte{0x8F, byte(0xA1 + row), byte(0xA1 + cell)})
			}
		}
		for _, o := range eucJPOverrides {
			row := &t.jis0208[o.from>>8-0xA1]
			for b := o.from & 0xFF; b <= o.to&0xFF; b++ {
				if o.first < 0 {
					row[b-0xA1] = 0
				} else {
					row[b-0xA1] = uint16(o.first + rune(b-o.from&0xFF))
				}
			}
		}
		for b := 0; b < 0x80; b++ {
			t.fromChar.set(uint16(b), packBytes(byte(b)))
		}
		t.fromChar.set(0x00A5, packBytes(0x5C))
		t.fromChar.set(0x203E, packBytes(0x7E))
		for row := 0; row < 94; row++ {
			for cell := 0; cell < 94; cell++ {
				if c := t.jis0208[row][cell]; c != 0 {
					t.fromChar.setIfAbsent(c, packBytes(byte(0xA1+row), byte(0xA1+cell)))
				}
			}
		}
		for b := 0xA1; b <= 0xDF; b++ {
			t.fromChar.setIfAbsent(uint16(0xFF61+b-0xA1), packBytes(0x8E, byte(b)))
		}
		for row := 0; row < 94; row++ {
			for cell := 0; cell < 94; cell++ {
				if c := t.jis0212[row][cell]; c != 0 {
					t.fromChar.setIfAbsent(c, packBytes(0x8F, byte(0xA1+row), byte(0xA1+cell)))
				}
			}
		}
		jisTablesData = t
	})
	return jisTablesData
}

// jisLookup returns the char of a two-byte JIS code (both bytes 0x21 to
// 0x7E when GL, 0xA1 to 0xFE when GR), or 0.
func jisLookup(table *[94][94]uint16, b1, b2 byte) uint16 {
	b1, b2 = b1&0x7F, b2&0x7F
	if b1 < 0x21 || b1 > 0x7E || b2 < 0x21 || b2 > 0x7E {
		return 0
	}
	return table[b1-0x21][b2-0x21]
}

type eucJPCoder struct{ t *jisTables }

func newEUCJP() charsetCoder { return eucJPCoder{loadJISTables()} }

func (c eucJPCoder) decode(b []byte, report bool) (string, codingError) {
	return decodeSteps(c, b, report)
}

func (c eucJPCoder) step(out []uint16, b []byte, i int) ([]uint16, int, codingError) {
	b1 := b[i]
	rest := len(b) - i
	switch {
	case b1 < 0x80:
		return append(out, uint16(b1)), 1, codingError{}
	case b1 == 0x8F:
		if rest < 3 {
			return out, 0, malformed(rest)
		}
		if b[i+1] >= 0xA1 && b[i+2] >= 0xA1 {
			if ch := jisLookup(&c.t.jis0212, b[i+1], b[i+2]); ch != 0 {
				return append(out, ch), 3, codingError{}
			}
		}
		return out, 0, unmappable(3)
	case rest < 2:
		return out, 0, malformed(1)
	case b1 == 0x8E:
		if b2 := b[i+1]; b2 >= 0xA1 && b2 <= 0xDF {
			return append(out, uint16(0xFF61+int(b2)-0xA1)), 2, codingError{}
		}
		return out, 0, unmappable(2)
	default:
		if b1 >= 0xA1 && b[i+1] >= 0xA1 {
			if ch := jisLookup(&c.t.jis0208, b1, b[i+1]); ch != 0 {
				return append(out, ch), 2, codingError{}
			}
		}
		return out, 0, unmappable(2)
	}
}

func (c eucJPCoder) canEncode(ch uint16) bool { return c.t.fromChar.get(ch) != 0 }

func (c eucJPCoder) newEncoder() charsetEncoder { return tableEncoder{&c.t.fromChar} }
