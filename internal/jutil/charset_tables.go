package jutil

import (
	"slices"
	"sync"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

// Go-only: language-adapter — the single-byte charsets. Their tables come
// from golang.org/x/text/encoding/charmap
// (LICENSES/BSD-3-Clause-golang-x-text.txt), with the differences from the
// JDK applied as overrides below. Those were found black-box, by decoding
// every byte and encoding every char through the JDK's public
// CharsetDecoder and CharsetEncoder API; the digest rows of
// tools/oracle/java/JutilCharsetVectors.java keep checking all of them.

// sbcsCoder is a single-byte charset: each byte decodes to one char, or is
// an erroneous sequence of length 1 on its own.
type sbcsCoder struct {
	toChar [256]uint16
	mapped [256]bool
	// badByte is what an unmapped byte reports. Observed behaviour
	// (black-box): unmappable, except in US-ASCII, where the bytes 0x80 to
	// 0xFF are malformed.
	badByte  codingError
	fromChar charTable
}

// sbcsOverride maps byte b to char c, or unmaps it when c < 0.
type sbcsOverride struct {
	b byte
	c rune
}

// newSBCS builds a single-byte charset from an x/text table. With
// isoControls, the bytes 0x80 to 0x9F decode to the C1 controls U+0080 to
// U+009F, as the JDK's ISO 8859 charsets do (x/text leaves them unmapped).
// Each char encodes to the byte that decodes to it (no JDK table here
// maps two bytes to one char, and none encodes a char one way).
func newSBCS(cm *charmap.Charmap, isoControls bool, overrides ...sbcsOverride) func() charsetCoder {
	return func() charsetCoder {
		c := &sbcsCoder{badByte: unmappable(1)}
		for b := 0; b < 256; b++ {
			r := cm.DecodeByte(byte(b))
			if isoControls && b >= 0x80 && b < 0xA0 {
				r = rune(b)
			}
			if r != utf8.RuneError {
				c.toChar[b], c.mapped[b] = uint16(r), true
			}
		}
		for _, o := range overrides {
			c.toChar[o.b], c.mapped[o.b] = uint16(o.c), o.c >= 0
		}
		c.buildEncoder()
		return c
	}
}

// newUSASCII is US-ASCII: the bytes 0x00 to 0x7F.
func newUSASCII() charsetCoder {
	c := &sbcsCoder{badByte: malformed(1)}
	for b := 0; b < 0x80; b++ {
		c.toChar[b], c.mapped[b] = uint16(b), true
	}
	c.buildEncoder()
	return c
}

func (c *sbcsCoder) buildEncoder() {
	for b := 0; b < 256; b++ {
		if c.mapped[b] {
			c.fromChar.setIfAbsent(c.toChar[b], packBytes(byte(b)))
		}
	}
}

func (c *sbcsCoder) decode(b []byte, report bool) (string, codingError) {
	units := make([]uint16, 0, len(b))
	for _, x := range b {
		if c.mapped[x] {
			units = append(units, c.toChar[x])
			continue
		}
		if report {
			return "", c.badByte
		}
		units = append(units, 0xFFFD)
	}
	return StringFromUTF16(units), codingError{}
}

func (c *sbcsCoder) canEncode(ch uint16) bool { return c.fromChar.get(ch) != 0 }

func (c *sbcsCoder) newEncoder() charsetEncoder { return tableEncoder{&c.fromChar} }

// The single-byte charsets and their overrides (observed black-box: the
// JDK's table differs from x/text's at these bytes).
var (
	// windows-1255: the JDK leaves 0xCA unmapped (x/text: U+05BA).
	sbcsWindows1255 = newSBCS(charmap.Windows1255, false, sbcsOverride{0xCA, -1})
	// KOI8-U: the JDK follows RFC 2319, with box-drawing chars at 0xAE and
	// 0xBE (x/text: U+045E and U+040E).
	sbcsKOI8U = newSBCS(charmap.KOI8U, false, sbcsOverride{0xAE, 0x255D}, sbcsOverride{0xBE, 0x256C})
	// x-MacCyrillic: the JDK has U+00A2 at 0xA2, U+2202 at 0xB6 and U+00A4
	// at 0xFF (x/text: U+0490, U+0491 and U+20AC).
	sbcsMacCyrillic = newSBCS(charmap.MacintoshCyrillic, false,
		sbcsOverride{0xA2, 0x00A2}, sbcsOverride{0xB6, 0x2202}, sbcsOverride{0xFF, 0x00A4})
	// TIS-620 is windows-874 without the nine chars at 0x80 to 0x97.
	sbcsTIS620 = newSBCS(charmap.Windows874, false,
		sbcsOverride{0x80, -1}, sbcsOverride{0x85, -1}, sbcsOverride{0x91, -1}, sbcsOverride{0x92, -1},
		sbcsOverride{0x93, -1}, sbcsOverride{0x94, -1}, sbcsOverride{0x95, -1}, sbcsOverride{0x96, -1},
		sbcsOverride{0x97, -1})
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

// Go-only: language-adapter — ISO-2022-JP (RFC 1468), with the JIS X 0208
// table of EUC-JP (charset_euc_jp.go). Escape sequences switch between
// ASCII (ESC ( B), JIS X 0201 Roman (ESC ( J), JIS X 0201 katakana
// (ESC ( I) and JIS X 0208 (ESC $ @ or ESC $ B).
//
// Observed behaviour of the JDK's decoder (black-box, REPORT mode):
//   - ESC followed by anything but '$' or '(' is malformed, length 2; ESC
//     '$' or ESC '(' followed by an unknown final byte is malformed, length
//     3; an escape sequence cut short by the end of the input is malformed
//     with the length of what is left;
//   - SO (0x0E) switches to katakana and remembers the mode it left; SI
//     (0x0F) returns to the remembered mode (ASCII if SO never came);
//   - in ASCII and Roman mode, bytes from 0x80 are malformed, length 1; in
//     Roman mode 0x5C is the yen sign and 0x7E the overline;
//   - in katakana mode the bytes 0x00 to 0x5F decode to U+FF40 + the byte
//     (0x21 to 0x5F are the half-width katakana), and the others are
//     malformed, length 1;
//   - in JIS X 0208 mode a byte from 0x80, or a byte at the end of the
//     input, is malformed, length 1, and any other two bytes that are not
//     a mapped char are unmappable, length 2.
//
// Observed behaviour of the encoder (black-box, CharsetEncoder and
// OutputStreamWriter in REPLACE mode): each char switches to the first
// mode that can encode it, in the order ASCII (all of 0x00 to 0x7F, even
// ESC, SO and SI), Roman (U+00A5, U+203E), katakana, JIS X 0208; an
// unmappable char (a surrogate pair counts as one) becomes the replacement
// 0x21 0x29 (U+FF1F in JIS X 0208), after switching to JIS X 0208, but an
// unpaired surrogate becomes those two bytes in the current mode; the end
// of the output switches back to ASCII.

const (
	iso2022ASCII = iota
	iso2022Roman
	iso2022Kana
	iso2022JIS0208
)

type iso2022JPCoder struct{ t *jisTables }

func newISO2022JP() charsetCoder { return iso2022JPCoder{loadJISTables()} }

func (c iso2022JPCoder) decode(b []byte, report bool) (string, codingError) {
	return decodeSteps(&iso2022JPDecoder{t: c.t}, b, report)
}

type iso2022JPDecoder struct {
	t          *jisTables
	mode, prev int
}

func (d *iso2022JPDecoder) step(out []uint16, b []byte, i int) ([]uint16, int, codingError) {
	b1 := b[i]
	rest := len(b) - i
	switch b1 {
	case 0x1B:
		if rest == 1 {
			return out, 0, malformed(1)
		}
		if b2 := b[i+1]; b2 != '$' && b2 != '(' {
			return out, 0, malformed(2)
		}
		if rest == 2 {
			return out, 0, malformed(2)
		}
		switch string(b[i+1 : i+3]) {
		case "$@", "$B":
			d.mode = iso2022JIS0208
		case "(B":
			d.mode = iso2022ASCII
		case "(J":
			d.mode = iso2022Roman
		case "(I":
			d.mode = iso2022Kana
		default:
			return out, 0, malformed(3)
		}
		return out, 3, codingError{}
	case 0x0E:
		d.prev, d.mode = d.mode, iso2022Kana
		return out, 1, codingError{}
	case 0x0F:
		d.mode = d.prev
		return out, 1, codingError{}
	}
	switch d.mode {
	case iso2022Kana:
		if b1 > 0x5F {
			return out, 0, malformed(1)
		}
		return append(out, 0xFF40+uint16(b1)), 1, codingError{}
	case iso2022JIS0208:
		if b1 >= 0x80 || rest < 2 {
			return out, 0, malformed(1)
		}
		if b2 := b[i+1]; b1 >= 0x21 && b1 <= 0x7E && b2 >= 0x21 && b2 <= 0x7E {
			if ch := d.t.jis0208[b1-0x21][b2-0x21]; ch != 0 {
				return append(out, ch), 2, codingError{}
			}
		}
		return out, 0, unmappable(2)
	}
	if b1 >= 0x80 {
		return out, 0, malformed(1)
	}
	if d.mode == iso2022Roman {
		switch b1 {
		case 0x5C:
			return append(out, 0x00A5), 1, codingError{}
		case 0x7E:
			return append(out, 0x203E), 1, codingError{}
		}
	}
	return append(out, uint16(b1)), 1, codingError{}
}

// jis0208Bytes returns the GL bytes of a JIS X 0208 char, or ok=false.
func (c iso2022JPCoder) jis0208Bytes(ch uint16) (b1, b2 byte, ok bool) {
	v := c.t.fromChar.get(ch)
	if v>>24 != 2 || v>>8&0xFF < 0xA1 {
		return 0, 0, false
	}
	return byte(v>>8) - 0x80, byte(v) - 0x80, true
}

func (c iso2022JPCoder) canEncode(ch uint16) bool {
	if ch < 0x80 || ch == 0x00A5 || ch == 0x203E || ch >= 0xFF61 && ch <= 0xFF9F {
		return true
	}
	_, _, ok := c.jis0208Bytes(ch)
	return ok
}

func (c iso2022JPCoder) newEncoder() charsetEncoder { return &iso2022JPEncoder{coder: c} }

type iso2022JPEncoder struct {
	coder iso2022JPCoder
	mode  int
}

var iso2022Escapes = [...]string{iso2022ASCII: "\x1b(B", iso2022Roman: "\x1b(J", iso2022Kana: "\x1b(I", iso2022JIS0208: "\x1b$B"}

func (e *iso2022JPEncoder) switchTo(dst []byte, mode int) []byte {
	if e.mode != mode {
		e.mode = mode
		dst = append(dst, iso2022Escapes[mode]...)
	}
	return dst
}

func (e *iso2022JPEncoder) encodeRune(dst []byte, r rune) []byte {
	switch {
	case r < 0x80:
		return append(e.switchTo(dst, iso2022ASCII), byte(r))
	case r == 0x00A5:
		return append(e.switchTo(dst, iso2022Roman), 0x5C)
	case r == 0x203E:
		return append(e.switchTo(dst, iso2022Roman), 0x7E)
	case r >= 0xFF61 && r <= 0xFF9F:
		return append(e.switchTo(dst, iso2022Kana), byte(r-0xFF61+0x21))
	}
	if r < 0x10000 {
		if b1, b2, ok := e.coder.jis0208Bytes(uint16(r)); ok {
			return append(e.switchTo(dst, iso2022JIS0208), b1, b2)
		}
	}
	return append(e.switchTo(dst, iso2022JIS0208), 0x21, 0x29)
}

func (e *iso2022JPEncoder) encodeMalformed(dst []byte) []byte { return append(dst, 0x21, 0x29) }

func (e *iso2022JPEncoder) finish(dst []byte) []byte { return e.switchTo(dst, iso2022ASCII) }
