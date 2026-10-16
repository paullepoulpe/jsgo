package jutil

// Go-only: language-adapter — java.io.PrintWriter over a Writer, without
// automatic flushing, written from its Javadoc: the print methods write the
// string form of their argument, println terminates the line with the line
// separator ("\n" on the platforms the port runs on), and no method throws
// an IOException: an I/O error of the underlying writer only sets the flag
// that CheckError reports. Once the writer is closed, closing it again has
// no effect and every write only sets that flag.

// PrintWriter is java.io.PrintWriter.
type PrintWriter struct {
	out           Writer // nil once closed
	errorOccurred bool
}

// NewPrintWriter is new PrintWriter(Writer out); it throws a
// NullPointerException when out is null.
func NewPrintWriter(out Writer) *PrintWriter {
	if IsNil(out) {
		panic(NewNullPointerException())
	}
	return &PrintWriter{out: out}
}

func (w *PrintWriter) guard(f func()) {
	if w.out == nil {
		w.errorOccurred = true
		return
	}
	Try(f, func(IIOException) { w.errorOccurred = true })
}

// Write is PrintWriter.write(String).
func (w *PrintWriter) Write(s string) {
	w.guard(func() { w.out.Write(s) })
}

// Print is PrintWriter.print(String).
func (w *PrintWriter) Print(s string) { w.Write(s) }

// PrintNullable is PrintWriter.print(String) for a String that may be null
// (printed as "null").
func (w *PrintWriter) PrintNullable(s *string) {
	if s == nil {
		w.Write("null")
		return
	}
	w.Write(*s)
}

// PrintChar is PrintWriter.print(char).
func (w *PrintWriter) PrintChar(c uint16) { w.Write(StringFromChar(c)) }

// PrintlnEmpty is PrintWriter.println().
func (w *PrintWriter) PrintlnEmpty() { w.Write("\n") }

// Println is PrintWriter.println(String).
func (w *PrintWriter) Println(s string) {
	w.Print(s)
	w.PrintlnEmpty()
}

// Flush is PrintWriter.flush().
func (w *PrintWriter) Flush() {
	w.guard(func() { w.out.Flush() })
}

// Close is PrintWriter.close().
func (w *PrintWriter) Close() {
	if w.out == nil {
		return
	}
	w.guard(func() { w.out.Close() })
	w.out = nil
}

// CheckError is PrintWriter.checkError(): it flushes the stream if it is not
// closed and reports whether an error occurred on it.
func (w *PrintWriter) CheckError() bool {
	if w.out != nil {
		w.Flush()
	}
	return w.errorOccurred
}
