package jutil

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// Go-only: language-adapter — java.lang.Appendable. An IOException is
// thrown as a panic (jutil.IOException), like every other Java exception.

// Appendable is java.lang.Appendable.
type Appendable interface {
	// AppendCharSequence is Appendable.append(CharSequence).
	AppendCharSequence(csq string)
	// AppendCharacter is Appendable.append(char).
	AppendCharacter(c uint16)
}

// Go-only: language-adapter — java.io.Writer. An IOException is thrown as
// a panic (jutil.IOException), like every other Java exception.

// Writer is java.io.Writer (the methods ported code uses).
type Writer interface {
	// Write is Writer.write(String).
	Write(s string)
	// TryWrite writes s and returns an IIOException instead of panicking on I/O failure.
	TryWrite(s string) IIOException
	// Flush is Writer.flush().
	Flush()
	// TryFlush flushes the writer and returns an IIOException instead of panicking on I/O failure.
	TryFlush() IIOException
	// Close is Writer.close().
	Close()
}

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

// TryWrite writes s to the buffer (never fails).
func (w *StringWriter) TryWrite(s string) IIOException {
	w.buf.WriteString(s)
	return nil
}

// Flush is StringWriter.flush().
func (w *StringWriter) Flush() {}

// TryFlush is a no-op on StringWriter (never fails).
func (w *StringWriter) TryFlush() IIOException { return nil }

// Close is StringWriter.close().
func (w *StringWriter) Close() {}

// String is StringWriter.toString(): the buffer's current value.
func (w *StringWriter) String() string { return w.buf.String() }

// AppendCharSequence is StringWriter.append(CharSequence) (Appendable).
func (w *StringWriter) AppendCharSequence(csq string) { w.buf.WriteString(csq) }

// AppendCharacter is StringWriter.append(char) (Appendable).
func (w *StringWriter) AppendCharacter(c uint16) {
	if c < 0x80 {
		w.buf.WriteByte(byte(c))
		return
	}
	w.buf.WriteString(StringFromChar(c))
}

// WriteString is io.StringWriter (Writer.write(String)), for callers that
// take a Go io.StringWriter.
func (w *StringWriter) WriteString(s string) (int, error) {
	w.buf.WriteString(s)
	return len(s), nil
}

// Go-only: language-adapter — java.io.BufferedReader over a StringReader or
// over an InputStreamReader(stream, UTF_8). The whole input is decoded on
// the first read (Java decodes incrementally; the characters are the same).

// BufferedReader is java.io.BufferedReader.
type BufferedReader struct {
	src    io.Reader
	text   string
	loaded bool
	pos    int
}

// NewBufferedReaderFromString is new BufferedReader(new StringReader(s)).
func NewBufferedReaderFromString(s string) *BufferedReader {
	return &BufferedReader{text: s, loaded: true}
}

// NewBufferedReaderUTF8 is new BufferedReader(new InputStreamReader(stream,
// UTF_8)): malformed input is replaced like new String(bytes, UTF_8).
func NewBufferedReaderUTF8(stream io.Reader) *BufferedReader {
	return &BufferedReader{src: stream}
}

func (r *BufferedReader) load() {
	if r.loaded {
		return
	}
	b, err := io.ReadAll(r.src)
	if err != nil {
		panic(NewIOException(err.Error()))
	}
	r.text = NewStringFromUTF8Bytes(b)
	r.loaded = true
}

// NextLineFromString scans the next line starting at pos in text, ending at
// '\n', '\r' or "\r\n" (not included). Returns ok == false when pos >= len(text).
func NextLineFromString(text string, pos int) (line string, nextPos int, ok bool) {
	if pos >= len(text) {
		return "", pos, false
	}
	rest := text[pos:]
	nl := strings.IndexByte(rest, '\n')
	if nl >= 0 {
		cr := strings.IndexByte(rest[:nl], '\r')
		if cr < 0 {
			return rest[:nl], pos + nl + 1, true
		}
		if cr == nl-1 {
			return rest[:cr], pos + nl + 1, true
		}
		return rest[:cr], pos + cr + 1, true
	}
	cr := strings.IndexByte(rest, '\r')
	if cr < 0 {
		return rest, len(text), true
	}
	return rest[:cr], pos + cr + 1, true
}

// ReadLine is BufferedReader.readLine(): a line ends at '\n', '\r' or
// "\r\n" (not included); ok is false (Java null) at end of input.
func (r *BufferedReader) ReadLine() (line string, ok bool) {
	r.load()
	line, r.pos, ok = NextLineFromString(r.text, r.pos)
	return line, ok
}

// Close is BufferedReader.close().
func (r *BufferedReader) Close() {}

// Go-only: language-adapter — Guava com.google.common.io.CharSource, as
// returned by Files.asCharSource(File, Charset).

// CharSource is Guava's CharSource over a file.
type CharSource struct {
	file    *File
	charset *Charset
}

// Files_AsCharSource is Guava Files.asCharSource(File, Charset).
func Files_AsCharSource(file *File, charset *Charset) *CharSource {
	return &CharSource{file: CheckNotNull(file), charset: CheckNotNull(charset)}
}

// OpenStream is CharSource.openStream(): a new InputStreamReader over a new
// FileInputStream for the file. It panics with FileNotFoundException when
// the file cannot be opened.
func (c *CharSource) OpenStream() JavaReader {
	f, err := os.Open(c.file.GetPath())
	if err == nil {
		var st os.FileInfo
		if st, err = f.Stat(); err == nil && st.IsDir() {
			f.Close()
			panic(NewFileNotFoundException(c.file.GetPath() + " (Is a directory)"))
		}
	}
	if err != nil {
		panic(FileNotFoundExceptionFromError(err, c.file.GetPath()))
	}
	return NewInputStreamReader(f, c.charset)
}

// Read is CharSource.read(): the whole content of the source as a string
// (for a file, decoded with the source's charset; it panics with an
// IOException (FileNotFoundException) when the file cannot be read).
func (c *CharSource) Read() string {
	r := c.OpenStream()
	defer r.Close()
	return CharStreams_ToString(r)
}

// TryRead reads the whole content of the source as a string, returning an
// IIOException instead of panicking if opening or reading the file fails.
func (c *CharSource) TryRead() (string, IIOException) {
	f, err := os.Open(c.file.GetPath())
	if err == nil {
		var st os.FileInfo
		if st, err = f.Stat(); err == nil && st.IsDir() {
			f.Close()
			return "", NewFileNotFoundException(c.file.GetPath() + " (Is a directory)")
		}
	}
	if err != nil {
		return "", FileNotFoundExceptionFromError(err, c.file.GetPath())
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return "", IOExceptionFromError(err, "").(IIOException)
	}
	if (c.charset == nil || c.charset == StandardCharsets_UTF_8) && utf8.Valid(b) {
		return string(b), nil
	}
	return charStreamsRoundTrip(c.charset.Decode(b)), nil
}

// TryReadLines reads and splits the lines of the source like
// CharStreams.readLines(openStream()), returning an IIOException instead of
// panicking if opening or reading the file fails.
func (c *CharSource) TryReadLines() ([]string, IIOException) {
	s, ex := c.TryRead()
	if ex != nil {
		return nil, ex
	}
	return splitLinesUTF16(StringToUTF16(s)), nil
}

// CharStreams_ReadLines is Guava CharStreams.readLines(Readable): the lines
// of the input, split at "\n", "\r" and "\r\n" (Guava's LineReader). The
// line terminators are not included, and an input that ends with a
// terminator has no empty last line.
func CharStreams_ReadLines(r JavaReader) []string {
	return splitLinesUTF16(StringToUTF16(CharStreams_ToString(r)))
}

func splitLinesUTF16(units []uint16) []string {
	var result []string
	start := 0
	for i := 0; i < len(units); i++ {
		switch units[i] {
		case '\n':
			result = append(result, StringFromUTF16(units[start:i]))
			start = i + 1
		case '\r':
			result = append(result, StringFromUTF16(units[start:i]))
			if i+1 < len(units) && units[i+1] == '\n' {
				i++
			}
			start = i + 1
		}
	}
	if start < len(units) {
		result = append(result, StringFromUTF16(units[start:]))
	}
	return result
}

// Go-only: language-adapter — java.io.FileInputStream as a Go
// io.ReadCloser over an *os.File.

// NewFileInputStream is new FileInputStream(String name): it throws
// FileNotFoundException("<name> (<reason>)") when the file does not exist,
// is a directory, or cannot be opened for reading.
func NewFileInputStream(name string) io.ReadCloser {
	f, err := TryNewFileInputStream(name)
	if err != nil {
		panic(err)
	}
	return f
}

// TryNewFileInputStream opens name like NewFileInputStream, returning an
// IIOException instead of panicking on failure.
func TryNewFileInputStream(name string) (io.ReadCloser, IIOException) {
	f, err := os.Open(name)
	if err == nil {
		if st, serr := f.Stat(); serr == nil && st.IsDir() {
			f.Close()
			return nil, NewFileNotFoundException(name + " (Is a directory)")
		}
		return f, nil
	}
	return nil, FileNotFoundExceptionFromError(err, name)
}

// Go-only: language-adapter — java.io.FileOutputStream as a Go
// io.WriteCloser over an *os.File.

// NewFileOutputStream is new FileOutputStream(String name): it creates or
// truncates the file, and throws FileNotFoundException("<name> (<reason>)")
// when the file cannot be opened for writing. A write error is thrown as
// IOException by the writers that wrap it.
func NewFileOutputStream(name string) io.WriteCloser {
	f, err := TryNewFileOutputStream(name)
	if err != nil {
		panic(err)
	}
	return f
}

// TryNewFileOutputStream opens name for writing like NewFileOutputStream,
// returning an IIOException instead of panicking on failure.
func TryNewFileOutputStream(name string) (io.WriteCloser, IIOException) {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		return nil, FileNotFoundExceptionFromError(err, name)
	}
	return f, nil
}

// Go-only: language-adapter — the reading methods of java.nio.file.Files
// (see files.go).

// Files_ReadString is Files.readString(Path, Charset): it panics with
// MalformedInputException on malformed input, and with an IOException
// (NoSuchFileException, ...) when the file cannot be read.
func Files_ReadString(path *Path, cs *Charset) string {
	return Files_ReadStringRaw(path.String(), cs)
}

// Files_ReadStringRaw is Files_ReadString for a pre-normalized path string.
func Files_ReadStringRaw(pathStr string, cs *Charset) string {
	s, err := TryFiles_ReadStringRaw(pathStr, cs)
	if err != nil {
		panic(err)
	}
	return s
}

// TryFiles_ReadStringRaw reads pathStr like Files_ReadStringRaw, returning an
// IIOException instead of panicking on I/O or charset decoding failure.
func TryFiles_ReadStringRaw(pathStr string, cs *Charset) (string, IIOException) {
	b, err := os.ReadFile(pathStr)
	if err != nil {
		return "", IOExceptionFromError(err, pathStr).(IIOException)
	}
	if (cs == nil || cs == StandardCharsets_UTF_8) && utf8.Valid(b) {
		return string(b), nil
	}
	return cs.TryDecodeStrict(b)
}

// Files_NewBufferedReader is Files.newBufferedReader(Path, Charset): the
// file is opened now, and malformed input is reported while reading.
func Files_NewBufferedReader(path *Path, cs *Charset) JavaReader {
	return Files_NewBufferedReaderRaw(path.String(), cs)
}

// Files_NewBufferedReaderRaw is Files_NewBufferedReader for a pre-normalized path string.
func Files_NewBufferedReaderRaw(pathStr string, cs *Charset) JavaReader {
	f, err := os.Open(pathStr)
	if err != nil {
		panic(IOExceptionFromError(err, pathStr))
	}
	r := NewInputStreamReader(f, cs)
	r.strict = true
	return r
}

// Files_ReadStringUTF8 is java.nio.file.Files.readString(Path), which
// decodes with UTF-8.
func Files_ReadStringUTF8(path *Path) string {
	return Files_ReadString(path, StandardCharsets_UTF_8)
}

// TryFiles_ReadStringUTF8 reads path with UTF-8 like Files_ReadStringUTF8,
// returning an IIOException instead of panicking on failure.
func TryFiles_ReadStringUTF8(path *Path) (string, IIOException) {
	return TryFiles_ReadStringRaw(path.String(), StandardCharsets_UTF_8)
}

// Go-only: language-adapter — java.util.zip.GZIPInputStream over Go's
// compress/gzip, as an io.ReadCloser. A malformed header is thrown as
// ZipException("Not in GZIP format"), as the Javadoc specifies for the
// constructor; other read errors are thrown as IOException.

type gzipInputStream struct {
	in io.Closer
	zr *gzip.Reader
}

// NewGZIPInputStreamWithSize is new GZIPInputStream(InputStream, int size):
// it reads and checks the GZIP header immediately. The buffer size only
// affects performance.
func NewGZIPInputStreamWithSize(in io.ReadCloser, size int) io.ReadCloser {
	r, err := TryNewGZIPInputStreamWithSize(in, size)
	if err != nil {
		panic(err)
	}
	return r
}

// TryNewGZIPInputStreamWithSize wraps in like NewGZIPInputStreamWithSize,
// returning an IIOException instead of panicking on header/EOF failure.
func TryNewGZIPInputStreamWithSize(in io.ReadCloser, size int) (io.ReadCloser, IIOException) {
	zr, err := gzip.NewReader(bufio.NewReaderSize(in, size))
	if err != nil {
		if err == gzip.ErrHeader {
			return nil, NewZipException("Not in GZIP format")
		}
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			// Java: new EOFException(), with a null message.
			e := &EOFException{}
			InitThrowable(e, "java.io.EOFException", nil, nil)
			return nil, e
		}
		return nil, NewIOException(err.Error())
	}
	return &gzipInputStream{in: in, zr: zr}, nil
}

func (g *gzipInputStream) Read(b []byte) (int, error) {
	return g.zr.Read(b)
}

func (g *gzipInputStream) Close() error {
	g.zr.Close()
	return g.in.Close()
}

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
		return string(b)
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
