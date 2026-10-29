package jutil

import (
	"bytes"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// Go-only: language-adapter — java.io.PrintStream over a Go io.Writer,
// encoding with UTF-8 (unpaired surrogates become '?', as the JDK encoder
// replaces them) and with "\n" as the line separator. Like PrintStream,
// it never throws on I/O errors.

// PrintStream is java.io.PrintStream.
type PrintStream struct {
	out io.Writer
}

// NewPrintStream is new PrintStream(OutputStream).
func NewPrintStream(out io.Writer) *PrintStream {
	return &PrintStream{out: out}
}

func (p *PrintStream) write(s string) {
	if utf8.ValidString(s) {
		_, _ = io.WriteString(p.out, s)
		return
	}
	_, _ = p.out.Write(StringGetBytesUTF8(s))
}

// Print is PrintStream.print(String).
func (p *PrintStream) Print(s string) { p.write(s) }

// Println is PrintStream.println(String).
func (p *PrintStream) Println(s string) {
	if p == nil || p.out == nil {
		return
	}
	if len(s) < 512 && utf8.ValidString(s) {
		var buf [512]byte
		copy(buf[:], s)
		buf[len(s)] = '\n'
		_, _ = p.out.Write(buf[:len(s)+1])
		return
	}
	p.write(s)
	p.write("\n")
}

// PrintlnNullable is PrintStream.println(String) for a possibly-null String.
func (p *PrintStream) PrintlnNullable(s *string) {
	if s == nil {
		p.Println("null")
		return
	}
	p.Println(*s)
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
		_ = f.Flush()
	}
}

// Write is PrintStream.write(byte[]) (which PrintStream inherits from
// FilterOutputStream and implements with write(byte[], int, int)): the
// bytes are written unchanged and, as PrintStream never throws, I/O errors
// are ignored. It makes a *PrintStream an io.Writer, for code taking an
// OutputStream.
func (p *PrintStream) Write(b []byte) (int, error) {
	_, _ = p.out.Write(b)
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
		_ = c.Close()
	}
}

// Go-only: language-adapter — the JVM's native.encoding, which is also
// stdout.encoding and stderr.encoding when the streams are not a console
// (the System.out and System.err PrintStreams encode text with it). On
// Linux the JVM takes it from the C library: setlocale(LC_ALL, "") then
// nl_langinfo(CODESET), and the C library falls back to the "C" locale
// (ASCII, which Java calls US-ASCII) when the locale named by LC_ALL,
// LC_CTYPE or LANG (the first one set) is not installed. Installed locales
// are the built-in C.* ones, the directories of /usr/lib/locale and the
// entries of /usr/lib/locale/locale-archive. Only the UTF-8 and ISO-8859-1
// codesets are mapped: how the JVM names other codesets (EUC-JP, for
// one, may become a charset jutil lacks) is not checked here, so any other
// codeset is taken as US-ASCII.

// System_NativeCharset is Charset.forName(System.getProperty(
// "native.encoding")) for the process's locale.
func System_NativeCharset() *Charset {
	name := ""
	for _, v := range []string{"LC_ALL", "LC_CTYPE", "LANG"} {
		if name = os.Getenv(v); name != "" {
			break
		}
	}
	if name == "" || name == "C" || name == "POSIX" {
		return StandardCharsets_US_ASCII
	}
	if at := strings.IndexByte(name, '@'); at >= 0 {
		name = name[:at]
	}
	dot := strings.IndexByte(name, '.')
	if dot < 0 {
		return StandardCharsets_US_ASCII
	}
	lang, codeset := name[:dot], nativeEncoding_normalizeCodeset(name[dot+1:])
	var cs *Charset
	switch codeset {
	case "utf8":
		cs = StandardCharsets_UTF_8
	case "iso88591":
		cs = standardCharsets_ISO_8859_1
	default:
		return StandardCharsets_US_ASCII
	}
	if lang == "C" || nativeEncoding_localeInstalled(lang+"."+codeset) {
		return cs
	}
	return StandardCharsets_US_ASCII
}

// nativeEncoding_normalizeCodeset is glibc's codeset normalization: ASCII
// letters lowercased, digits kept, everything else dropped.
func nativeEncoding_normalizeCodeset(codeset string) string {
	var sb strings.Builder
	for i := 0; i < len(codeset); i++ {
		c := codeset[i]
		switch {
		case c >= 'A' && c <= 'Z':
			sb.WriteByte(c + 'a' - 'A')
		case (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9'):
			sb.WriteByte(c)
		}
	}
	return sb.String()
}

func nativeEncoding_localeInstalled(normalized string) bool {
	if st, err := os.Stat("/usr/lib/locale/" + normalized); err == nil && st.IsDir() {
		return true
	}
	archive, err := os.ReadFile("/usr/lib/locale/locale-archive")
	return err == nil && bytes.Contains(archive, append([]byte(normalized), 0))
}

// NewPrintStreamWithCharset is new PrintStream(out, true, cs), for cs
// US-ASCII, ISO-8859-1 or UTF-8 and a stream that is only printed to
// (print/println of strings): the text is encoded with cs, each character
// cs cannot encode (a surrogate pair counts as one) becoming '?'. Bytes
// written with PrintStream.write(byte[]) are encoded the same way, so such
// a stream must not be used for raw bytes unless cs is UTF-8.
func NewPrintStreamWithCharset(out io.Writer, cs *Charset) *PrintStream {
	if cs == StandardCharsets_UTF_8 {
		return NewPrintStream(out)
	}
	return NewPrintStream(&nativeEncoding_writer{out: out, enc: cs.impl().newEncoder()})
}

// nativeEncoding_writer re-encodes the UTF-8 that PrintStream writes (in
// which unpaired surrogates are already '?').
type nativeEncoding_writer struct {
	out     io.Writer
	enc     charsetEncoder
	pending []byte
}

func (w *nativeEncoding_writer) Write(b []byte) (int, error) {
	in := append(w.pending, b...)
	w.pending = nil
	out := make([]byte, 0, len(in))
	for len(in) > 0 {
		r, size := utf8.DecodeRune(in)
		if r == utf8.RuneError && size <= 1 && !utf8.FullRune(in) {
			w.pending = append([]byte{}, in...)
			break
		}
		in = in[size:]
		if r == utf8.RuneError && size <= 1 {
			out = append(out, '?')
		} else {
			out = w.enc.encodeRune(out, r)
		}
	}
	if _, err := w.out.Write(out); err != nil {
		return 0, err
	}
	return len(b), nil
}

func (w *nativeEncoding_writer) Flush() error {
	if f, ok := w.out.(interface{ Flush() error }); ok {
		return f.Flush()
	}
	return nil
}

func (w *nativeEncoding_writer) Close() error {
	if c, ok := w.out.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
