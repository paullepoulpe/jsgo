package jutil

import (
	"slices"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

// Go-only: language-adapter — the double-byte charsets Shift_JIS,
// windows-31j, GBK, GB2312, Big5, EUC-KR and x-windows-949: ASCII (and, in
// Shift_JIS and windows-31j, half-width katakana) in one byte, everything
// else in a lead byte and a trail byte. The pair tables come from
// golang.org/x/text/encoding (LICENSES/BSD-3-Clause-golang-x-text.txt),
// whose encodings follow the WHATWG Encoding Standard: each JDK charset
// takes the part of the x/text table in its own code space, and the pairs
// where the JDK's table differs are listed as overrides. Those were found
// black-box, by decoding every pair and encoding every char through the
// JDK's public CharsetDecoder and CharsetEncoder API; the digest rows of
// tools/oracle/java/JutilCharsetVectors.java keep checking all of them.
//
// Observed behaviour (black-box, CharsetDecoder in REPORT mode), the same
// in all seven charsets:
//   - a byte that is neither a char nor the lead byte of a mapped pair is
//     malformed, length 1; so is a lead byte at the end of the input;
//   - a lead byte followed by a byte that does not complete a mapped pair
//     is malformed with length 1 when that byte can start a sequence (a
//     char or a lead byte), which is then decoded on its own, and
//     unmappable with length 2 (both bytes consumed) otherwise.
//
// Each char encodes to its single byte or to the first pair (in byte
// order) that decodes to it, except for the few chars of dbcsSpec.oneWay.

const (
	dbcsLead    = -1
	dbcsInvalid = -2
)

type dbcsCoder struct {
	// single is the char of each byte, or dbcsLead or dbcsInvalid.
	single [256]int32
	// pairs holds the char of each pair, indexed by lead byte - 0x80 and
	// trail byte (0: unmapped).
	pairs    [128][256]uint16
	fromChar charTable
}

// dbcsRun maps the pairs from..to (one lead byte, consecutive trail bytes)
// to consecutive chars starting at first, or unmaps them when first < 0.
type dbcsRun struct {
	from, to uint16
	first    rune
}

// dbcsSpec describes a double-byte charset.
type dbcsSpec struct {
	base encoding.Encoding
	// The JDK's code space: lead bytes leadLo..leadHi (and, for Shift_JIS,
	// leadLo2..leadHi2), trail bytes in trails.
	leadLo, leadHi, leadLo2, leadHi2 byte
	trails                           [][2]byte
	// kana: the bytes 0xA1 to 0xDF are the half-width katakana U+FF61 to
	// U+FF9F (Shift_JIS).
	kana bool
	// overrides are applied in order after the x/text table is read.
	overrides []dbcsRun
	// encLast are lead bytes whose pairs the encoder uses only for chars
	// that no other pair decodes to.
	encLast []byte
	// oneWay are chars that encode to bytes that do not decode to them
	// (or not only to them), with those bytes.
	oneWay []dbcsOneWay
}

type dbcsOneWay struct {
	c     uint16
	bytes uint32 // packed (packBytes)
}

func (s *dbcsSpec) isLead(b byte) bool {
	return b >= s.leadLo && b <= s.leadHi || s.leadLo2 != 0 && b >= s.leadLo2 && b <= s.leadHi2
}

func (s *dbcsSpec) isTrail(b byte) bool {
	for _, t := range s.trails {
		if b >= t[0] && b <= t[1] {
			return true
		}
	}
	return false
}

func (s *dbcsSpec) build() charsetCoder {
	c := &dbcsCoder{}
	d := s.base.NewDecoder()
	var dst [16]byte
	src := make([]byte, 2)
	for b1 := 0x80; b1 < 0x100; b1++ {
		if !s.isLead(byte(b1)) {
			continue
		}
		for b2 := 0; b2 < 0x100; b2++ {
			if !s.isTrail(byte(b2)) {
				continue
			}
			src[0], src[1] = byte(b1), byte(b2)
			d.Reset()
			n, used, err := d.Transform(dst[:], src, true)
			if err != nil || used != 2 {
				continue
			}
			r, size := utf8.DecodeRune(dst[:n])
			if size != n || r == utf8.RuneError || r >= 0x10000 {
				continue
			}
			c.pairs[b1-0x80][b2] = uint16(r)
		}
	}
	for _, o := range s.overrides {
		row := &c.pairs[o.from>>8-0x80]
		for t := o.from & 0xFF; t <= o.to&0xFF; t++ {
			if o.first < 0 {
				row[t] = 0
			} else {
				row[t] = uint16(o.first + rune(t-o.from&0xFF))
			}
		}
	}
	for b := 0; b < 0x100; b++ {
		switch {
		case b < 0x80:
			c.single[b] = int32(b)
		case s.kana && b >= 0xA1 && b <= 0xDF:
			c.single[b] = int32(0xFF61 + b - 0xA1)
		default:
			c.single[b] = dbcsInvalid
			for _, ch := range c.pairs[b-0x80] {
				if ch != 0 {
					c.single[b] = dbcsLead
					break
				}
			}
		}
	}
	for b := 0; b < 0x100; b++ {
		if v := c.single[b]; v >= 0 {
			c.fromChar.setIfAbsent(uint16(v), packBytes(byte(b)))
		}
	}
	for pass := 0; pass < 2; pass++ {
		for b1 := 0x80; b1 < 0x100; b1++ {
			if slices.Contains(s.encLast, byte(b1)) != (pass == 1) {
				continue
			}
			for b2, ch := range c.pairs[b1-0x80] {
				if ch != 0 {
					c.fromChar.setIfAbsent(ch, packBytes(byte(b1), byte(b2)))
				}
			}
		}
	}
	for _, o := range s.oneWay {
		c.fromChar.set(o.c, o.bytes)
	}
	return c
}

func (c *dbcsCoder) decode(b []byte, report bool) (string, codingError) {
	return decodeSteps(c, b, report)
}

func (c *dbcsCoder) step(out []uint16, b []byte, i int) ([]uint16, int, codingError) {
	b1 := b[i]
	switch v := c.single[b1]; v {
	case dbcsInvalid:
		return out, 0, malformed(1)
	case dbcsLead:
		if i+1 == len(b) {
			return out, 0, malformed(1)
		}
		b2 := b[i+1]
		if ch := c.pairs[b1-0x80][b2]; ch != 0 {
			return append(out, ch), 2, codingError{}
		}
		if c.single[b2] != dbcsInvalid {
			return out, 0, malformed(1)
		}
		return out, 0, unmappable(2)
	default:
		return append(out, uint16(v)), 1, codingError{}
	}
}

func (c *dbcsCoder) canEncode(ch uint16) bool { return c.fromChar.get(ch) != 0 }

func (c *dbcsCoder) newEncoder() charsetEncoder { return tableEncoder{&c.fromChar} }

func newShiftJIS() charsetCoder {
	s := &dbcsSpec{
		base: japanese.ShiftJIS, leadLo: 0x81, leadHi: 0x9F, leadLo2: 0xE0, leadHi2: 0xEF,
		trails: [][2]byte{{0x40, 0x7E}, {0x80, 0xFC}}, kana: true,
		overrides: []dbcsRun{
			// JIS X 0208 mappings where x/text follows Windows-31J.
			{0x815C, 0x815C, 0x2014}, {0x8160, 0x8160, 0x301C}, {0x8161, 0x8161, 0x2016},
			{0x817C, 0x817C, 0x2212}, {0x8191, 0x8192, 0x00A2}, {0x81CA, 0x81CA, 0x00AC},
			// The Windows-31J extensions (NEC row 13, the NEC-selected IBM
			// extensions) are not in the JDK's Shift_JIS.
			{0x8740, 0x877E, -1}, {0x8780, 0x87FC, -1},
			{0xED40, 0xED7E, -1}, {0xED80, 0xEDFC, -1}, {0xEE40, 0xEE7E, -1}, {0xEE80, 0xEEFC, -1},
		},
		// The yen sign and overline encode to the JIS X 0201 Roman bytes.
		oneWay: []dbcsOneWay{{0x00A5, packBytes(0x5C)}, {0x203E, packBytes(0x7E)}},
	}
	return s.build()
}

func newWindows31J() charsetCoder {
	// The user-defined area F040 to F9FC maps to the Private Use Area,
	// 188 chars per lead byte (x/text leaves it unmapped).
	var pua []dbcsRun
	for b1 := 0xF0; b1 <= 0xF9; b1++ {
		first := rune(0xE000 + (b1-0xF0)*188)
		pua = append(pua, dbcsRun{uint16(b1<<8 | 0x40), uint16(b1<<8 | 0x7E), first},
			dbcsRun{uint16(b1<<8 | 0x80), uint16(b1<<8 | 0xFC), first + 63})
	}
	s := &dbcsSpec{
		base: japanese.ShiftJIS, leadLo: 0x81, leadHi: 0x9F, leadLo2: 0xE0, leadHi2: 0xFC,
		trails: [][2]byte{{0x40, 0x7E}, {0x80, 0xFC}}, kana: true,
		overrides: pua,
		// The NEC-selected IBM extensions duplicate the IBM extensions
		// (FA40 to FC4B), which the encoder prefers.
		encLast: []byte{0xED, 0xEE},
		// Best-fit encodings of Latin-1 and other chars.
		oneWay: []dbcsOneWay{
			{0x00A2, packBytes(0x81, 0x91)}, {0x00A3, packBytes(0x81, 0x92)}, {0x00A5, packBytes(0x5C)},
			{0x00AB, packBytes(0x81, 0xE1)}, {0x00AC, packBytes(0x81, 0xCA)}, {0x00AF, packBytes(0x81, 0x50)},
			{0x00B5, packBytes(0x83, 0xCA)}, {0x00B7, packBytes(0x81, 0x45)}, {0x00B8, packBytes(0x81, 0x43)},
			{0x00BB, packBytes(0x81, 0xE2)}, {0x203E, packBytes(0x7E)}, {0x3094, packBytes(0x83, 0x94)},
		},
	}
	return s.build()
}

func newGBK() charsetCoder {
	s := &dbcsSpec{
		base: simplifiedchinese.GBK, leadLo: 0x81, leadHi: 0xFE,
		trails: [][2]byte{{0x40, 0x7E}, {0x80, 0xFE}},
		overrides: append(append([]dbcsRun{}, gbkCommonOverrides...),
			// The JDK keeps in the Private Use Area some codes that x/text
			// maps to chars GB 18030-2005 assigned them, and maps A892 to
			// U+2641 (x/text: U+2295).
			dbcsRun{0xA6D9, 0xA6DF, 0xE78D}, dbcsRun{0xA6EC, 0xA6ED, 0xE794}, dbcsRun{0xA6F3, 0xA6F3, 0xE796},
			dbcsRun{0xA892, 0xA892, 0x2641}, dbcsRun{0xA8BC, 0xA8BC, 0xE7C7}, dbcsRun{0xA8BF, 0xA8BF, 0xE7C8},
			dbcsRun{0xA989, 0xA995, 0xE7E7}, dbcsRun{0xFE50, 0xFE7E, 0xE815}, dbcsRun{0xFE80, 0xFEA0, 0xE844}),
	}
	return s.build()
}

func newGB2312() charsetCoder {
	s := &dbcsSpec{
		base: simplifiedchinese.GBK, leadLo: 0xA1, leadHi: 0xF7,
		trails: [][2]byte{{0xA1, 0xFE}},
		overrides: []dbcsRun{
			// x/text: U+00B7 and U+2014.
			{0xA1A4, 0xA1A4, 0x30FB}, {0xA1AA, 0xA1AA, 0x2015},
			// Chars that GBK added inside the GB 2312 code space.
			{0xA2A1, 0xA2AA, -1}, {0xA2E3, 0xA2E3, -1}, {0xA6E0, 0xA6EB, -1},
			{0xA6EE, 0xA6F2, -1}, {0xA6F4, 0xA6F5, -1}, {0xA8BB, 0xA8BB, -1},
			{0xA8BD, 0xA8C0, -1},
		},
	}
	return s.build()
}

func newBig5() charsetCoder {
	s := &dbcsSpec{
		base: traditionalchinese.Big5, leadLo: 0xA1, leadHi: 0xF9,
		trails: [][2]byte{{0x40, 0x7E}, {0xA1, 0xFE}},
		overrides: []dbcsRun{
			// Symbols mapped differently (x/text follows Big5-HKSCS).
			{0xA145, 0xA145, 0x2022}, {0xA14E, 0xA14E, 0xFF64}, {0xA15A, 0xA15A, 0xFF3F},
			{0xA1C2, 0xA1C2, 0x203E}, {0xA1E3, 0xA1E3, 0x223C}, {0xA1F2, 0xA1F2, 0x2641},
			{0xA1F3, 0xA1F3, 0x2609}, {0xA1FE, 0xA1FE, 0x2571}, {0xA240, 0xA240, 0x2572},
			{0xA241, 0xA241, 0xFF0F}, {0xA242, 0xA242, 0xFF3C}, {0xA244, 0xA244, 0x00A5},
			{0xA246, 0xA247, 0x00A2},
			// Codes the JDK leaves unmapped.
			{0xA1C3, 0xA1C3, -1}, {0xA1C5, 0xA1C5, -1}, {0xA3C0, 0xA3FE, -1},
			{0xC7FD, 0xC7FE, -1}, {0xC840, 0xC87E, -1}, {0xC8A1, 0xC8FE, -1}, {0xF9D6, 0xF9FE, -1},
			// The JDK's kana, Cyrillic and numbered symbols at C6A1 to C7FC
			// (x/text has the Big5-HKSCS layout there).
			{0xC6A1, 0xC6A1, 0x30FE}, {0xC6A2, 0xC6A3, 0x309D}, {0xC6A4, 0xC6A4, 0x3005},
			{0xC6A5, 0xC6F7, 0x3041}, {0xC6F8, 0xC6FE, 0x30A1}, {0xC740, 0xC77E, 0x30A8},
			{0xC7A1, 0xC7B0, 0x30E7}, {0xC7B1, 0xC7B2, 0x0414}, {0xC7B3, 0xC7B3, 0x0401},
			{0xC7B4, 0xC7BA, 0x0416}, {0xC7BB, 0xC7CD, 0x0423}, {0xC7CE, 0xC7CE, 0x0451},
			{0xC7CF, 0xC7E8, 0x0436}, {0xC7E9, 0xC7F2, 0x2460}, {0xC7F3, 0xC7FC, 0x2474},
		},
		// Five chars that two pairs decode to encode to the second one.
		oneWay: []dbcsOneWay{
			{0x2571, packBytes(0xA2, 0xAC)}, {0x2572, packBytes(0xA2, 0xAD)}, {0x5341, packBytes(0xA4, 0x51)},
			{0x5345, packBytes(0xA4, 0xCA)}, {0xFF3F, packBytes(0xA1, 0xC4)},
		},
	}
	return s.build()
}

func newEUCKR() charsetCoder {
	s := &dbcsSpec{
		base: korean.EUCKR, leadLo: 0xA1, leadHi: 0xFE,
		trails: [][2]byte{{0xA1, 0xFE}},
		// KS X 1001:2002 added the postal mark at A2E8 (unmapped in x/text).
		overrides: []dbcsRun{{0xA2E8, 0xA2E8, 0x327E}},
	}
	return s.build()
}

func newWindows949() charsetCoder {
	s := &dbcsSpec{
		base: korean.EUCKR, leadLo: 0x81, leadHi: 0xFE,
		trails: [][2]byte{{0x41, 0x5A}, {0x61, 0x7A}, {0x81, 0xFE}},
		// The user-defined rows C9 and FE map to the Private Use Area
		// (x/text leaves them unmapped).
		overrides: []dbcsRun{{0xC9A1, 0xC9FE, 0xE000}, {0xFEA1, 0xFEFE, 0xE05E}},
	}
	return s.build()
}
