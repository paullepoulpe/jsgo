package jutil

import "strings"

// Go-only: language-adapter — java.io.StringWriter, written from its
// Javadoc: a character stream that collects its output in a string buffer.
// Flushing and closing have no effect, and the buffer can still be used
// after the writer is closed.

// StringWriter is java.io.StringWriter.
type StringWriter struct {
	buf strings.Builder
}

// NewStringWriter is new StringWriter().
func NewStringWriter() *StringWriter { return &StringWriter{} }

// Write is StringWriter.write(String).
func (w *StringWriter) Write(s string) { w.buf.WriteString(s) }

// Flush is StringWriter.flush().
func (w *StringWriter) Flush() {}

// Close is StringWriter.close().
func (w *StringWriter) Close() {}

// String is StringWriter.toString(): the buffer's current value.
func (w *StringWriter) String() string { return w.buf.String() }

// AppendCharSequence is StringWriter.append(CharSequence) (Appendable).
func (w *StringWriter) AppendCharSequence(csq string) { w.buf.WriteString(csq) }

// AppendCharacter is StringWriter.append(char) (Appendable).
func (w *StringWriter) AppendCharacter(c uint16) { w.buf.WriteString(StringFromChar(c)) }

// WriteString is io.StringWriter (Writer.write(String)), for callers that
// take a Go io.StringWriter.
func (w *StringWriter) WriteString(s string) (int, error) {
	w.buf.WriteString(s)
	return len(s), nil
}
