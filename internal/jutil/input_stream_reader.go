package jutil

import (
	"io"
	"strings"
	"unicode/utf8"
	"unsafe"
)

// Go-only: language-adapter — java.io.InputStreamReader as a JavaReader.
// The whole stream is read and decoded on the first read (Java decodes
// incrementally; the characters are the same). Malformed input is replaced
// (InputStreamReader's decoder uses CodingErrorAction.REPLACE), unless the
// reader was created for Files.newBufferedReader, which reports it.

// InputStreamReader is java.io.InputStreamReader.
type InputStreamReader struct {
	in     io.Reader
	cs     *Charset
	strict bool
	text   *StringReader
	closed bool
}

// NewInputStreamReader is new InputStreamReader(in, charset).
func NewInputStreamReader(in io.Reader, cs *Charset) *InputStreamReader {
	return &InputStreamReader{in: in, cs: CheckNotNull(cs)}
}

func (r *InputStreamReader) load() {
	if r.text != nil {
		return
	}
	r.text = NewStringReader(r.decodeAll())
}

// decodeAll reads and decodes the whole stream.
func (r *InputStreamReader) decodeAll() string {
	b, err := io.ReadAll(r.in)
	if err != nil {
		panic(IOExceptionFromError(err, ""))
	}
	if (r.cs == nil || r.cs == StandardCharsets_UTF_8) && utf8.Valid(b) {
		return unsafe.String(unsafe.SliceData(b), len(b))
	}
	if r.strict {
		return r.cs.DecodeStrict(b)
	}
	return r.cs.Decode(b)
}

// Read is Reader.read(char[], int, int).
func (r *InputStreamReader) Read(cbuf []uint16, off, length int) int {
	if r.closed {
		panic(NewIOException("Stream closed"))
	}
	r.load()
	return r.text.Read(cbuf, off, length)
}

// Close is Reader.close().
func (r *InputStreamReader) Close() {
	r.closed = true
	if c, ok := r.in.(io.Closer); ok {
		c.Close()
	}
}

// CharStreams_ToString is Guava CharStreams.toString(Readable): reads all
// characters from r (it does not close it).
func CharStreams_ToString(r JavaReader) string {
	// Go-only: performance — an open InputStreamReader that nothing has read
	// yet decodes its whole stream at once, and that text is the result:
	// return it without the round trip through UTF-16 chunks (valid UTF-8
	// comes back from that round trip unchanged), and leave the reader at
	// its end, as reading it through would.
	// The same goes for a StringReader nothing has read yet.
	if isr, ok := r.(*InputStreamReader); ok && !isr.closed && isr.text == nil {
		s := isr.decodeAll()
		isr.text = &StringReader{converted: true}
		return charStreamsRoundTrip(s)
	}
	if sr, ok := r.(*StringReader); ok && !sr.converted {
		s := sr.s
		*sr = StringReader{converted: true}
		return charStreamsRoundTrip(s)
	}
	var units []uint16
	buf := make([]uint16, 2048)
	for {
		n := r.Read(buf, 0, len(buf))
		if n == -1 {
			break
		}
		units = append(units, buf[:n]...)
	}
	return StringFromUTF16(units)
}

// charStreamsRoundTrip is StringFromUTF16(StringToUTF16(s)), which is s
// itself when s is valid UTF-8.
func charStreamsRoundTrip(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	return StringFromUTF16(StringToUTF16(s))
}

// IOExceptionFromError is the IOException the JDK throws for a failed file
// operation on path: NoSuchFileException / AccessDeniedException (message:
// the path) or an IOException with the system's message.
func IOExceptionFromError(err error, path string) Throwable {
	e := &IOException{}
	msg := err.Error()
	switch {
	case path != "" && strings.Contains(msg, "no such file or directory"):
		InitThrowable(e, "java.nio.file.NoSuchFileException", &path, nil)
	case path != "" && strings.Contains(msg, "permission denied"):
		InitThrowable(e, "java.nio.file.AccessDeniedException", &path, nil)
	case strings.Contains(msg, "is a directory"):
		m := "Is a directory"
		InitThrowable(e, "java.io.IOException", &m, nil)
	default:
		InitThrowable(e, "java.io.IOException", &msg, nil)
	}
	return e
}
