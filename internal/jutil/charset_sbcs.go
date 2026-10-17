package jutil

import (
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
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
