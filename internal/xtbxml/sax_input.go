package xtbxml

// Go-only: language-adapter — the input side of the SAX parser in
// sax_parser.go: encoding detection, decoding, and the character buffer the
// scanner (sax_scanner.go) reads from.
//
// Everything here is written from the XML 1.0 (Fifth Edition) specification
// (§4.3.3 and Appendix F for encoding detection, §2.11 for line ends) and
// checked black-box against the JDK (tools/oracle/java/JutilSaxVectors.java).
// The parts XML 1.0 leaves to the implementation are observable through
// XtbMessageBundle, because it turns every ContentHandler.characters() call
// into a separate message part, and the JDK ends a characters() call where
// its input buffer ends. So the buffer below reproduces how the JDK fills
// its own, as measured from outside:
//
//   - Observed behaviour (black-box): the buffer holds at most 8192 chars.
//     The first read takes the first 32 bytes of the stream (a byte order
//     mark included), in two parts: the four bytes read to detect the
//     encoding, then the rest. Every later read asks for (8192 - kept)
//     chars, where kept is the number of chars carried over from the old
//     buffer, and reads that many bytes (UTF-16: twice as many, at most
//     8192). A read whose last byte starts a multi-byte sequence reads on
//     to complete that character, so a read of n bytes can yield n+1 chars
//     (a supplementary character at the end).
//   - Observed behaviour (black-box): while the XML declaration is scanned,
//     reads past the first 32 bytes take one byte (plus completion) at a
//     time, so the buffer ends exactly where the declaration ends.
//   - Observed behaviour (black-box): a malformed byte sequence ends the
//     read that meets it; the chars decoded before it are delivered, and
//     the error is thrown by the next read (immediately when the read has
//     nothing to deliver).
//
// The scanner decides what to keep on each read; see the load methods.

import (
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

const (
	inputBufferChars = 8192
	firstReadBytes   = 32
)

type xmlEncoding int

const (
	encUTF8 xmlEncoding = iota
	encUTF16BE
	encUTF16LE
	encLatin1
	encASCII
	// encUTF8Lenient is the reader the JDK uses for a UTF-8 declared
	// by a Java-only alias ("UTF8"). Observed behaviour (black-box): its
	// reads count chars, not bytes, and a malformed sequence becomes
	// U+FFFD instead of an error: each maximal prefix of a well-formed
	// sequence, or each other byte, gives one U+FFFD; an encoded surrogate
	// (ED A0..BF xx) gives one for its three bytes.
	encUTF8Lenient
	// encUTF16BELenient and encUTF16LELenient are the readers the
	// JDK uses for a UTF-16 declared by a Java-only alias ("UTF16",
	// "UnicodeBig", "UTF_16LE", ...). Observed behaviour (black-box): their
	// reads count chars, not bytes, and malformed input becomes U+FFFD: a
	// low surrogate on its own gives one, a high surrogate and the unit
	// after it (when that is not a low surrogate) give one, and so does
	// what is left at the end of the input (an odd byte, or a high
	// surrogate without a whole unit after it).
	encUTF16BELenient
	encUTF16LELenient
	// encUTF16Lenient and encUTF16LELenientBOM are those readers for
	// the aliases that look for a byte order mark: when the first two
	// bytes they read are one, it is dropped and sets the byte order;
	// otherwise the order is big-endian or little-endian respectively.
	encUTF16Lenient
	encUTF16LELenientBOM
	// encASCIILenient is the reader for US-ASCII declared by a Java-only
	// alias ("ASCII7"). Observed behaviour (black-box): a byte above 0x7F
	// becomes U+FFFD instead of an error.
	encASCIILenient
)

// charInput is the decoded character buffer over a byte stream.
type charInput struct {
	src    io.Reader
	ahead  []byte // bytes read for encoding detection, not yet decoded
	srcEOF bool
	enc    xmlEncoding
	bomLen int

	// deferred is a decoding error met by the last read after it had
	// decoded something; the next read throws it.
	deferred string
	raw      []byte

	buf   []uint16
	count int // chars in buf
	pos   int // next char to scan
	base  int // offset in the input of buf[0], in chars

	line, column int

	firstRead bool
	splitRead bool
	declMode  bool

	s *scanner
}

func newCharInput(s *scanner, src io.Reader) *charInput {
	return &charInput{
		src:       src,
		buf:       make([]uint16, inputBufferChars+2),
		line:      1,
		column:    1,
		firstRead: true,
		s:         s,
	}
}

func (in *charInput) failed() bool { return in.s.err != nil }

// detectEncoding looks at the first bytes of a byte stream, as XML 1.0
// Appendix F describes: a byte order mark, or the first characters of
// "<?xml" in UTF-16. Everything else is read as UTF-8 until an encoding
// declaration says otherwise.
func (in *charInput) detectEncoding() {
	in.ahead = in.readRaw(4)
	b := in.ahead
	switch {
	case len(b) >= 3 && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF:
		in.enc, in.bomLen = encUTF8, 3
	case len(b) >= 2 && b[0] == 0xFE && b[1] == 0xFF:
		in.enc, in.bomLen = encUTF16BE, 2
	case len(b) >= 2 && b[0] == 0xFF && b[1] == 0xFE:
		in.enc, in.bomLen = encUTF16LE, 2
	case len(b) >= 4 && b[0] == 0x00 && b[1] == 0x3C && b[2] == 0x00 && b[3] == 0x3F:
		in.enc = encUTF16BE
	case len(b) >= 4 && b[0] == 0x3C && b[1] == 0x00 && b[2] == 0x3F && b[3] == 0x00:
		in.enc = encUTF16LE
	case len(b) >= 4 && ((b[0] == 0 && b[1] == 0) || (b[2] == 0 && b[3] == 0) ||
		(b[0] == 0x4C && b[1] == 0x6F && b[2] == 0xA7 && b[3] == 0x94)):
		panic(jutil.NewUnsupportedOperationExceptionMsg("SAX parser: UCS-4 and EBCDIC input are not emulated"))
	default:
		in.enc = encUTF8
	}
	in.ahead = in.ahead[in.bomLen:]
}

// readRaw reads up to n bytes, fewer only at the end of the stream.
func (in *charInput) readRaw(n int) []byte {
	if n <= 0 || in.failed() {
		return nil
	}
	if cap(in.raw) < n {
		in.raw = make([]byte, n)
	}
	buf := in.raw[:n]
	copied := copy(buf, in.ahead)
	in.ahead = in.ahead[copied:]
	out := buf[:copied]
	for len(out) < n && !in.srcEOF {
		m, err := in.src.Read(buf[len(out):n])
		out = buf[:len(out)+m]
		if err == io.EOF {
			in.srcEOF = true
		} else if err != nil {
			in.s.setErr(jutil.NewIOException(err.Error()))
			return nil
		}
	}
	if len(out) == 0 {
		return nil
	}
	// Return a copy the caller may keep while readRaw is called again for
	// continuation bytes.
	return slices.Clone(out)
}

func (in *charInput) readByte() (byte, bool) {
	b := in.readRaw(1)
	if len(b) == 0 {
		return 0, false
	}
	return b[0], true
}

// setEncoding switches the decoder for the bytes after the encoding
// declaration named name. Observed behaviour (black-box): "UTF-16" keeps
// the byte order of a stream detected as UTF-16; any other name switches
// to its own decoder, whatever was detected.
func (in *charInput) setEncoding(enc xmlEncoding, name string) {
	if strings.EqualFold(name, "UTF-16") && (in.enc == encUTF16BE || in.enc == encUTF16LE) {
		return
	}
	in.enc = enc
}

// readChars appends decoded chars to buf[count:], reading about n bytes; it
// returns the number of chars appended (0 at the end of the input).
func (in *charInput) readChars(n int) int {
	if in.failed() {
		return 0
	}
	if in.deferred != "" {
		msg := in.deferred
		in.deferred = ""
		in.s.fatal(msg)
		return 0
	}
	switch in.enc {
	case encUTF8Lenient:
		c := in.readLenientUTF8(n)
		if in.failed() {
			return 0
		}
		return c
	case encUTF16BELenient, encUTF16LELenient, encUTF16Lenient, encUTF16LELenientBOM:
		c := in.readLenientUTF16(n)
		if in.failed() {
			return 0
		}
		return c
	}
	start := in.count
	if (in.enc == encUTF16BE || in.enc == encUTF16LE) && !in.splitRead {
		// Observed behaviour (black-box): the UTF-16 reader reads two bytes
		// per char asked for, at most 8192.
		n = min(2*n, inputBufferChars)
	}
	if in.splitRead {
		// Observed behaviour (black-box): the first read decodes the bytes
		// taken for encoding detection on their own (completing a
		// character that straddles their end), then reads the rest.
		in.splitRead = false
		head := in.ahead
		in.ahead = nil
		if !in.decodeBlock(head, start) {
			if in.failed() {
				return 0
			}
			return in.count - start
		}
		n -= len(head)
	}
	raw := in.readRaw(n)
	if in.failed() {
		return 0
	}
	in.decodeBlock(raw, start)
	if in.failed() {
		return 0
	}
	return in.count - start
}

// decodeBlock decodes bytes into buf; it returns false when it stopped at
// a malformed sequence (the error then waits for the next read).
func (in *charInput) decodeBlock(bytes []byte, start int) bool {
	for i := 0; i < len(bytes); {
		var msg string
		switch in.enc {
		case encUTF16BE, encUTF16LE:
			i, msg = in.decodeUTF16(bytes, i)
		case encLatin1:
			in.put(uint16(bytes[i]))
			i++
		case encASCIILenient:
			if bytes[i] >= 0x80 {
				in.put(0xFFFD)
			} else {
				in.put(uint16(bytes[i]))
			}
			i++
		case encASCII:
			if bytes[i] >= 0x80 {
				msg = "Byte \"" + strconv.Itoa(int(bytes[i])) + "\" is not a member of the (7-bit) ASCII character set."
			} else {
				in.put(uint16(bytes[i]))
				i++
			}
		default:
			i, msg = in.decodeUTF8(bytes, i)
		}
		if msg != "" {
			if in.failed() {
				return false
			}
			// Observed behaviour (black-box): a UTF-8 error waits for the
			// next read when this one has chars to deliver, except the one
			// about the planes of a 4-byte sequence; an ASCII one does not.
			if in.count > start && in.enc == encUTF8 && !strings.HasPrefix(msg, "High surrogate") {
				in.deferred = msg
				return false
			}
			in.s.fatal(msg)
			return false
		}
	}
	return true
}

func (in *charInput) put(c uint16) {
	in.buf[in.count] = c
	in.count++
}

// moreByte returns the byte at bytes[i], reading it from the stream when
// the read block has ended (completion of a character that straddles the
// end of the block).
func (in *charInput) moreByte(bytes *[]byte, i int) (byte, bool) {
	if i < len(*bytes) {
		return (*bytes)[i], true
	}
	b, ok := in.readByte()
	if !ok {
		return 0, false
	}
	*bytes = append(*bytes, b)
	return b, true
}

// decodeUTF8 decodes one character at bytes[i] per RFC 3629. The messages
// for malformed input were determined black-box (they name the position of
// the offending byte in its sequence).
func (in *charInput) decodeUTF8(bytes []byte, i int) (int, string) {
	b0 := bytes[i]
	if b0 < 0x80 {
		in.put(uint16(b0))
		return i + 1, ""
	}
	var size int
	switch {
	case b0 >= 0xC2 && b0 <= 0xDF:
		size = 2
	case b0 >= 0xE0 && b0 <= 0xEF:
		size = 3
	case b0 >= 0xF0 && b0 <= 0xF7:
		size = 4
	default:
		return i, "Invalid byte 1 of 1-byte UTF-8 sequence."
	}
	sizeText := strconv.Itoa(size)
	var bs [4]byte
	bs[0] = b0
	for k := 1; k < size; k++ {
		b, ok := in.moreByte(&bytes, i+k)
		if !ok {
			return i, "Expected byte " + strconv.Itoa(k+1) + " of " + sizeText + "-byte UTF-8 sequence."
		}
		bad := !isUTF8Continuation(b)
		if k == 1 {
			switch {
			case b0 == 0xE0 && b < 0xA0, b0 == 0xED && b >= 0xA0, b0 == 0xF0 && b < 0x90:
				bad = true
			}
		}
		if bad {
			return i, "Invalid byte " + strconv.Itoa(k+1) + " of " + sizeText + "-byte UTF-8 sequence."
		}
		bs[k] = b
	}
	switch size {
	case 2:
		in.put(uint16(b0&0x1F)<<6 | uint16(bs[1]&0x3F))
	case 3:
		in.put(uint16(b0&0x0F)<<12 | uint16(bs[1]&0x3F)<<6 | uint16(bs[2]&0x3F))
	default:
		planes := int(b0&0x07)<<2 | int(bs[1]&0x30)>>4
		if planes > 0x10 {
			return i, "High surrogate bits in UTF-8 sequence must not exceed 0x10 but found 0x" +
				strconv.FormatInt(int64(planes), 16) + "."
		}
		cp := rune(b0&0x07)<<18 | rune(bs[1]&0x3F)<<12 | rune(bs[2]&0x3F)<<6 | rune(bs[3]&0x3F)
		cp -= 0x10000
		in.put(uint16(0xD800 + cp>>10))
		in.put(uint16(0xDC00 + cp&0x3FF))
	}
	return i + size, ""
}

func (in *charInput) decodeUTF16(bytes []byte, i int) (int, string) {
	b1, ok := in.moreByte(&bytes, i+1)
	if !ok {
		// Observed behaviour (black-box): the JDK reports an odd byte at
		// the end with this message (naming UTF-8).
		return i, "Expected byte 2 of 2-byte UTF-8 sequence."
	}
	if in.enc == encUTF16BE {
		in.put(uint16(bytes[i])<<8 | uint16(b1))
	} else {
		in.put(uint16(b1)<<8 | uint16(bytes[i]))
	}
	return i + 2, ""
}

// readLenientUTF8 reads up to n chars with the lenient UTF-8 reader (see
// encUTF8Lenient). A supplementary character that does not fit is left
// for the next read.
func (in *charInput) readLenientUTF8(n int) int {
	start := in.count
	for in.count-start < n {
		var seq [4]byte
		k, ok, eof := 0, true, false
		// more reads one more byte of the sequence if valid accepts it;
		// otherwise the byte is given back.
		more := func(valid func(byte) bool) bool {
			b, got := in.readByte()
			if !got {
				eof = true
				return false
			}
			if !valid(b) {
				in.ahead = append([]byte{b}, in.ahead...)
				return false
			}
			seq[k] = b
			k++
			return true
		}
		if seq[0], ok = in.readByte(); !ok {
			break
		}
		k = 1
		b0 := seq[0]
		var units [2]uint16
		size := 1
		units[0] = 0xFFFD
		switch {
		case b0 < 0x80:
			units[0] = uint16(b0)
		case b0 >= 0xC2 && b0 <= 0xDF:
			if more(isUTF8Continuation) {
				units[0] = uint16(b0&0x1F)<<6 | uint16(seq[1]&0x3F)
			}
		case b0 >= 0xE0 && b0 <= 0xEF:
			if more(func(b byte) bool { return isUTF8Continuation(b) && !(b0 == 0xE0 && b < 0xA0) }) &&
				more(isUTF8Continuation) {
				c := uint16(b0&0x0F)<<12 | uint16(seq[1]&0x3F)<<6 | uint16(seq[2]&0x3F)
				if c < 0xD800 || c > 0xDFFF {
					units[0] = c
				}
			}
		case b0 >= 0xF0 && b0 <= 0xF4:
			if more(func(b byte) bool {
				return isUTF8Continuation(b) && !(b0 == 0xF0 && b < 0x90) && !(b0 == 0xF4 && b >= 0x90)
			}) && more(isUTF8Continuation) && more(isUTF8Continuation) {
				cp := rune(b0&0x07)<<18 | rune(seq[1]&0x3F)<<12 | rune(seq[2]&0x3F)<<6 | rune(seq[3]&0x3F)
				cp -= 0x10000
				units[0], units[1] = uint16(0xD800+cp>>10), uint16(0xDC00+cp&0x3FF)
				size = 2
			}
		}
		if in.count-start+size > n || eof && in.count > start {
			// A sequence cut by the end of the input becomes U+FFFD in the
			// next read.
			in.ahead = slices.Concat(seq[:k], in.ahead)
			break
		}
		for _, u := range units[:size] {
			in.put(u)
		}
	}
	return in.count - start
}

func isUTF8Continuation(b byte) bool { return b&0xC0 == 0x80 }

// readLenientUTF16 reads up to n chars with a lenient UTF-16 reader (see
// encUTF16BELenient). A supplementary character that does not fit is
// left for the next read.
func (in *charInput) readLenientUTF16(n int) int {
	if in.enc == encUTF16Lenient || in.enc == encUTF16LELenientBOM {
		order := encUTF16BELenient
		if in.enc == encUTF16LELenientBOM {
			order = encUTF16LELenient
		}
		b0, ok0 := in.readByte()
		b1, ok1 := in.readByte()
		switch {
		case ok1 && b0 == 0xFE && b1 == 0xFF:
			order = encUTF16BELenient
		case ok1 && b0 == 0xFF && b1 == 0xFE:
			order = encUTF16LELenient
		case ok1:
			in.ahead = append([]byte{b0, b1}, in.ahead...)
		case ok0:
			in.ahead = append([]byte{b0}, in.ahead...)
		}
		in.enc = order
	}
	// unit reads one code unit; whole is false when fewer than two bytes
	// were left (got says how many).
	var got []byte
	unit := func() (u uint16, whole bool) {
		b0, ok := in.readByte()
		if !ok {
			return 0, false
		}
		got = append(got, b0)
		b1, ok := in.readByte()
		if !ok {
			return 0, false
		}
		got = append(got, b1)
		if in.enc == encUTF16LELenient {
			b0, b1 = b1, b0
		}
		return uint16(b0)<<8 | uint16(b1), true
	}
	start := in.count
	for in.count-start < n {
		got = got[:0]
		u, whole := unit()
		if len(got) == 0 {
			break
		}
		if !whole && in.count > start {
			// What is left at the end of the input becomes U+FFFD in the
			// next read.
			in.ahead = slices.Concat(got, in.ahead)
			break
		}
		if !whole || isLowSurrogate(int(u)) {
			in.put(0xFFFD)
			continue
		}
		if !isHighSurrogate(int(u)) {
			in.put(u)
			continue
		}
		low, whole := unit()
		if !whole && in.count > start {
			in.ahead = slices.Concat(got, in.ahead)
			break
		}
		if !whole || !isLowSurrogate(int(low)) {
			in.put(0xFFFD)
			continue
		}
		if in.count-start+2 > n {
			in.ahead = slices.Concat(got, in.ahead)
			break
		}
		in.put(u)
		in.put(low)
	}
	return in.count - start
}

// readSize is the number of bytes the next read asks for, given the number
// of chars kept in the buffer.
func (in *charInput) readSize(kept int) int {
	switch {
	case in.firstRead:
		in.firstRead = false
		in.splitRead = true
		return firstReadBytes - in.bomLen
	case in.declMode:
		return 1
	}
	return inputBufferChars - kept
}

// load moves buf[keepFrom:count] to the front of the buffer and reads more
// after it. It returns the number of chars read (0 at the end of input).
func (in *charInput) load(keepFrom int) int {
	kept := in.count - keepFrom
	copy(in.buf, in.buf[keepFrom:in.count])
	in.pos -= keepFrom
	in.base += keepFrom
	in.count = kept
	return in.readChars(in.readSize(kept))
}

// loadMore makes room for more chars without dropping any: it reads into
// the free end of the buffer, or, when the buffer is full, keeps the chars
// from pos on (like load(pos)).
func (in *charInput) loadMore() int {
	if in.count < inputBufferChars && !in.firstRead {
		if in.declMode {
			return in.readChars(1)
		}
		return in.readChars(inputBufferChars - in.count)
	}
	return in.load(in.pos)
}

// The position helpers below keep line and column in step with pos. A
// newline is counted when the scanner consumes it.

func (in *charInput) advance(n int) {
	in.pos += n
	in.column += n
}

func (in *charInput) newline(n int) {
	in.pos += n
	in.line++
	in.column = 1
}

// encodingByName maps an encoding name from an XML declaration to a
// decoder. Names are compared case-insensitively, except that only
// "UTF-16BE" and "UTF-16LE" in upper case get the JDK's own UTF-16 readers
// (other spellings get the lenient ones). Only the names below are
// emulated, each checked black-box for the reader the JDK picks for it; the
// JDK knows many more.
func encodingByName(name string) (xmlEncoding, bool) {
	switch name {
	case "UTF-16BE":
		return encUTF16BE, true
	case "UTF-16LE":
		return encUTF16LE, true
	}
	switch strings.ToUpper(name) {
	case "UTF-8":
		return encUTF8, true
	case "UTF8":
		return encUTF8Lenient, true
	case "UTF-16":
		return encUTF16BE, true
	case "UTF16", "UTF_16", "UNICODE", "UNICODEBIG":
		return encUTF16Lenient, true
	case "UNICODELITTLE", "X-UTF-16LE-BOM":
		return encUTF16LELenientBOM, true
	case "UTF-16BE", "UTF_16BE", "X-UTF-16BE", "UNICODEBIGUNMARKED":
		return encUTF16BELenient, true
	case "UTF-16LE", "UTF_16LE", "X-UTF-16LE", "UNICODELITTLEUNMARKED":
		return encUTF16LELenient, true
	case "ISO-8859-1", "ISO8859-1", "ISO_8859-1", "ISO8859_1", "ISO_8859_1", "ISO-IR-100", "LATIN1", "L1",
		"CP819", "IBM819", "IBM-819", "CSISOLATIN1":
		return encLatin1, true
	case "US-ASCII", "ASCII", "US", "ISO646-US", "IBM367", "IBM-367", "CP367", "CSASCII", "ANSI_X3.4-1968",
		"ANSI_X3.4-1986", "ISO-IR-6":
		return encASCII, true
	case "ASCII7":
		return encASCIILenient, true
	}
	return 0, false
}

// unsupportedEncodingException is java.io.UnsupportedEncodingException.
type unsupportedEncodingException struct{ jutil.IOException }

// isKnownUnsupportedEncoding reports whether name is a common encoding (IANA name)
// that the JDK supports but this parser does not emulate; any other name
// is treated as unknown to the JDK.
func isKnownUnsupportedEncoding(name string) bool {
	// numbered reports whether n is prefix followed by digits only.
	numbered := func(n, prefix string) bool {
		rest, ok := strings.CutPrefix(n, prefix)
		if !ok || rest == "" {
			return false
		}
		for _, c := range rest {
			if c < '0' || c > '9' {
				return false
			}
		}
		return true
	}
	switch n := strings.ToUpper(name); {
	case strings.HasPrefix(n, "WINDOWS-125"), numbered(n, "ISO-8859-"), strings.HasPrefix(n, "UTF-32"),
		numbered(n, "IBM"), numbered(n, "IBM-"), numbered(n, "CP"), strings.HasPrefix(n, "X-"):
		return true
	default:
		switch n {
		case "KOI8-R", "KOI8-U", "SHIFT_JIS", "EUC-JP", "EUC-KR", "ISO-2022-JP", "ISO-2022-KR", "GBK", "GB2312",
			"GB18030", "BIG5", "BIG5-HKSCS", "TIS-620", "UTF-16-LE", "MACROMAN", "ISO-10646-UCS-2",
			"ISO-10646-UCS-4", "CESU-8":
			return true
		}
	}
	return false
}

// preScanDeclaration looks for an XML declaration and its version before
// the document is reported (the JDK reads the version first, to choose
// between XML 1.0 and 1.1). It consumes "<?xml", white space, "version",
// white space, '=' and white space, as far as they are there, then looks
// at the next five chars (the version literal and one more) without
// consuming them; the declaration is then scanned from where it stopped.
// Observed behaviour (black-box):
//   - the end of input during the pre-scan is a "Premature end of file."
//     error at line and column -1, before setDocumentLocator; a decoding
//     error is reported where the pre-scan met it, also before;
//   - after the pre-scan the position is line 1, column one more than
//     the number of chars consumed: line ends in the pre-scanned part
//     count as columns. But when the version literal is not "1.0" in a
//     matching pair of quotes, the column at the literal's quote is 15
//     (that of "<?xml version='"), and a literal of four chars or more
//     gets as many spaces after its fourth char as the pre-scanned part is
//     longer than "<?xml version=" (the spaces count as columns too);
//   - when the pre-scan (look-ahead included) had to read past the first
//     read of the input (32 bytes, so always in UTF-16), the column is
//     that of the shortest declaration in the same state ("<?xml ",
//     "<?xml version", "<?xml version="), and there are no spaces;
//   - when the pre-scan stayed in the first read and consumed more than
//     "<?xml version=", the pseudo attribute after the version needs no
//     white space before it;
//   - when "<?xml" is followed by a name char instead (a processing
//     instruction such as "<?xml-stylesheet"), the scan restarts at the
//     beginning but keeps the five columns it skipped, so positions on
//     line 1 are five columns too far.
func (s *scanner) preScanDeclaration() {
	in := s.in
	if !s.skipASCII("<?xml") || s.err != nil {
		return
	}
	eof := func() { s.fatalAt("Premature end of file.", -1, -1) }
	c := s.peek()
	if c == stopEOF {
		eof()
		return
	}
	if !isXMLSpace(c) {
		in.pos = 0
		return
	}
	s.hasDecl = true
	in.declMode = true
	firstEnd := in.base + in.count
	beyond := func() bool { return in.base+in.count > firstEnd }
	spaces := func() bool {
		s.skipSpaces()
		if s.peek() == stopEOF {
			eof()
			return false
		}
		return true
	}
	// finish sets the position after the pre-scan; shortest is the length
	// of the shortest declaration prefix in the same state.
	finish := func(shortest int) {
		in.line, in.column = 1, in.base+in.pos+1
		if beyond() {
			in.column = shortest + 1
		}
	}
	if !spaces() {
		return
	}
	if !s.skipASCII("version") {
		if s.err == nil {
			finish(len("<?xml "))
		}
		return
	}
	s.declStage = 1
	if !spaces() {
		return
	}
	if s.peek() != '=' {
		finish(len("<?xml version"))
		return
	}
	in.advance(1)
	s.declStage = 2
	if !spaces() {
		return
	}
	// The literal is read one char at a time: a decoding error is
	// reported after the chars read so far.
	for k := 1; k <= 5; k++ {
		for in.count-in.pos < k {
			line, column := in.line, in.column
			for _, u := range in.buf[in.pos : in.pos+k-1] {
				if u == '\n' || u == '\r' {
					in.line, in.column = in.line+1, 1
				} else {
					in.column++
				}
			}
			n := in.loadMore()
			in.line, in.column = line, column
			if n == 0 || s.err != nil {
				eof()
				return
			}
		}
	}
	b := in.buf[in.pos : in.pos+5]
	quoted := (b[0] == '"' || b[0] == '\'') && b[4] == b[0]
	if quoted && b[1] == '1' && b[2] == '.' && b[3] == '1' {
		panicUnsupported("XML 1.1 documents are")
	}
	s.versionLoose = !beyond() && in.base+in.pos > len("<?xml version=")
	if quoted && b[1] == '1' && b[2] == '.' && b[3] == '0' {
		finish(len("<?xml version="))
		return
	}
	if !beyond() {
		s.versionPad = in.base + in.pos - len("<?xml version=")
	}
	in.line, in.column = 1, len("<?xml version=")+1
}

// xmlDecl scans the rest of the XML declaration (production [23]) from
// where the pre-scan stopped (s.declStage: 0 before the first pseudo
// attribute, 1 before the '=' of version, 2 before its value).
func (s *scanner) xmlDecl() {
	seen := 0
	var encoding string
	haveEncoding, haveStandalone := false, false
	for {
		var which string
		sawSpace := true
		stage := 0
		if seen == 0 {
			stage = s.declStage
		}
		if stage == 0 {
			if seen > 0 {
				sawSpace = s.skipSpaces() || seen == 1 && s.versionLoose
			}
			c := s.peek()
			if c == '?' || c == stopEOF {
				if seen == 0 {
					s.fatal("The version is required in the XML declaration.")
					return
				}
				break
			}
			switch {
			case s.skipASCII("version"):
				which = "version"
			case s.skipASCII("encoding"):
				which = "encoding"
			case s.skipASCII("standalone"):
				which = "standalone"
			default:
				s.fatal("A pseudo attribute name is expected.")
				return
			}
			s.skipSpaces()
		} else {
			which = "version"
		}
		if stage < 2 {
			if !s.skipChar('=') {
				s.fatalOrEOF("The ' = ' character must follow \"" + which + "\" in the XML declaration.")
				return
			}
			s.skipSpaces()
		}
		quote := s.peek()
		if quote != '"' && quote != '\'' {
			s.fatalOrEOF("The value following \"" + which + "\" in the XML declaration must be a quoted string.")
			return
		}
		s.in.advance(1)
		if s.in.pos == s.in.count-1 {
			// Observed behaviour (black-box): when the value starts with
			// the last char of the buffer, more is read before it.
			s.in.load(s.in.pos)
			if s.err != nil {
				return
			}
		}
		var value []uint16
		padded := seen != 0
		for {
			c := s.peek()
			if c == stopEOF {
				s.fatal(s.eofMessage())
				return
			}
			if isHighSurrogate(c) && isLowSurrogate(s.peekSecond()) {
				value = append(value, uint16(c), s.in.buf[s.in.pos+1])
				s.in.advance(2)
			} else if isHighSurrogate(c) {
				s.loneHighSurrogate(c)
				return
			} else {
				if !isXMLChar(c) {
					s.fatal("An invalid XML character (Unicode: " + formatHex(c) + ") was found in the XML declaration.")
					return
				}
				if c = s.next(); c == quote {
					break
				}
				value = append(value, uint16(c))
			}
			if !padded && len(value) >= 4 {
				for range s.versionPad {
					value = append(value, ' ')
				}
				s.in.column += s.versionPad
				padded = true
			}
		}
		v := jutil.StringFromUTF16(value)
		switch {
		case seen == 0 && which != "version":
			s.fatal("The version is required in the XML declaration.")
		case seen == 0:
			if v != "1.0" {
				s.fatal("XML version \"" + v + "\" is not supported, only XML 1.0 is supported.")
			}
		// Observed behaviour (black-box): the messages for a pseudo
		// attribute out of order.
		case haveStandalone:
			s.fatal("No more pseudo attributes are allowed.")
		case haveEncoding && which != "standalone":
			s.fatal("The standalone name in XML declaration may be misspelled.")
		case which == "version":
			s.fatal("The encoding declaration is required in the text declaration.")
		case !sawSpace:
			// Observed behaviour (black-box): the message names encoding
			// for standalone too.
			s.fatal("White space is required before the encoding pseudo attribute in the XML declaration.")
		case which == "encoding":
			haveEncoding = true
			encoding = v
		default:
			if v != "yes" && v != "no" {
				s.fatal("The standalone document declaration value must be \"yes\" or \"no\", not \"" + v + "\".")
			}
			haveStandalone = true
			s.standalone = v == "yes"
		}
		if s.err != nil {
			return
		}
		seen++
	}
	if !s.skipChar('?') || !s.skipChar('>') {
		s.fatalOrEOF("The XML declaration must end with \"?>\".")
		return
	}
	s.in.declMode = false
	if haveEncoding {
		if !isEncName(encoding) {
			s.fatal("Invalid encoding name \"" + encoding + "\".")
			return
		}
		enc, ok := encodingByName(encoding)
		if !ok {
			if isKnownUnsupportedEncoding(encoding) {
				panicUnsupported("encoding \"" + encoding + "\" is")
			}
			// Observed behaviour (black-box): the JDK throws
			// java.io.UnsupportedEncodingException for an encoding it does
			// not know.
			e := &unsupportedEncodingException{}
			jutil.InitThrowableMsg(e, "java.io.UnsupportedEncodingException", encoding, nil)
			s.setErr(e)
			return
		}
		s.in.setEncoding(enc, encoding)
	}
}

// isEncName is production [81] EncName.
func isEncName(name string) bool {
	if name == "" {
		return false
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		letter := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
		if i == 0 {
			if !letter {
				return false
			}
		} else if !letter && !(c >= '0' && c <= '9') && c != '.' && c != '_' && c != '-' {
			return false
		}
	}
	return true
}
