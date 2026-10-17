package jutil

import (
	"io"
	"os"
)

// Go-only: language-adapter — java.io.PrintStream over a Go io.Writer,
// encoding with UTF-8 (unpaired surrogates become '?', as the JDK encoder
// replaces them) and with "\n" as the line separator. Like PrintStream,
// it never throws: a write error sets the flag that CheckError reports.

// PrintStream is java.io.PrintStream.
type PrintStream struct {
	out           io.Writer
	errorOccurred bool
}

// NewPrintStream is new PrintStream(OutputStream).
func NewPrintStream(out io.Writer) *PrintStream {
	return &PrintStream{out: out}
}

func (p *PrintStream) write(s string) {
	if _, err := p.out.Write(StringGetBytesUTF8(s)); err != nil {
		p.errorOccurred = true
	}
}

// Print is PrintStream.print(String).
func (p *PrintStream) Print(s string) { p.write(s) }

// PrintNullable is PrintStream.print(String) for a possibly-null String
// ("null" when nil).
func (p *PrintStream) PrintNullable(s *string) {
	if s == nil {
		p.write("null")
		return
	}
	p.write(*s)
}

// Println is PrintStream.println(String).
func (p *PrintStream) Println(s string) { p.write(s + "\n") }

// PrintlnNullable is PrintStream.println(String) for a possibly-null String.
func (p *PrintStream) PrintlnNullable(s *string) {
	p.PrintNullable(s)
	p.write("\n")
}

// PrintlnEmpty is PrintStream.println().
func (p *PrintStream) PrintlnEmpty() { p.write("\n") }

// Printf is PrintStream.printf(String, Object...) (String.format).
func (p *PrintStream) Printf(format string, args ...any) *PrintStream {
	p.write(StringFormat(format, args...))
	return p
}

// Append is PrintStream.append(CharSequence).
func (p *PrintStream) Append(s string) *PrintStream {
	p.write(s)
	return p
}

// AppendCharSequence is Appendable.append(CharSequence).
func (p *PrintStream) AppendCharSequence(csq string) { p.write(csq) }

// AppendCharacter is Appendable.append(char).
func (p *PrintStream) AppendCharacter(c uint16) { p.write(StringFromChar(c)) }

// Flush is PrintStream.flush().
func (p *PrintStream) Flush() {
	if f, ok := p.out.(interface{ Flush() error }); ok {
		if err := f.Flush(); err != nil {
			p.errorOccurred = true
		}
	}
}

// CheckError is PrintStream.checkError().
func (p *PrintStream) CheckError() bool {
	p.Flush()
	return p.errorOccurred
}

// Write is PrintStream.write(byte[]) (which PrintStream inherits from
// FilterOutputStream and implements with write(byte[], int, int)): the
// bytes are written unchanged and, as PrintStream never throws, an I/O error
// only sets the flag that CheckError reports. It makes a *PrintStream an
// io.Writer, for code taking an OutputStream.
func (p *PrintStream) Write(b []byte) (int, error) {
	if _, err := p.out.Write(b); err != nil {
		p.errorOccurred = true
	}
	return len(b), nil
}

// NewPrintStreamFromFile is new PrintStream(File): it creates (or
// truncates) the file, and panics with FileNotFoundException when it cannot.
func NewPrintStreamFromFile(file *File) *PrintStream {
	f, err := os.OpenFile(file.GetPath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		panic(FileNotFoundExceptionFromError(err, file.GetPath()))
	}
	return NewPrintStream(f)
}

// Close is PrintStream.close(): flushes the stream and closes the
// underlying stream, if it can be closed. Like PrintStream, it never
// throws.
func (p *PrintStream) Close() {
	p.Flush()
	if c, ok := p.out.(io.Closer); ok {
		if err := c.Close(); err != nil {
			p.errorOccurred = true
		}
	}
}
