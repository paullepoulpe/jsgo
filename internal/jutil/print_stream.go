package jutil

import "io"

// Go-only: language-adapter — java.io.PrintStream over a Go io.Writer,
// encoding with UTF-8 (unpaired surrogates become '?', as the JDK encoder
// replaces them) and with "\n" as the line separator. Like PrintStream,
// it never throws: a write error sets the flag that CheckError reports.

// PrintStream is java.io.PrintStream.
type PrintStream struct {
	out     io.Writer
	trouble bool
}

// NewPrintStream is new PrintStream(OutputStream).
func NewPrintStream(out io.Writer) *PrintStream {
	return &PrintStream{out: out}
}

func (p *PrintStream) write(s string) {
	if _, err := p.out.Write(StringGetBytesUTF8(s)); err != nil {
		p.trouble = true
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
			p.trouble = true
		}
	}
}

// CheckError is PrintStream.checkError().
func (p *PrintStream) CheckError() bool {
	p.Flush()
	return p.trouble
}
