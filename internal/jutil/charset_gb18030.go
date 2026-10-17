package jutil

import (
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// Go-only: language-adapter — GB18030, which the JDK (21) implements as
// GB 18030-2022. One-byte ASCII; two-byte codes (lead 0x81 to 0xFE, trail
// 0x40 to 0x7E or 0x80 to 0xFE), all mapped; four-byte codes (bytes
// 0x81-0xFE, 0x30-0x39, 0x81-0xFE, 0x30-0x39) numbered linearly, the BMP
// chars that have no two-byte code from 0x81308130 and U+10000 to U+10FFFF
// from 0x90308130, as the standard defines. The two-byte table and the
// order of the four-byte BMP chars come from
// golang.org/x/text/encoding/simplifiedchinese
// (LICENSES/BSD-3-Clause-golang-x-text.txt), which follows GB 18030-2005;
// black-box comparison with the JDK (as in charset_dbcs.go) found the
// differences listed in gb18030Overrides: the Private Use Area mappings of
// the user-defined areas and of the codes the standard left free, and the
// 2022 changes, where a char moved from a four-byte code to a two-byte one
// whose PUA char took the four-byte code in exchange.
//
// Observed behaviour of the JDK's decoder (black-box, REPORT mode):
//   - 0x80 and 0xFF are malformed, length 1, and so is a lead byte at the
//     end of the input or followed by a byte below 0x30;
//   - a lead byte followed by 0x3A..0x3F, 0x7F or 0xFF is malformed,
//     length 2;
//   - a lead byte and a digit start a four-byte code: if fewer than two
//     bytes follow, the rest of the input is malformed; a third byte
//     outside 0x81..0xFE is malformed with length 3, a fourth byte outside
//     0x30..0x39 with length 4, and a four-byte code after 0x8431A439
//     (U+FFFF) and before 0x90308130, or after 0xE3329A35 (U+10FFFF), is
//     malformed, length 4.

// gbkCommonOverrides are the Private Use Area mappings that GBK and
// GB18030 share in the JDK.
var gbkCommonOverrides = []dbcsRun{
	{0xA140, 0xA17E, 0xE4C6}, {0xA180, 0xA1A0, 0xE505}, {0xA240, 0xA27E, 0xE526},
	{0xA280, 0xA2A0, 0xE565}, {0xA2AB, 0xA2B0, 0xE766}, {0xA2E4, 0xA2E4, 0xE76D},
	{0xA2EF, 0xA2F0, 0xE76E}, {0xA2FD, 0xA2FE, 0xE770}, {0xA340, 0xA37E, 0xE586},
	{0xA380, 0xA3A0, 0xE5C5}, {0xA440, 0xA47E, 0xE5E6}, {0xA480, 0xA4A0, 0xE625},
	{0xA4F4, 0xA4FE, 0xE772}, {0xA540, 0xA57E, 0xE646}, {0xA580, 0xA5A0, 0xE685},
	{0xA5F7, 0xA5FE, 0xE77D}, {0xA640, 0xA67E, 0xE6A6}, {0xA680, 0xA6A0, 0xE6E5},
	{0xA6B9, 0xA6C0, 0xE785}, {0xA6F6, 0xA6FE, 0xE797}, {0xA740, 0xA77E, 0xE706},
	{0xA780, 0xA7A0, 0xE745}, {0xA7C2, 0xA7D0, 0xE7A0}, {0xA7F2, 0xA7FE, 0xE7AF},
	{0xA896, 0xA8A0, 0xE7BC}, {0xA8C1, 0xA8C4, 0xE7C9}, {0xA8EA, 0xA8FE, 0xE7CD},
	{0xA958, 0xA958, 0xE7E2}, {0xA95B, 0xA95B, 0xE7E3}, {0xA95D, 0xA95F, 0xE7E4},
	{0xA997, 0xA9A3, 0xE7F4}, {0xA9F0, 0xA9FE, 0xE801}, {0xAAA1, 0xAAFE, 0xE000},
	{0xABA1, 0xABFE, 0xE05E}, {0xACA1, 0xACFE, 0xE0BC}, {0xADA1, 0xADFE, 0xE11A},
	{0xAEA1, 0xAEFE, 0xE178}, {0xAFA1, 0xAFFE, 0xE1D6}, {0xD7FA, 0xD7FE, 0xE810},
	{0xF8A1, 0xF8FE, 0xE234}, {0xF9A1, 0xF9FE, 0xE292}, {0xFAA1, 0xFAFE, 0xE2F0},
	{0xFBA1, 0xFBFE, 0xE34E}, {0xFCA1, 0xFCFE, 0xE3AC}, {0xFDA1, 0xFDFE, 0xE40A},
	{0xFEA1, 0xFEFE, 0xE468},
}

// gb18030Overrides are the JDK's GB18030 two-byte codes that differ from
// x/text beyond gbkCommonOverrides.
var gb18030Overrides = []dbcsRun{
	// Vertical forms (x/text: unmapped).
	{0xA6D9, 0xA6D9, 0xFE10}, {0xA6DA, 0xA6DA, 0xFE12}, {0xA6DB, 0xA6DB, 0xFE11},
	{0xA6DC, 0xA6DF, 0xFE13}, {0xA6EC, 0xA6ED, 0xFE17}, {0xA6F3, 0xA6F3, 0xFE19},
	{0xA8BC, 0xA8BC, 0x1E3F},
	// Radicals left in the PUA, and the chars GB 18030-2022 moved here.
	{0xFE51, 0xFE53, 0xE816}, {0xFE59, 0xFE59, 0x9FB4}, {0xFE61, 0xFE61, 0x9FB5},
	{0xFE66, 0xFE67, 0x9FB6}, {0xFE6C, 0xFE6C, 0xE831}, {0xFE6D, 0xFE6D, 0x9FB8},
	{0xFE76, 0xFE76, 0xE83B}, {0xFE7E, 0xFE7E, 0x9FB9}, {0xFE90, 0xFE90, 0x9FBA},
	{0xFE91, 0xFE91, 0xE855}, {0xFEA0, 0xFEA0, 0x9FBB},
}

// gb18030FourByteOverrides are the four-byte BMP codes where the JDK
// differs from x/text.
var gb18030FourByteOverrides = []struct {
	code uint32
	c    uint16
}{
	{0x8135F437, 0xE7C7},
	{0x82359037, 0xE81E}, {0x82359038, 0xE826}, {0x82359039, 0xE82B}, {0x82359130, 0xE82C},
	{0x82359131, 0xE832}, {0x82359132, 0xE843}, {0x82359133, 0xE854}, {0x82359134, 0xE864},
	{0x84318236, 0xE78D}, {0x84318237, 0xE78F}, {0x84318238, 0xE78E}, {0x84318239, 0xE790},
	{0x84318330, 0xE791}, {0x84318331, 0xE792}, {0x84318332, 0xE793}, {0x84318333, 0xE794},
	{0x84318334, 0xE795}, {0x84318335, 0xE796},
	// U+FFFD (x/text decodes it but the table reader cannot tell it from
	// an error).
	{0x8431A437, 0xFFFD},
}

// gb18030Linear is the linear index of a four-byte code.
func gb18030Linear(b1, b2, b3, b4 byte) int {
	return ((int(b1)-0x81)*10+int(b2)-0x30)*1260 + (int(b3)-0x81)*10 + int(b4) - 0x30
}

const (
	gb18030BMPCount     = 39420  // 0x81308130 to 0x8431A439
	gb18030SupplStart   = 189000 // linear index of 0x90308130
	gb18030EncFourBytes = 0x80000000
)

type gb18030Coder struct {
	two *dbcsCoder
	// bmp holds the char of each four-byte BMP code.
	bmp [gb18030BMPCount]uint16
	// fromChar packs two-byte codes (packBytes) and four-byte codes
	// (gb18030EncFourBytes | linear index).
	fromChar charTable
}

func newGB18030() charsetCoder {
	s := &dbcsSpec{
		base: simplifiedchinese.GB18030, leadLo: 0x81, leadHi: 0xFE,
		trails:    [][2]byte{{0x40, 0x7E}, {0x80, 0xFE}},
		overrides: append(append([]dbcsRun{}, gbkCommonOverrides...), gb18030Overrides...),
	}
	c := &gb18030Coder{two: s.build().(*dbcsCoder)}
	d := simplifiedchinese.GB18030.NewDecoder()
	var dst [16]byte
	src := make([]byte, 4)
	for i := 0; i < gb18030BMPCount; i++ {
		src[0], src[1], src[2], src[3] = byte(0x81+i/12600), byte(0x30+i/1260%10), byte(0x81+i/10%126), byte(0x30+i%10)
		d.Reset()
		n, used, err := d.Transform(dst[:], src, true)
		if err != nil || used != 4 {
			continue
		}
		if r, size := utf8.DecodeRune(dst[:n]); size == n && r != utf8.RuneError && r < 0x10000 {
			c.bmp[i] = uint16(r)
		}
	}
	for _, o := range gb18030FourByteOverrides {
		c.bmp[gb18030Linear(byte(o.code>>24), byte(o.code>>16), byte(o.code>>8), byte(o.code))] = o.c
	}
	// The two-byte coder's table, extended in place (the two-byte coder
	// is only used for its pairs).
	c.fromChar = c.two.fromChar
	for i, ch := range c.bmp {
		if ch != 0 {
			c.fromChar.setIfAbsent(ch, gb18030EncFourBytes|uint32(i))
		}
	}
	return c
}

func (c *gb18030Coder) decode(b []byte, report bool) (string, codingError) {
	return decodeSteps(c, b, report)
}

func (c *gb18030Coder) step(out []uint16, b []byte, i int) ([]uint16, int, codingError) {
	b1 := b[i]
	rest := len(b) - i
	switch {
	case b1 < 0x80:
		return append(out, uint16(b1)), 1, codingError{}
	case b1 == 0x80 || b1 == 0xFF || rest < 2:
		return out, 0, malformed(1)
	}
	b2 := b[i+1]
	switch {
	case b2 < 0x30:
		return out, 0, malformed(1)
	case b2 > 0x39:
		if ch := c.two.pairs[b1-0x80][b2]; ch != 0 {
			return append(out, ch), 2, codingError{}
		}
		return out, 0, malformed(2)
	case rest < 4:
		return out, 0, malformed(rest)
	}
	b3, b4 := b[i+2], b[i+3]
	switch {
	case b3 < 0x81 || b3 == 0xFF:
		return out, 0, malformed(3)
	case b4 < 0x30 || b4 > 0x39:
		return out, 0, malformed(4)
	}
	n := gb18030Linear(b1, b2, b3, b4)
	switch {
	case n < gb18030BMPCount:
		if ch := c.bmp[n]; ch != 0 {
			return append(out, ch), 4, codingError{}
		}
	case n >= gb18030SupplStart && n-gb18030SupplStart < 0x100000:
		v := n - gb18030SupplStart
		return append(out, uint16(0xD800+v>>10), uint16(0xDC00+v&0x3FF)), 4, codingError{}
	}
	return out, 0, malformed(4)
}

func (c *gb18030Coder) canEncode(ch uint16) bool { return c.fromChar.get(ch) != 0 }

func (c *gb18030Coder) newEncoder() charsetEncoder { return gb18030Encoder{c} }

type gb18030Encoder struct{ c *gb18030Coder }

func appendGB18030FourBytes(dst []byte, n int) []byte {
	return append(dst, byte(0x81+n/12600), byte(0x30+n/1260%10), byte(0x81+n/10%126), byte(0x30+n%10))
}

func (e gb18030Encoder) encodeRune(dst []byte, r rune) []byte {
	if r >= 0x10000 {
		return appendGB18030FourBytes(dst, gb18030SupplStart+int(r-0x10000))
	}
	v := e.c.fromChar.get(uint16(r))
	switch {
	case v == 0:
		return append(dst, '?')
	case v&gb18030EncFourBytes != 0:
		return appendGB18030FourBytes(dst, int(v&^gb18030EncFourBytes))
	}
	return appendPacked(dst, v)
}

func (gb18030Encoder) encodeMalformed(dst []byte) []byte { return append(dst, '?') }
func (gb18030Encoder) finish(dst []byte) []byte          { return dst }
