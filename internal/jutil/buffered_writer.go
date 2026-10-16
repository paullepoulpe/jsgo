package jutil

import "errors"

// Go-only: language-adapter — java.io.BufferedWriter (see its Javadoc and
// that of java.io.Writer) over a jutil.Writer, with the default buffer of
// 8192 chars.
//
// Characters are collected in the buffer and handed to the underlying
// writer only when the buffer is full, on Flush and on Close.
//
// Observed behaviour (black-box, JDK 21, with a Writer that logs its calls):
//   - the buffer is passed on as soon as it is full, so writing exactly 8192
//     chars into an empty buffer reaches the underlying writer at once, and
//     a write that does not fit is split at the buffer boundary (which may
//     fall between the two halves of a surrogate pair);
//   - Flush passes the buffered chars on (if any) and then flushes the
//     underlying writer;
//   - Close passes the buffered chars on and closes the underlying writer
//     without flushing it. The underlying writer is closed, and this one
//     marked closed, even when passing the chars on fails; that failure is
//     then rethrown. A second Close does nothing;
//   - after Close, every write, append and flush (even of an empty string)
//     throws IOException("Stream closed").

// bufferedWriterSize is the default buffer size of BufferedWriter, in chars.
const bufferedWriterSize = 8192

// BufferedWriter is java.io.BufferedWriter.
type BufferedWriter struct {
	out    Writer
	buf    []uint16
	closed bool
}

// NewBufferedWriter is new BufferedWriter(out).
func NewBufferedWriter(out Writer) *BufferedWriter {
	if IsNil(out) {
		// Observed behaviour (black-box): the JDK 21 constructor fails with
		// this helpful NullPointerException message.
		panic(NewNullPointerExceptionMsg(`Cannot invoke "Object.getClass()" because "writer" is null`))
	}
	return &BufferedWriter{out: out, buf: make([]uint16, 0, bufferedWriterSize)}
}

// checkNotClosed throws IOException "Stream closed" once Close has run.
func (w *BufferedWriter) checkNotClosed() {
	if w.closed {
		panic(NewIOException("Stream closed"))
	}
}

// FlushBuffer passes the buffered chars on to the underlying writer without
// flushing it (BufferedWriter.flushBuffer()).
func (w *BufferedWriter) FlushBuffer() {
	w.checkNotClosed()
	if len(w.buf) == 0 {
		return
	}
	chunk := StringFromUTF16(w.buf)
	w.buf = w.buf[:0]
	w.out.Write(chunk)
}

// Write is BufferedWriter.write(String).
func (w *BufferedWriter) Write(s string) {
	w.checkNotClosed()
	w.writeUnits(StringToUTF16(s))
}

func (w *BufferedWriter) writeUnits(units []uint16) {
	for len(units) > 0 {
		n := min(bufferedWriterSize-len(w.buf), len(units))
		w.buf = append(w.buf, units[:n]...)
		units = units[n:]
		if len(w.buf) >= bufferedWriterSize {
			w.FlushBuffer()
		}
	}
}

// Append is BufferedWriter.append(CharSequence), which writes the string.
func (w *BufferedWriter) Append(csq string) *BufferedWriter {
	w.Write(csq)
	return w
}

// AppendChar is BufferedWriter.append(char).
func (w *BufferedWriter) AppendChar(c uint16) *BufferedWriter {
	w.checkNotClosed()
	w.writeUnits([]uint16{c})
	return w
}

// WriteString implements io.StringWriter for Go callers: it is Write, with
// an IOException returned as an error instead of thrown.
func (w *BufferedWriter) WriteString(s string) (n int, err error) {
	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(*IOException)
			if !ok {
				panic(r)
			}
			n, err = 0, errors.New(e.GetMessage())
		}
	}()
	w.Write(s)
	return len(s), nil
}

// Flush is BufferedWriter.flush().
func (w *BufferedWriter) Flush() {
	w.FlushBuffer()
	w.out.Flush()
}

// Close is BufferedWriter.close().
func (w *BufferedWriter) Close() {
	if w.closed {
		return
	}
	var flushFailure any
	func() {
		defer func() { flushFailure = recover() }()
		w.FlushBuffer()
	}()
	w.closed = true
	w.buf = nil
	if flushFailure != nil {
		// The failure to pass the chars on wins over one from closing (Java
		// adds the latter as a suppressed exception).
		func() {
			defer func() { _ = recover() }()
			w.out.Close()
		}()
		panic(flushFailure)
	}
	w.out.Close()
}
