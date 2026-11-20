// Clean-room specification based on java.io.InputStream, OutputStream, FileInputStream, FileOutputStream, BufferedInputStream, BufferedOutputStream Javadoc.

package javaio

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

const (
	defaultStreamBufSize = 64 * 1024
	rawWriteScratchSize  = 256 * 1024
)

var (
	bufferedReaderPool = sync.Pool{
		New: func() any {
			return bufio.NewReaderSize(nil, defaultStreamBufSize)
		},
	}
	bufferedWriterPool = sync.Pool{
		New: func() any {
			return bufio.NewWriterSize(nil, defaultStreamBufSize)
		},
	}
	rawWriteScratchPool = sync.Pool{
		New: func() any {
			return new([rawWriteScratchSize]byte)
		},
	}
)

func ignoreAlreadyClosed(err error) error {
	if err == nil || errors.Is(err, os.ErrClosed) {
		return nil
	}
	return err
}

// InputStream represents an input stream of bytes.
type InputStream interface {
	Read(b []byte) (int, error)
	ReadByte() (byte, error)
	Skip(n int64) (int64, error)
	Available() int
	Close() error
}

// ReaderInputStream adapts an arbitrary io.Reader into an InputStream.
type ReaderInputStream struct {
	r      io.Reader
	closed bool
}

// NewReaderInputStream wraps r as an InputStream.
func NewReaderInputStream(r io.Reader) *ReaderInputStream {
	if r == nil {
		return nil
	}
	return &ReaderInputStream{r: r}
}

func (s *ReaderInputStream) Read(b []byte) (int, error) {
	return s.r.Read(b)
}

func (s *ReaderInputStream) ReadByte() (byte, error) {
	if br, ok := s.r.(io.ByteReader); ok {
		return br.ReadByte()
	}
	var one [1]byte
	n, err := s.r.Read(one[:])
	if n == 0 {
		if err == nil {
			err = io.EOF
		}
		return 0, err
	}
	return one[0], nil
}

func (s *ReaderInputStream) Skip(n int64) (int64, error) {
	if n <= 0 {
		return 0, nil
	}
	return io.CopyN(io.Discard, s.r, n)
}

func (s *ReaderInputStream) Available() int {
	return 0
}

func (s *ReaderInputStream) Close() error {
	if s == nil || s.closed {
		return nil
	}
	s.closed = true
	if c, ok := s.r.(io.Closer); ok {
		return ignoreAlreadyClosed(c.Close())
	}
	return nil
}

// RemainingBytesHinter is an optional interface implemented by InputStreams
// that can report their exact remaining unread byte count without consuming it.
type RemainingBytesHinter interface {
	RemainingBytesHint() (int64, bool)
}

func readerRemainingBytesHint(r io.Reader) (int64, bool) {
	if r == nil {
		return 0, false
	}
	if h, ok := r.(RemainingBytesHinter); ok {
		return h.RemainingBytesHint()
	}
	if l, ok := r.(interface{ Len() int }); ok {
		n := l.Len()
		if n >= 0 {
			return int64(n), true
		}
	}
	return 0, false
}

func (s *ReaderInputStream) RemainingBytesHint() (int64, bool) {
	if s == nil {
		return 0, false
	}
	return readerRemainingBytesHint(s.r)
}

// ReadAllBytes reads all remaining bytes from in (Java 9+ InputStream.readAllBytes).
func ReadAllBytes(in io.Reader) ([]byte, error) {
	if rem, ok := readerRemainingBytesHint(in); ok && rem >= 0 && rem < 1<<31-1 {
		n := int(rem)
		buf := make([]byte, 0, n+1)
		for {
			if len(buf) == cap(buf) {
				buf = append(buf, 0)[:len(buf)]
			}
			m, err := in.Read(buf[len(buf):cap(buf)])
			buf = buf[:len(buf)+m]
			if err != nil {
				if err == io.EOF {
					return buf, nil
				}
				return buf, err
			}
		}
	}
	return io.ReadAll(in)
}

// ByteArrayInputStream contains an internal buffer of bytes to be read from the stream.
type ByteArrayInputStream struct {
	buf []byte
	pos int
}

// NewByteArrayInputStream creates a ByteArrayInputStream using buf.
func NewByteArrayInputStream(buf []byte) *ByteArrayInputStream {
	return &ByteArrayInputStream{buf: buf, pos: 0}
}

func (s *ByteArrayInputStream) Read(b []byte) (int, error) {
	if s.pos >= len(s.buf) {
		return 0, io.EOF
	}
	n := copy(b, s.buf[s.pos:])
	s.pos += n
	return n, nil
}

func (s *ByteArrayInputStream) ReadByte() (byte, error) {
	if s.pos >= len(s.buf) {
		return 0, io.EOF
	}
	b := s.buf[s.pos]
	s.pos++
	return b, nil
}

func (s *ByteArrayInputStream) Skip(n int64) (int64, error) {
	if n <= 0 {
		return 0, nil
	}
	avail := int64(len(s.buf) - s.pos)
	skipped := min(n, avail)
	s.pos += int(skipped)
	return skipped, nil
}

func (s *ByteArrayInputStream) Available() int {
	return len(s.buf) - s.pos
}

func (s *ByteArrayInputStream) RemainingBytesHint() (int64, bool) {
	if s == nil {
		return 0, false
	}
	rem := len(s.buf) - s.pos
	if rem < 0 {
		return 0, true
	}
	return int64(rem), true
}

func (s *ByteArrayInputStream) Close() error {
	return nil
}

// FileInputStream reads raw bytes from a filesystem file.
type FileInputStream struct {
	f      *os.File
	fd     atomic.Int32
	closed atomic.Bool
}

// NewFileInputStream opens path for reading.
func NewFileInputStream(path string) (*FileInputStream, error) {
	if hasRawFileFD {
		fd, err := openRawReadFD(path)
		if err != nil {
			if isRawNotExist(err) {
				return nil, NewFileNotFoundException(err.Error())
			}
			return nil, NewIOException(err.Error())
		}
		s := &FileInputStream{}
		s.fd.Store(int32(fd))
		return s, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, NewFileNotFoundException(err.Error())
		}
		return nil, NewIOException(err.Error())
	}
	s := &FileInputStream{f: f}
	s.fd.Store(-1)
	return s, nil
}

// NewFileInputStreamFromFile opens f.GetPath() for reading.
func NewFileInputStreamFromFile(f *File) (*FileInputStream, error) {
	return NewFileInputStream(f.GetPath())
}

func (s *FileInputStream) Read(b []byte) (int, error) {
	if s == nil {
		return 0, os.ErrClosed
	}
	if s.f != nil {
		return s.f.Read(b)
	}
	fd := int(s.fd.Load())
	if fd < 0 {
		return 0, os.ErrClosed
	}
	return readRawFD(fd, b)
}

func (s *FileInputStream) ReadByte() (byte, error) {
	var one [1]byte
	n, err := s.Read(one[:])
	if n == 0 {
		if err == nil {
			err = io.EOF
		}
		return 0, err
	}
	return one[0], nil
}

func (s *FileInputStream) Skip(n int64) (int64, error) {
	if n <= 0 {
		return 0, nil
	}
	return io.CopyN(io.Discard, s, n)
}

func (s *FileInputStream) Available() int {
	rem, ok := s.RemainingBytesHint()
	if !ok || rem <= 0 {
		return 0
	}
	return int(rem)
}

func (s *FileInputStream) RemainingBytesHint() (int64, bool) {
	if s == nil {
		return 0, false
	}
	if s.f != nil {
		fi, err := s.f.Stat()
		if err != nil || !fi.Mode().IsRegular() {
			return 0, false
		}
		pos, err := s.f.Seek(0, io.SeekCurrent)
		if err != nil || fi.Size() < pos {
			return 0, false
		}
		return fi.Size() - pos, true
	}
	fd := int(s.fd.Load())
	if fd < 0 {
		return 0, false
	}
	return remainingRegularRawFD(fd)
}

func (s *FileInputStream) Close() error {
	if s == nil {
		return nil
	}
	if s.f != nil {
		if !s.closed.CompareAndSwap(false, true) {
			return nil
		}
		return ignoreAlreadyClosed(s.f.Close())
	}
	fd := int(s.fd.Swap(-1))
	if fd < 0 {
		return nil
	}
	return ignoreAlreadyClosed(closeRawFD(fd))
}

// BufferedInputStream adds 64 KiB pooled buffering to an InputStream / io.Reader.
type BufferedInputStream struct {
	br     *bufio.Reader
	src    io.Reader
	closer io.Closer
	closed atomic.Bool
}

// NewBufferedInputStream creates a BufferedInputStream wrapping r.
func NewBufferedInputStream(r io.Reader) *BufferedInputStream {
	var c io.Closer
	if cl, ok := r.(io.Closer); ok {
		c = cl
	}
	br := bufferedReaderPool.Get().(*bufio.Reader)
	br.Reset(r)
	return &BufferedInputStream{br: br, src: r, closer: c}
}

func (b *BufferedInputStream) Read(p []byte) (int, error) {
	if b == nil || b.closed.Load() || b.br == nil {
		return 0, os.ErrClosed
	}
	return b.br.Read(p)
}

func (b *BufferedInputStream) ReadByte() (byte, error) {
	if b == nil || b.closed.Load() || b.br == nil {
		return 0, os.ErrClosed
	}
	return b.br.ReadByte()
}

func (b *BufferedInputStream) Skip(n int64) (int64, error) {
	if b == nil || b.closed.Load() || b.br == nil {
		return 0, os.ErrClosed
	}
	if n <= 0 {
		return 0, nil
	}
	skipped, err := b.br.Discard(int(n))
	return int64(skipped), err
}

func (b *BufferedInputStream) Available() int {
	if b == nil || b.closed.Load() || b.br == nil {
		return 0
	}
	return b.br.Buffered()
}

func (b *BufferedInputStream) RemainingBytesHint() (int64, bool) {
	if b == nil || b.closed.Load() || b.br == nil {
		return 0, false
	}
	rem, ok := readerRemainingBytesHint(b.src)
	if !ok {
		return 0, false
	}
	return int64(b.br.Buffered()) + rem, true
}

func (b *BufferedInputStream) Close() error {
	if b == nil || !b.closed.CompareAndSwap(false, true) {
		return nil
	}
	if b.br != nil {
		b.br.Reset(nil)
		bufferedReaderPool.Put(b.br)
		b.br = nil
	}
	if b.closer != nil {
		return ignoreAlreadyClosed(b.closer.Close())
	}
	return nil
}

// OutputStream represents an output stream of bytes.
type OutputStream interface {
	Write(b []byte) (int, error)
	WriteByte(b byte) error
	Flush() error
	Close() error
}

// ByteArrayOutputStream implements an output stream in which the data is written into a byte array.
type ByteArrayOutputStream struct {
	buf bytes.Buffer
}

// NewByteArrayOutputStream creates a new byte array output stream.
func NewByteArrayOutputStream() *ByteArrayOutputStream {
	return &ByteArrayOutputStream{}
}

func (b *ByteArrayOutputStream) Write(p []byte) (int, error) {
	return b.buf.Write(p)
}

func (b *ByteArrayOutputStream) WriteByte(c byte) error {
	return b.buf.WriteByte(c)
}

func (b *ByteArrayOutputStream) Flush() error {
	return nil
}

func (b *ByteArrayOutputStream) Close() error {
	return nil
}

// ToByteArray returns a newly allocated byte array with the current contents.
func (b *ByteArrayOutputStream) ToByteArray() []byte {
	out := make([]byte, b.buf.Len())
	copy(out, b.buf.Bytes())
	return out
}

// ToString converts the buffer's contents into a string.
func (b *ByteArrayOutputStream) ToString() string {
	return b.buf.String()
}

// Size returns the current size of the buffer.
func (b *ByteArrayOutputStream) Size() int {
	return b.buf.Len()
}

// Reset resets the count field of this byte array output stream to zero.
func (b *ByteArrayOutputStream) Reset() {
	b.buf.Reset()
}

// FileOutputStream writes raw bytes to a filesystem file.
type FileOutputStream struct {
	f      *os.File
	fd     atomic.Int32
	closed atomic.Bool
}

// NewFileOutputStream opens or creates path for writing (truncating existing content).
func NewFileOutputStream(path string) (*FileOutputStream, error) {
	return NewFileOutputStreamAppend(path, false)
}

// NewFileOutputStreamAppend opens or creates path for writing, either appending or truncating.
func NewFileOutputStreamAppend(path string, appendMode bool) (*FileOutputStream, error) {
	flags := os.O_WRONLY | os.O_CREATE
	if appendMode {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}
	if hasRawFileFD {
		fd, err := openRawWriteFD(path, flags)
		if err != nil {
			return nil, NewFileNotFoundException(err.Error())
		}
		s := &FileOutputStream{}
		s.fd.Store(int32(fd))
		return s, nil
	}
	f, err := os.OpenFile(path, flags, 0o666)
	if err != nil {
		return nil, NewFileNotFoundException(err.Error())
	}
	s := &FileOutputStream{f: f}
	s.fd.Store(-1)
	return s, nil
}

// NewFileOutputStreamFlags opens path with explicit os.O_* flags (used by javanio.Files_NewOutputStream).
func NewFileOutputStreamFlags(path string, flags int) (*FileOutputStream, error) {
	if hasRawFileFD {
		fd, err := openRawWriteFD(path, flags)
		if err != nil {
			return nil, NewIOException(err.Error())
		}
		s := &FileOutputStream{}
		s.fd.Store(int32(fd))
		return s, nil
	}
	f, err := os.OpenFile(path, flags, 0o666)
	if err != nil {
		return nil, NewIOException(err.Error())
	}
	s := &FileOutputStream{f: f}
	s.fd.Store(-1)
	return s, nil
}

// NewFileOutputStreamFromFile opens or creates f.GetPath() for writing.
func NewFileOutputStreamFromFile(f *File) (*FileOutputStream, error) {
	return NewFileOutputStream(f.GetPath())
}

// NewFileOutputStreamFromOSFile wraps an already opened *os.File.
func NewFileOutputStreamFromOSFile(f *os.File) *FileOutputStream {
	s := &FileOutputStream{f: f}
	s.fd.Store(-1)
	return s
}

func (s *FileOutputStream) Write(p []byte) (int, error) {
	if s == nil {
		return 0, os.ErrClosed
	}
	if s.f != nil {
		if s.closed.Load() {
			return 0, os.ErrClosed
		}
		return s.f.Write(p)
	}
	fd := int(s.fd.Load())
	if fd < 0 {
		return 0, os.ErrClosed
	}
	return writeRawFD(fd, p)
}

func (s *FileOutputStream) WriteString(str string) (int, error) {
	if s == nil {
		return 0, os.ErrClosed
	}
	if s.f != nil {
		if s.closed.Load() {
			return 0, os.ErrClosed
		}
		return s.f.WriteString(str)
	}
	fd := int(s.fd.Load())
	if fd < 0 {
		return 0, os.ErrClosed
	}
	if len(str) == 0 {
		return 0, nil
	}
	if len(str) <= 4096 {
		var tmp [4096]byte
		n := copy(tmp[:], str)
		return writeRawFD(fd, tmp[:n])
	}
	buf := rawWriteScratchPool.Get().(*[rawWriteScratchSize]byte)
	total := 0
	rem := str
	for len(rem) > 0 {
		chunk := min(len(rem), len(buf))
		copy(buf[:chunk], rem[:chunk])
		n, err := writeRawFD(fd, buf[:chunk])
		total += n
		if err != nil {
			rawWriteScratchPool.Put(buf)
			return total, err
		}
		rem = rem[chunk:]
	}
	rawWriteScratchPool.Put(buf)
	return total, nil
}

func (s *FileOutputStream) WriteByte(c byte) error {
	var one [1]byte
	one[0] = c
	_, err := s.Write(one[:])
	return err
}

func (s *FileOutputStream) Flush() error { return nil }

func (s *FileOutputStream) Close() error {
	if s == nil {
		return nil
	}
	if s.f != nil {
		if !s.closed.CompareAndSwap(false, true) {
			return nil
		}
		return ignoreAlreadyClosed(s.f.Close())
	}
	fd := int(s.fd.Swap(-1))
	if fd < 0 {
		return nil
	}
	return ignoreAlreadyClosed(closeRawFD(fd))
}

// BufferedOutputStream buffers writes to an underlying io.Writer using a pooled 64 KiB buffer.
type BufferedOutputStream struct {
	bw     *bufio.Writer
	closer io.Closer
	closed atomic.Bool
}

// NewBufferedOutputStream creates a BufferedOutputStream wrapping w.
func NewBufferedOutputStream(w io.Writer) *BufferedOutputStream {
	var c io.Closer
	if cl, ok := w.(io.Closer); ok {
		c = cl
	}
	bw := bufferedWriterPool.Get().(*bufio.Writer)
	bw.Reset(w)
	return &BufferedOutputStream{bw: bw, closer: c}
}

func (b *BufferedOutputStream) Write(p []byte) (int, error) {
	if b == nil || b.closed.Load() || b.bw == nil {
		return 0, os.ErrClosed
	}
	return b.bw.Write(p)
}

func (b *BufferedOutputStream) WriteByte(c byte) error {
	if b == nil || b.closed.Load() || b.bw == nil {
		return os.ErrClosed
	}
	return b.bw.WriteByte(c)
}

func (b *BufferedOutputStream) Flush() error {
	if b == nil || b.closed.Load() || b.bw == nil {
		return nil
	}
	return b.bw.Flush()
}

func (b *BufferedOutputStream) Close() error {
	if b == nil || !b.closed.CompareAndSwap(false, true) {
		return nil
	}
	var flushErr error
	if b.bw != nil {
		flushErr = ignoreAlreadyClosed(b.bw.Flush())
		b.bw.Reset(nil)
		bufferedWriterPool.Put(b.bw)
		b.bw = nil
	}
	if b.closer != nil {
		if err := ignoreAlreadyClosed(b.closer.Close()); err != nil && flushErr == nil {
			return err
		}
	}
	return flushErr
}
