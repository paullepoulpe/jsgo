package jutil

import (
	"bufio"
	"compress/gzip"
	"hash/crc32"
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

// AppendCharSequence is StringBuilder.append(CharSequence) as an
// Appendable.
func (b *StringBuilder) AppendCharSequence(csq string) { b.Append(csq) }

// AppendCharacter is StringBuilder.append(char) as an Appendable.
func (b *StringBuilder) AppendCharacter(c uint16) { b.AppendChar(c) }

// Go-only: language-adapter — java.io.Writer. An IOException is thrown as
// a panic (jutil.IOException), like every other Java exception.

// Writer is java.io.Writer (the methods ported code uses).
type Writer interface {
	// Write is Writer.write(String).
	Write(s string)
	// Flush is Writer.flush().
	Flush()
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
	var sb StringBuilder
	buf := make([]uint16, 8192)
	for {
		n := r.Read(buf, 0, len(buf))
		if n < 0 {
			break
		}
		for _, u := range buf[:n] {
			sb.AppendChar(u)
		}
	}
	return sb.String()
}

// CharStreams_ReadLines is Guava CharStreams.readLines(Readable): the lines
// of the input, split at "\n", "\r" and "\r\n" (Guava's LineReader). The
// line terminators are not included, and an input that ends with a
// terminator has no empty last line.
func CharStreams_ReadLines(r JavaReader) []string {
	units := StringToUTF16(CharStreams_ToString(r))
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
	f, err := os.Open(name)
	if err == nil {
		if st, serr := f.Stat(); serr == nil && st.IsDir() {
			f.Close()
			panic(NewFileNotFoundException(name + " (Is a directory)"))
		}
		return f
	}
	panic(FileNotFoundExceptionFromError(err, name))
}

// Go-only: language-adapter — java.io.FileOutputStream as a Go
// io.WriteCloser over an *os.File.

// NewFileOutputStream is new FileOutputStream(String name): it creates or
// truncates the file, and throws FileNotFoundException("<name> (<reason>)")
// when the file cannot be opened for writing. A write error is thrown as
// IOException by the writers that wrap it.
func NewFileOutputStream(name string) io.WriteCloser {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		panic(FileNotFoundExceptionFromError(err, name))
	}
	return f
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
	b, err := os.ReadFile(pathStr)
	if err != nil {
		panic(IOExceptionFromError(err, pathStr))
	}
	if (cs == nil || cs == StandardCharsets_UTF_8) && utf8.Valid(b) {
		return string(b)
	}
	return cs.DecodeStrict(b)
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
	zr, err := gzip.NewReader(bufio.NewReaderSize(in, size))
	if err != nil {
		if err == gzip.ErrHeader {
			panic(NewZipException("Not in GZIP format"))
		}
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			// Java: new EOFException(), with a null message.
			e := &EOFException{}
			InitThrowable(e, "java.io.EOFException", nil, nil)
			panic(e)
		}
		panic(NewIOException(err.Error()))
	}
	return &gzipInputStream{in: in, zr: zr}
}

func (g *gzipInputStream) Read(b []byte) (int, error) {
	return g.zr.Read(b)
}

func (g *gzipInputStream) Close() error {
	g.zr.Close()
	return g.in.Close()
}

// Go-only: language-adapter — java.util.zip.GZIPOutputStream, written from
// its Javadoc ("This class implements a stream filter for writing
// compressed data in the GZIP file format") and RFC 1952, over a Go
// io.Writer. The compressed data is a raw deflate stream from a Deflater
// with the default compression level and nowrap (no zlib header and
// checksum); the protected field def is DeflaterOutputStream.Def. Like
// Deflater's, the compressed bytes can differ from Java's; the member
// decompresses to exactly the data written.
//
// Observed behaviour (black-box, JDK 21, JutilDeflaterVectors.java):
//   - the constructor writes the 10-byte member header at once: ID1 0x1f,
//     ID2 0x8b, CM 8 (deflate), FLG 0, MTIME 0, XFL 0 and OS 0xff
//     (unknown), whatever the compression level set later;
//   - finish writes the rest of the deflate stream, then the trailer: the
//     CRC-32 of the uncompressed data and its length modulo 2^32, both
//     little-endian. A second finish writes nothing.

// gzipOutputStreamHeader is the member header that GZIPOutputStream
// writes.
var gzipOutputStreamHeader = []byte{0x1f, 0x8b, 8, 0, 0, 0, 0, 0, 0, 0xff}

// GZIPOutputStream is java.util.zip.GZIPOutputStream.
type GZIPOutputStream struct {
	DeflaterOutputStream
	// crc is the protected field crc: the CRC-32 of the uncompressed data.
	crc     uint32
	written uint32 // uncompressed length modulo 2^32
}

// NewGZIPOutputStream is new GZIPOutputStream(OutputStream out), with the
// default buffer size of 512 bytes.
func NewGZIPOutputStream(out io.Writer) *GZIPOutputStream {
	return NewGZIPOutputStreamWithSizeAndSyncFlush(out, 512, false)
}

// NewGZIPOutputStreamWithSize is new GZIPOutputStream(OutputStream out,
// int size).
func NewGZIPOutputStreamWithSize(out io.Writer, size int) *GZIPOutputStream {
	return NewGZIPOutputStreamWithSizeAndSyncFlush(out, size, false)
}

// NewGZIPOutputStreamWithSizeAndSyncFlush is new GZIPOutputStream(
// OutputStream out, int size, boolean syncFlush).
func NewGZIPOutputStreamWithSizeAndSyncFlush(out io.Writer, size int, syncFlush bool) *GZIPOutputStream {
	g := &GZIPOutputStream{}
	g.InitGZIPOutputStream(out, size, syncFlush)
	return g
}

// InitGZIPOutputStream is the constructor body, for subclasses.
func (g *GZIPOutputStream) InitGZIPOutputStream(out io.Writer, size int, syncFlush bool) {
	g.InitDeflaterOutputStream(
		out, NewDeflaterWithLevelAndNowrap(Deflater_DEFAULT_COMPRESSION, true), size, syncFlush)
	g.finishSelf = g.Finish
	g.writeOut(gzipOutputStreamHeader)
	g.crc = 0
	g.written = 0
}

// WriteByte is write(int b).
func (g *GZIPOutputStream) WriteByte(b byte) error {
	_, err := g.Write([]byte{b})
	return err
}

// Write is write(byte[] buf, int off, int len) for the slice b: writes an
// array of bytes to the compressed output stream.
func (g *GZIPOutputStream) Write(b []byte) (int, error) {
	g.DeflaterOutputStream.writeBytes(b)
	g.crc = crc32.Update(g.crc, crc32.IEEETable, b)
	g.written += uint32(len(b))
	return len(b), nil
}

// Finish is finish(): finishes writing compressed data to the output
// stream without closing the underlying stream.
func (g *GZIPOutputStream) Finish() {
	if g.Def.Finished() {
		return
	}
	g.DeflaterOutputStream.finishDeflater()
	trailer := []byte{
		byte(g.crc), byte(g.crc >> 8), byte(g.crc >> 16), byte(g.crc >> 24),
		byte(g.written), byte(g.written >> 8), byte(g.written >> 16), byte(g.written >> 24),
	}
	g.writeOut(trailer)
}
