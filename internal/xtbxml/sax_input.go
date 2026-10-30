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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

const (
	saxBufferChars    = 8192
	saxFirstReadBytes = 32
)

type saxEncoding int

const (
	saxEncUTF8 saxEncoding = iota
	saxEncUTF16BE
	saxEncUTF16LE
	saxEncLatin1
	saxEncASCII
	// saxEncUTF8Lenient is the reader the JDK uses for a UTF-8 declared
	// by a Java-only alias ("UTF8"). Observed behaviour (black-box): its
	// reads count chars, not bytes, and a malformed sequence becomes
	// U+FFFD instead of an error: each maximal prefix of a well-formed
	// sequence, or each other byte, gives one U+FFFD; an encoded surrogate
	// (ED A0..BF xx) gives one for its three bytes.
	saxEncUTF8Lenient
	// saxEncUTF16BELenient and saxEncUTF16LELenient are the readers the
	// JDK uses for a UTF-16 declared by a Java-only alias ("UTF16",
	// "UnicodeBig", "UTF_16LE", ...). Observed behaviour (black-box): their
	// reads count chars, not bytes, and malformed input becomes U+FFFD: a
	// low surrogate on its own gives one, a high surrogate and the unit
	// after it (when that is not a low surrogate) give one, and so does
	// what is left at the end of the input (an odd byte, or a high
	// surrogate without a whole unit after it).
	saxEncUTF16BELenient
	saxEncUTF16LELenient
	// saxEncUTF16Lenient and saxEncUTF16LELenientBOM are those readers for
	// the aliases that look for a byte order mark: when the first two
	// bytes they read are one, it is dropped and sets the byte order;
	// otherwise the order is big-endian or little-endian respectively.
	saxEncUTF16Lenient
	saxEncUTF16LELenientBOM
	// saxEncASCIILenient is the reader for US-ASCII declared by a Java-only
	// alias ("ASCII7"). Observed behaviour (black-box): a byte above 0x7F
	// becomes U+FFFD instead of an error.
	saxEncASCIILenient
)

// saxInput is the decoded character buffer over a byte stream.
type saxInput struct {
	src    io.Reader
	ahead  []byte // bytes read for encoding detection, not yet decoded
	srcEOF bool
	enc    saxEncoding
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

	// fail reports a fatal error at the current position (set by the
	// scanner; it does not return).
	fail func(msg string)
}

func newSaxInput(src io.Reader) *saxInput {
	return &saxInput{
		src:       src,
		buf:       make([]uint16, saxBufferChars+2),
		line:      1,
		column:    1,
		firstRead: true,
	}
}

// detectEncoding looks at the first bytes of a byte stream, as XML 1.0
// Appendix F describes: a byte order mark, or the first characters of
// "<?xml" in UTF-16. Everything else is read as UTF-8 until an encoding
// declaration says otherwise.
func (in *saxInput) detectEncoding() {
	in.ahead = in.readRaw(4)
	b := in.ahead
	switch {
	case len(b) >= 3 && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF:
		in.enc, in.bomLen = saxEncUTF8, 3
	case len(b) >= 2 && b[0] == 0xFE && b[1] == 0xFF:
		in.enc, in.bomLen = saxEncUTF16BE, 2
	case len(b) >= 2 && b[0] == 0xFF && b[1] == 0xFE:
		in.enc, in.bomLen = saxEncUTF16LE, 2
	case len(b) >= 4 && b[0] == 0x00 && b[1] == 0x3C && b[2] == 0x00 && b[3] == 0x3F:
		in.enc = saxEncUTF16BE
	case len(b) >= 4 && b[0] == 0x3C && b[1] == 0x00 && b[2] == 0x3F && b[3] == 0x00:
		in.enc = saxEncUTF16LE
	case len(b) >= 4 && ((b[0] == 0 && b[1] == 0) || (b[2] == 0 && b[3] == 0) ||
		(b[0] == 0x4C && b[1] == 0x6F && b[2] == 0xA7 && b[3] == 0x94)):
		panic(jutil.NewUnsupportedOperationExceptionMsg("SAX parser: UCS-4 and EBCDIC input are not emulated"))
	default:
		in.enc = saxEncUTF8
	}
	in.ahead = in.ahead[in.bomLen:]
}

// readRaw reads up to n bytes, fewer only at the end of the stream.
func (in *saxInput) readRaw(n int) []byte {
	if n <= 0 {
		return nil
	}
	if cap(in.raw) < n {
		in.raw = make([]byte, n)
	}
	out := in.raw[:0]
	for len(in.ahead) > 0 && len(out) < n {
		out = append(out, in.ahead[0])
		in.ahead = in.ahead[1:]
	}
	for len(out) < n && !in.srcEOF {
		m, err := in.src.Read(in.raw[len(out):n])
		out = in.raw[:len(out)+m]
		if err == io.EOF {
			in.srcEOF = true
		} else if err != nil {
			panic(saxBailout{err: jutil.NewIOException(err.Error())})
		}
	}
	// Return a copy the caller may keep while readRaw is called again for
	// continuation bytes.
	return append([]byte(nil), out...)
}

func (in *saxInput) readByte() (byte, bool) {
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
func (in *saxInput) setEncoding(enc saxEncoding, name string) {
	if strings.EqualFold(name, "UTF-16") && (in.enc == saxEncUTF16BE || in.enc == saxEncUTF16LE) {
		return
	}
	in.enc = enc
}

// readChars appends decoded chars to buf[count:], reading about n bytes; it
// returns the number of chars appended (0 at the end of the input).
func (in *saxInput) readChars(n int) int {
	if in.deferred != "" {
		in.fail(in.deferred)
	}
	switch in.enc {
	case saxEncUTF8Lenient:
		return in.readLenientUTF8(n)
	case saxEncUTF16BELenient, saxEncUTF16LELenient, saxEncUTF16Lenient, saxEncUTF16LELenientBOM:
		return in.readLenientUTF16(n)
	}
	start := in.count
	if (in.enc == saxEncUTF16BE || in.enc == saxEncUTF16LE) && !in.splitRead {
		// Observed behaviour (black-box): the UTF-16 reader reads two bytes
		// per char asked for, at most 8192.
		n = min(2*n, saxBufferChars)
	}
	if in.splitRead {
		// Observed behaviour (black-box): the first read decodes the bytes
		// taken for encoding detection on their own (completing a
		// character that straddles their end), then reads the rest.
		in.splitRead = false
		head := in.ahead
		in.ahead = nil
		if !in.decodeBlock(head, start) {
			return in.count - start
		}
		n -= len(head)
	}
	in.decodeBlock(in.readRaw(n), start)
	return in.count - start
}

// decodeBlock decodes bytes into buf; it returns false when it stopped at
// a malformed sequence (the error then waits for the next read).
func (in *saxInput) decodeBlock(bytes []byte, start int) bool {
	for i := 0; i < len(bytes); {
		var msg string
		switch in.enc {
		case saxEncUTF16BE, saxEncUTF16LE:
			i, msg = in.decodeUTF16(bytes, i)
		case saxEncLatin1:
			in.put(uint16(bytes[i]))
			i++
		case saxEncASCIILenient:
			if bytes[i] >= 0x80 {
				in.put(0xFFFD)
			} else {
				in.put(uint16(bytes[i]))
			}
			i++
		case saxEncASCII:
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
			// Observed behaviour (black-box): a UTF-8 error waits for the
			// next read when this one has chars to deliver, except the one
			// about the planes of a 4-byte sequence; an ASCII one does not.
			if in.count > start && in.enc == saxEncUTF8 && !strings.HasPrefix(msg, "High surrogate") {
				in.deferred = msg
				return false
			}
			in.fail(msg)
		}
	}
	return true
}

func (in *saxInput) put(c uint16) {
	in.buf[in.count] = c
	in.count++
}

// moreByte returns the byte at bytes[i], reading it from the stream when
// the read block has ended (completion of a character that straddles the
// end of the block).
func (in *saxInput) moreByte(bytes *[]byte, i int) (byte, bool) {
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
func (in *saxInput) decodeUTF8(bytes []byte, i int) (int, string) {
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
		bad := b&0xC0 != 0x80
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

func (in *saxInput) decodeUTF16(bytes []byte, i int) (int, string) {
	b1, ok := in.moreByte(&bytes, i+1)
	if !ok {
		// Observed behaviour (black-box): the JDK reports an odd byte at
		// the end with this message (naming UTF-8).
		return i, "Expected byte 2 of 2-byte UTF-8 sequence."
	}
	if in.enc == saxEncUTF16BE {
		in.put(uint16(bytes[i])<<8 | uint16(b1))
	} else {
		in.put(uint16(b1)<<8 | uint16(bytes[i]))
	}
	return i + 2, ""
}

// readLenientUTF8 reads up to n chars with the lenient UTF-8 reader (see
// saxEncUTF8Lenient). A supplementary character that does not fit is left
// for the next read.
func (in *saxInput) readLenientUTF8(n int) int {
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
			if more(saxIsContinuation) {
				units[0] = uint16(b0&0x1F)<<6 | uint16(seq[1]&0x3F)
			}
		case b0 >= 0xE0 && b0 <= 0xEF:
			if more(func(b byte) bool { return saxIsContinuation(b) && !(b0 == 0xE0 && b < 0xA0) }) &&
				more(saxIsContinuation) {
				c := uint16(b0&0x0F)<<12 | uint16(seq[1]&0x3F)<<6 | uint16(seq[2]&0x3F)
				if c < 0xD800 || c > 0xDFFF {
					units[0] = c
				}
			}
		case b0 >= 0xF0 && b0 <= 0xF4:
			if more(func(b byte) bool {
				return saxIsContinuation(b) && !(b0 == 0xF0 && b < 0x90) && !(b0 == 0xF4 && b >= 0x90)
			}) && more(saxIsContinuation) && more(saxIsContinuation) {
				cp := rune(b0&0x07)<<18 | rune(seq[1]&0x3F)<<12 | rune(seq[2]&0x3F)<<6 | rune(seq[3]&0x3F)
				cp -= 0x10000
				units[0], units[1] = uint16(0xD800+cp>>10), uint16(0xDC00+cp&0x3FF)
				size = 2
			}
		}
		if in.count-start+size > n || eof && in.count > start {
			// A sequence cut by the end of the input becomes U+FFFD in the
			// next read.
			in.ahead = append(append([]byte(nil), seq[:k]...), in.ahead...)
			break
		}
		for _, u := range units[:size] {
			in.put(u)
		}
	}
	return in.count - start
}

func saxIsContinuation(b byte) bool { return b&0xC0 == 0x80 }

// readLenientUTF16 reads up to n chars with a lenient UTF-16 reader (see
// saxEncUTF16BELenient). A supplementary character that does not fit is
// left for the next read.
func (in *saxInput) readLenientUTF16(n int) int {
	if in.enc == saxEncUTF16Lenient || in.enc == saxEncUTF16LELenientBOM {
		order := saxEncUTF16BELenient
		if in.enc == saxEncUTF16LELenientBOM {
			order = saxEncUTF16LELenient
		}
		b0, ok0 := in.readByte()
		b1, ok1 := in.readByte()
		switch {
		case ok1 && b0 == 0xFE && b1 == 0xFF:
			order = saxEncUTF16BELenient
		case ok1 && b0 == 0xFF && b1 == 0xFE:
			order = saxEncUTF16LELenient
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
		if in.enc == saxEncUTF16LELenient {
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
			in.ahead = append(append([]byte(nil), got...), in.ahead...)
			break
		}
		if !whole || saxIsLowSurrogate(int(u)) {
			in.put(0xFFFD)
			continue
		}
		if !saxIsHighSurrogate(int(u)) {
			in.put(u)
			continue
		}
		low, whole := unit()
		if !whole && in.count > start {
			in.ahead = append(append([]byte(nil), got...), in.ahead...)
			break
		}
		if !whole || !saxIsLowSurrogate(int(low)) {
			in.put(0xFFFD)
			continue
		}
		if in.count-start+2 > n {
			in.ahead = append(append([]byte(nil), got...), in.ahead...)
			break
		}
		in.put(u)
		in.put(low)
	}
	return in.count - start
}

// readSize is the number of bytes the next read asks for, given the number
// of chars kept in the buffer.
func (in *saxInput) readSize(kept int) int {
	switch {
	case in.firstRead:
		in.firstRead = false
		in.splitRead = true
		return saxFirstReadBytes - in.bomLen
	case in.declMode:
		return 1
	}
	return saxBufferChars - kept
}

// load moves buf[keepFrom:count] to the front of the buffer and reads more
// after it. It returns the number of chars read (0 at the end of input).
func (in *saxInput) load(keepFrom int) int {
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
func (in *saxInput) loadMore() int {
	if in.count < saxBufferChars && !in.firstRead {
		if in.declMode {
			return in.readChars(1)
		}
		return in.readChars(saxBufferChars - in.count)
	}
	return in.load(in.pos)
}

// The position helpers below keep line and column in step with pos. A
// newline is counted when the scanner consumes it.

func (in *saxInput) advance(n int) {
	in.pos += n
	in.column += n
}

func (in *saxInput) newline(n int) {
	in.pos += n
	in.line++
	in.column = 1
}

// saxEncodingByName maps an encoding name from an XML declaration to a
// decoder. Names are compared case-insensitively, except that only
// "UTF-16BE" and "UTF-16LE" in upper case get the JDK's own UTF-16 readers
// (other spellings get the lenient ones). Only the names below are
// emulated, each checked black-box for the reader the JDK picks for it; the
// JDK knows many more.
func saxEncodingByName(name string) (saxEncoding, bool) {
	switch name {
	case "UTF-16BE":
		return saxEncUTF16BE, true
	case "UTF-16LE":
		return saxEncUTF16LE, true
	}
	switch strings.ToUpper(name) {
	case "UTF-8":
		return saxEncUTF8, true
	case "UTF8":
		return saxEncUTF8Lenient, true
	case "UTF-16":
		return saxEncUTF16BE, true
	case "UTF16", "UTF_16", "UNICODE", "UNICODEBIG":
		return saxEncUTF16Lenient, true
	case "UNICODELITTLE", "X-UTF-16LE-BOM":
		return saxEncUTF16LELenientBOM, true
	case "UTF-16BE", "UTF_16BE", "X-UTF-16BE", "UNICODEBIGUNMARKED":
		return saxEncUTF16BELenient, true
	case "UTF-16LE", "UTF_16LE", "X-UTF-16LE", "UNICODELITTLEUNMARKED":
		return saxEncUTF16LELenient, true
	case "ISO-8859-1", "ISO8859-1", "ISO_8859-1", "ISO8859_1", "ISO_8859_1", "ISO-IR-100", "LATIN1", "L1",
		"CP819", "IBM819", "IBM-819", "CSISOLATIN1":
		return saxEncLatin1, true
	case "US-ASCII", "ASCII", "US", "ISO646-US", "IBM367", "IBM-367", "CP367", "CSASCII", "ANSI_X3.4-1968",
		"ANSI_X3.4-1986", "ISO-IR-6":
		return saxEncASCII, true
	case "ASCII7":
		return saxEncASCIILenient, true
	}
	return 0, false
}

// saxUnsupportedEncodingException is java.io.UnsupportedEncodingException.
type saxUnsupportedEncodingException struct{ jutil.IOException }

// saxKnownEncoding reports whether name is a common encoding (IANA name)
// that the JDK supports but this parser does not emulate; any other name
// is treated as unknown to the JDK.
func saxKnownEncoding(name string) bool {
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
