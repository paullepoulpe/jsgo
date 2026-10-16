package jutil

import (
	"io"
	"strings"
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
	b, err := io.ReadAll(r.in)
	if err != nil {
		panic(IOExceptionFromError(err, ""))
	}
	var s string
	if r.strict {
		s = r.cs.DecodeStrict(b)
	} else {
		s = r.cs.Decode(b)
	}
	r.text = NewStringReader(s)
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
