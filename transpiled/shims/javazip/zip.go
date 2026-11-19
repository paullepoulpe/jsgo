package javazip

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javanio"
)

const (
	Deflater_NO_COMPRESSION      = 0
	Deflater_BEST_SPEED          = 1
	Deflater_BEST_COMPRESSION    = 9
	Deflater_DEFAULT_COMPRESSION = -1
)

// ZipException models java.util.zip.ZipException.
type ZipException struct {
	javaio.IOException
}

func (*ZipException) zipException() {}

type IZipException interface {
	javaio.IIOException
	zipException()
}

func NewZipException(msg string) *ZipException {
	e := &ZipException{}
	javalang.InitThrowable(e, "java.util.zip.ZipException", msg, nil)
	return e
}

// ZipEntry models java.util.zip.ZipEntry.
type ZipEntry struct {
	name string
	size int64
	dir  bool
	file *zip.File
}

func NewZipEntry(name string) *ZipEntry {
	if name == "" {
		// Java allows non-null names, rejects null.
	}
	return &ZipEntry{
		name: name,
		size: -1,
		dir:  strings.HasSuffix(name, "/"),
	}
}

func zipEntryFromArchive(f *zip.File) *ZipEntry {
	return &ZipEntry{
		name: f.Name,
		size: int64(f.UncompressedSize64),
		dir:  f.FileInfo().IsDir() || strings.HasSuffix(f.Name, "/"),
		file: f,
	}
}

func (e *ZipEntry) GetName() string   { return e.name }
func (e *ZipEntry) String() string    { return e.name }
func (e *ZipEntry) IsDirectory() bool { return e.dir }
func (e *ZipEntry) GetSize() int64    { return e.size }

// ZipInputStream models java.util.zip.ZipInputStream and implements javaio.InputStream.
type ZipInputStream struct {
	in       io.Reader
	cs       *javanio.Charset
	zr       *zip.Reader
	loadErr  error
	loaded   bool
	idx      int
	curRead  io.ReadCloser
	closed   bool
}

func NewZipInputStream(in io.Reader) *ZipInputStream {
	return NewZipInputStreamCharset(in, javanio.StandardCharsets_UTF_8)
}

func NewZipInputStreamCharset(in io.Reader, cs *javanio.Charset) *ZipInputStream {
	if in == nil {
		panic(javalang.NewNullPointerExceptionMsg("in is null"))
	}
	if cs == nil {
		panic(javalang.NewNullPointerExceptionMsg("charset is null"))
	}
	return &ZipInputStream{in: in, cs: cs, idx: -1}
}

func (z *ZipInputStream) ensureLoaded() error {
	if z.loaded {
		return z.loadErr
	}
	z.loaded = true
	data, err := io.ReadAll(z.in)
	if err != nil {
		z.loadErr = javaio.NewIOException(err.Error())
		return z.loadErr
	}
	if len(data) == 0 {
		// Java ZipInputStream returns null on GetNextEntry() for an empty stream!
		return nil
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		z.loadErr = NewZipException(err.Error())
		return z.loadErr
	}
	z.zr = zr
	return nil
}

func (z *ZipInputStream) GetNextEntry() (*ZipEntry, error) {
	if z.closed {
		return nil, javaio.NewIOException("Stream closed")
	}
	if z.curRead != nil {
		_ = z.curRead.Close()
		z.curRead = nil
	}
	if err := z.ensureLoaded(); err != nil {
		return nil, err
	}
	if z.zr == nil {
		return nil, nil
	}
	z.idx++
	if z.idx >= len(z.zr.File) {
		return nil, nil
	}
	zf := z.zr.File[z.idx]
	rc, err := zf.Open()
	if err != nil {
		return nil, NewZipException(err.Error())
	}
	z.curRead = rc
	entry := zipEntryFromArchive(zf)
	if zf.NonUTF8 && z.cs != nil && !z.cs.Equals(javanio.StandardCharsets_UTF_8) {
		entry.name = z.cs.Decode([]byte(zf.Name))
	}
	return entry, nil
}

func (z *ZipInputStream) CloseEntry() error {
	if z.closed {
		return javaio.NewIOException("Stream closed")
	}
	if z.curRead != nil {
		err := z.curRead.Close()
		z.curRead = nil
		if err != nil {
			return javaio.NewIOException(err.Error())
		}
	}
	return nil
}

func (z *ZipInputStream) Read(b []byte) (int, error) {
	if z.closed {
		return 0, javaio.NewIOException("Stream closed")
	}
	if z.curRead == nil {
		return 0, io.EOF
	}
	return z.curRead.Read(b)
}

func (z *ZipInputStream) ReadByte() (byte, error) {
	var one [1]byte
	n, err := z.Read(one[:])
	if n == 0 {
		if err == nil {
			err = io.EOF
		}
		return 0, err
	}
	return one[0], nil
}

func (z *ZipInputStream) Skip(n int64) (int64, error) {
	if z.closed {
		return 0, javaio.NewIOException("Stream closed")
	}
	if n <= 0 || z.curRead == nil {
		return 0, nil
	}
	return io.CopyN(io.Discard, z.curRead, n)
}

func (z *ZipInputStream) Available() int {
	if z.closed || z.curRead == nil {
		return 0
	}
	return 1
}

func isAlreadyClosed(err error) bool {
	return err == nil || errors.Is(err, os.ErrClosed)
}

func (z *ZipInputStream) Close() error {
	if z.closed {
		return nil
	}
	z.closed = true
	if z.curRead != nil {
		_ = z.curRead.Close()
		z.curRead = nil
	}
	if c, ok := z.in.(io.Closer); ok {
		if err := c.Close(); !isAlreadyClosed(err) {
			return javaio.NewIOException(err.Error())
		}
	}
	return nil
}

// ZipFile models java.util.zip.ZipFile.
type ZipFile struct {
	name    string
	rc      *zip.ReadCloser
	byName  map[string]*ZipEntry
	entries []*ZipEntry
	closed  bool
}

func NewZipFile(name string) (*ZipFile, error) {
	rc, err := zip.OpenReader(name)
	if err != nil {
		return nil, NewZipException(err.Error())
	}
	byName := make(map[string]*ZipEntry, len(rc.File))
	entries := make([]*ZipEntry, len(rc.File))
	for i, f := range rc.File {
		e := zipEntryFromArchive(f)
		entries[i] = e
		if _, exists := byName[e.name]; !exists {
			byName[e.name] = e
		}
	}
	return &ZipFile{
		name:    name,
		rc:      rc,
		byName:  byName,
		entries: entries,
	}, nil
}

func NewZipFileFromFile(f *javaio.File) (*ZipFile, error) {
	if f == nil {
		panic(javalang.NewNullPointerExceptionMsg("file is null"))
	}
	return NewZipFile(f.GetPath())
}

func (zf *ZipFile) GetName() string { return zf.name }
func (zf *ZipFile) Size() int       { return len(zf.entries) }

func (zf *ZipFile) GetEntry(name string) *ZipEntry {
	if zf.closed {
		panic(javalang.NewIllegalStateException("zip file closed"))
	}
	return zf.byName[name]
}

func (zf *ZipFile) Entries() []*ZipEntry {
	if zf.closed {
		panic(javalang.NewIllegalStateException("zip file closed"))
	}
	out := make([]*ZipEntry, len(zf.entries))
	copy(out, zf.entries)
	return out
}

func (zf *ZipFile) GetInputStream(entry *ZipEntry) (javaio.InputStream, error) {
	if entry == nil {
		panic(javalang.NewNullPointerExceptionMsg("entry is null"))
	}
	if zf.closed {
		panic(javalang.NewIllegalStateException("zip file closed"))
	}
	target := entry.file
	if target == nil {
		if found := zf.byName[entry.name]; found != nil {
			target = found.file
		}
	}
	if target == nil {
		return nil, nil
	}
	r, err := target.Open()
	if err != nil {
		return nil, NewZipException(err.Error())
	}
	return javaio.NewReaderInputStream(r), nil
}

func (zf *ZipFile) Close() error {
	if zf.closed {
		return nil
	}
	zf.closed = true
	if err := zf.rc.Close(); !isAlreadyClosed(err) {
		return javaio.NewIOException(err.Error())
	}
	return nil
}

// GZIPInputStream models java.util.zip.GZIPInputStream and implements javaio.InputStream.
type GZIPInputStream struct {
	gr     *gzip.Reader
	closer io.Closer
	closed bool
}

func NewGZIPInputStream(in io.Reader) (*GZIPInputStream, error) {
	if in == nil {
		panic(javalang.NewNullPointerExceptionMsg("in is null"))
	}
	gr, err := gzip.NewReader(in)
	if err != nil {
		return nil, javaio.NewIOException(err.Error())
	}
	c, _ := in.(io.Closer)
	return &GZIPInputStream{gr: gr, closer: c}, nil
}

func NewGZIPInputStreamWithSize(in io.Reader, size int) (*GZIPInputStream, error) {
	if size <= 0 {
		panic(javalang.NewIllegalArgumentException("buffer size <= 0"))
	}
	return NewGZIPInputStream(in)
}

func (g *GZIPInputStream) Read(b []byte) (int, error) {
	if g.closed {
		return 0, javaio.NewIOException("Stream closed")
	}
	return g.gr.Read(b)
}

func (g *GZIPInputStream) ReadByte() (byte, error) {
	var one [1]byte
	n, err := g.Read(one[:])
	if n == 0 {
		if err == nil {
			err = io.EOF
		}
		return 0, err
	}
	return one[0], nil
}

func (g *GZIPInputStream) Skip(n int64) (int64, error) {
	if g.closed {
		return 0, javaio.NewIOException("Stream closed")
	}
	if n <= 0 {
		return 0, nil
	}
	return io.CopyN(io.Discard, g.gr, n)
}

func (g *GZIPInputStream) Available() int {
	if g.closed {
		return 0
	}
	return 1
}

func (g *GZIPInputStream) Close() error {
	if g.closed {
		return nil
	}
	g.closed = true
	err1 := g.gr.Close()
	var err2 error
	if g.closer != nil {
		err2 = g.closer.Close()
	}
	if !isAlreadyClosed(err1) {
		return javaio.NewIOException(err1.Error())
	}
	if !isAlreadyClosed(err2) {
		return javaio.NewIOException(err2.Error())
	}
	return nil
}

// GZIPOutputStream models java.util.zip.GZIPOutputStream and implements javaio.OutputStream.
type GZIPOutputStream struct {
	gw       *gzip.Writer
	closer   io.Closer
	finished bool
	closed   bool
}

func NewGZIPOutputStream(out io.Writer) (*GZIPOutputStream, error) {
	return NewGZIPOutputStreamLevel(out, Deflater_DEFAULT_COMPRESSION)
}

func NewGZIPOutputStreamLevel(out io.Writer, level int) (*GZIPOutputStream, error) {
	if out == nil {
		panic(javalang.NewNullPointerExceptionMsg("out is null"))
	}
	gw, err := gzip.NewWriterLevel(out, level)
	if err != nil {
		return nil, javaio.NewIOException(err.Error())
	}
	c, _ := out.(io.Closer)
	return &GZIPOutputStream{gw: gw, closer: c}, nil
}

// NewFastGzipOutputStream creates a BEST_SPEED GZIPOutputStream (used by FastGzipOutputStream).
func NewFastGzipOutputStream(out io.Writer) (*GZIPOutputStream, error) {
	return NewGZIPOutputStreamLevel(out, Deflater_BEST_SPEED)
}

func (g *GZIPOutputStream) Write(b []byte) (int, error) {
	if g.closed {
		return 0, javaio.NewIOException("Stream closed")
	}
	return g.gw.Write(b)
}

func (g *GZIPOutputStream) WriteByte(c byte) error {
	_, err := g.Write([]byte{c})
	return err
}

func (g *GZIPOutputStream) Flush() error {
	if g.closed {
		return javaio.NewIOException("Stream closed")
	}
	if err := g.gw.Flush(); err != nil {
		return javaio.NewIOException(err.Error())
	}
	return nil
}

func (g *GZIPOutputStream) Finish() error {
	if g.finished {
		return nil
	}
	g.finished = true
	if err := g.gw.Close(); err != nil {
		return javaio.NewIOException(err.Error())
	}
	return nil
}

func (g *GZIPOutputStream) Close() error {
	if g.closed {
		return nil
	}
	g.closed = true
	err1 := g.Finish()
	var err2 error
	if g.closer != nil {
		err2 = g.closer.Close()
	}
	if !isAlreadyClosed(err1) {
		return err1
	}
	if !isAlreadyClosed(err2) {
		return javaio.NewIOException(err2.Error())
	}
	return nil
}
