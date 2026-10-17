package jutil

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
