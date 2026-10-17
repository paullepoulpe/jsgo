package jutil

// Go-only: language-adapter — the conversions behind Charset, in the shape
// the jutil callers need: whole-array decoding in the two modes the Java
// callers use (CodingErrorAction.REPLACE for new String(bytes, cs) and
// InputStreamReader, REPORT for Files.readString and
// Files.newBufferedReader), and a streaming REPLACE-mode encoder for
// OutputStreamWriter (see the Javadoc of CharsetDecoder, CharsetEncoder
// and CodingErrorAction).

// charsetCoder holds the conversions of one charset.
type charsetCoder interface {
	// decode decodes all of b. In report mode it stops at the first
	// malformed or unmappable sequence and returns it; otherwise each such
	// sequence becomes U+FFFD (the decoders' replacement in every charset
	// here) and decoding resumes after it.
	decode(b []byte, report bool) (string, codingError)
	// canEncode is CharsetEncoder.canEncode(char) for a char that is not a
	// surrogate.
	canEncode(c uint16) bool
	// newEncoder returns an encoder for one output stream.
	newEncoder() charsetEncoder
}

// codingError describes a malformed or unmappable input sequence, as a
// CoderResult does (length 0: no error).
type codingError struct {
	length     int
	unmappable bool
}

// charsetStepper decodes one character (or skips a sequence that only
// changes state) at a time.
type charsetStepper interface {
	// step decodes the bytes at b[i:], the rest of the whole input. It
	// appends the decoded chars to out and returns the number of bytes
	// consumed, or returns a codingError for the erroneous sequence at i
	// (and appends nothing).
	step(out []uint16, b []byte, i int) ([]uint16, int, codingError)
}

// charsetEncoder encodes code points in REPLACE mode, keeping any shift
// state between calls.
type charsetEncoder interface {
	// encodeRune appends the encoding of r (any code point but a
	// surrogate), or the charset's replacement when r is unmappable.
	encodeRune(dst []byte, r rune) []byte
	// encodeMalformed appends the replacement for an unpaired surrogate.
	encodeMalformed(dst []byte) []byte
	// finish appends the bytes that end the output (the return to the
	// initial shift state), at the end of the input.
	finish(dst []byte) []byte
}

// decodeSteps is charsetCoder.decode for a decoder that works one
// sequence at a time.
func decodeSteps(d charsetStepper, b []byte, report bool) (string, codingError) {
	out := make([]uint16, 0, len(b))
	for i := 0; i < len(b); {
		var n int
		var e codingError
		out, n, e = d.step(out, b, i)
		if e.length > 0 {
			if report {
				return "", e
			}
			out = append(out, 0xFFFD)
			i += e.length
			continue
		}
		i += n
	}
	return StringFromUTF16(out), codingError{}
}

func malformed(n int) codingError  { return codingError{length: n} }
func unmappable(n int) codingError { return codingError{length: n, unmappable: true} }

// charTable maps BMP chars to uint32 values (0: absent), in pages of 256.
type charTable struct{ pages [256]*[256]uint32 }

func (t *charTable) get(c uint16) uint32 {
	if p := t.pages[c>>8]; p != nil {
		return p[c&0xFF]
	}
	return 0
}

func (t *charTable) set(c uint16, v uint32) {
	p := t.pages[c>>8]
	if p == nil {
		p = new([256]uint32)
		t.pages[c>>8] = p
	}
	p[c&0xFF] = v
}

// setIfAbsent records v for c unless c already has a value.
func (t *charTable) setIfAbsent(c uint16, v uint32) {
	if t.get(c) == 0 {
		t.set(c, v)
	}
}

// Encoded byte sequences of up to three bytes are packed into a uint32 as
// the length (top byte) and the bytes (big-endian, low three bytes).
func packBytes(b ...byte) uint32 {
	var v uint32
	for _, x := range b {
		v = v<<8 | uint32(x)
	}
	return uint32(len(b))<<24 | v
}

func appendPacked(dst []byte, v uint32) []byte {
	switch v >> 24 {
	case 1:
		return append(dst, byte(v))
	case 2:
		return append(dst, byte(v>>8), byte(v))
	default:
		return append(dst, byte(v>>16), byte(v>>8), byte(v))
	}
}

// tableEncoder is the stateless encoder of a charset whose encodable chars
// are all in a charTable of packed byte sequences, and whose replacement
// is '?'.
type tableEncoder struct{ table *charTable }

func (e tableEncoder) encodeRune(dst []byte, r rune) []byte {
	if r < 0x10000 {
		if v := e.table.get(uint16(r)); v != 0 {
			return appendPacked(dst, v)
		}
	}
	return append(dst, '?')
}

func (tableEncoder) encodeMalformed(dst []byte) []byte { return append(dst, '?') }
func (tableEncoder) finish(dst []byte) []byte          { return dst }
