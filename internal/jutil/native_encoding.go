package jutil

import (
	"bytes"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

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
		cs = StandardCharsets_ISO_8859_1
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
