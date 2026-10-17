package jutil

// Go-only: language-adapter — the UTF-32 charsets (UTF-32, UTF-32BE,
// UTF-32LE, X-UTF-32BE-BOM and X-UTF-32LE-BOM), written from the Unicode
// standard's definition of the encoding form and the Javadoc of
// java.nio.charset.Charset (byte-order marks).
//
// Observed behaviour (black-box, CharsetDecoder in REPORT and REPLACE
// mode, CharsetEncoder and OutputStreamWriter):
//   - a byte-order mark is honoured and dropped only in the first four
//     bytes. UTF-32 takes the order from it (00 00 FE FF big-endian, FF FE
//     00 00 little-endian) and defaults to big-endian; the other four drop
//     a mark in their own order and decode anything else as a char;
//   - a unit above U+10FFFF is malformed with length 4, a trailing partial
//     unit is malformed with its length (1 to 3), and a unit naming a
//     surrogate code point decodes to that lone surrogate;
//   - the encoders replace an unpaired surrogate with U+FFFD; only the two
//     -BOM charsets write a mark, before their first char.

type utf32Coder struct {
	bigEndian bool
	// detect: the order comes from a leading mark (UTF-32).
	detect bool
	// writeMark: the encoder writes a mark first (the -BOM charsets).
	writeMark bool
}

func (c utf32Coder) decode(b []byte, report bool) (string, codingError) {
	return decodeSteps(&utf32Decoder{bigEndian: c.bigEndian, detect: c.detect}, b, report)
}

func (utf32Coder) canEncode(uint16) bool { return true }

func (c utf32Coder) newEncoder() charsetEncoder { return &utf32Encoder{coder: c} }

type utf32Decoder struct {
	bigEndian bool
	detect    bool
}

func (d *utf32Decoder) step(out []uint16, b []byte, i int) ([]uint16, int, codingError) {
	if rest := len(b) - i; rest < 4 {
		return out, 0, malformed(rest)
	}
	be := uint32(b[i])<<24 | uint32(b[i+1])<<16 | uint32(b[i+2])<<8 | uint32(b[i+3])
	le := uint32(b[i+3])<<24 | uint32(b[i+2])<<16 | uint32(b[i+1])<<8 | uint32(b[i])
	if i == 0 {
		switch {
		case d.detect && be == 0xFEFF:
			return out, 4, codingError{}
		case d.detect && le == 0xFEFF:
			d.bigEndian = false
			return out, 4, codingError{}
		case !d.detect && d.bigEndian && be == 0xFEFF, !d.detect && !d.bigEndian && le == 0xFEFF:
			return out, 4, codingError{}
		}
	}
	v := le
	if d.bigEndian {
		v = be
	}
	switch {
	case v < 0x10000:
		return append(out, uint16(v)), 4, codingError{}
	case v <= 0x10FFFF:
		v -= 0x10000
		return append(out, uint16(0xD800+v>>10), uint16(0xDC00+v&0x3FF)), 4, codingError{}
	}
	return out, 0, malformed(4)
}

type utf32Encoder struct {
	coder   utf32Coder
	started bool
}

func (e *utf32Encoder) unit(dst []byte, v rune) []byte {
	if e.coder.writeMark && !e.started {
		dst = e.raw(dst, 0xFEFF)
	}
	e.started = true
	return e.raw(dst, v)
}

func (e *utf32Encoder) raw(dst []byte, v rune) []byte {
	if e.coder.bigEndian {
		return append(dst, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
	return append(dst, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
}

func (e *utf32Encoder) encodeRune(dst []byte, r rune) []byte { return e.unit(dst, r) }
func (e *utf32Encoder) encodeMalformed(dst []byte) []byte    { return e.unit(dst, 0xFFFD) }
func (e *utf32Encoder) finish(dst []byte) []byte             { return dst }
