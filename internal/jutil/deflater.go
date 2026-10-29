package jutil

import (
	"bytes"
	"compress/flate"
	"hash"
	"hash/adler32"
	"io"
)

// Go-only: language-adapter — java.util.zip.Deflater, written from its
// Javadoc, compressing with Go's compress/flate. The JDK compresses with
// zlib, so the compressed bytes can differ from Java's; the stream is a
// valid DEFLATE stream (RFC 1951), wrapped as RFC 1950 describes unless
// nowrap, and decompresses to exactly the input, which is all a reader of
// the output can rely on. Oracles and tools compare such output
// decompressed (tools/oracle/README.md).
//
// The stream contract follows the Javadoc and JDK 21's observed behaviour
// (tools/oracle/java/JutilDeflaterVectors.java):
//   - deflate consumes input only while its output buffer has room, so
//     needsInput() stays false until all the input set by setInput has
//     been taken; it returns 0 once it needs input;
//   - setLevel and setStrategy before the first deflate apply to all the
//     data; later, "the next invocation of deflate will compress the input
//     available so far with the old level (and may be flushed)": that
//     input is compressed and sync-flushed, and the following data goes to
//     a new compressor (a deflate stream may change its blocks' encoding
//     at any block boundary);
//   - SYNC_FLUSH makes all input so far decodable from the output;
//     FULL_FLUSH also stops later data from referring to earlier data;
//   - the exceptions: new Deflater with a bad level and setStrategy with a
//     bad strategy throw IllegalArgumentException without a message,
//     setLevel throws IllegalArgumentException("invalid compression
//     level"), deflate with a bad flush mode IllegalArgumentException
//     without a message; after end, deflate, reset and the getters throw
//     NullPointerException("Deflater has been closed"), while setInput,
//     setLevel, setStrategy, needsInput, finish and finished still work.
//
// FILTERED is accepted; Go's compressor has no filtered mode, so it only
// changes the compressed bytes in Java. HUFFMAN_ONLY is flate.HuffmanOnly.
// Preset dictionaries and ByteBuffer input and output are not emulated.

// Deflater constants (the Javadoc's "Constant Field Values").
const (
	Deflater_DEFLATED            = 8
	Deflater_NO_COMPRESSION      = 0
	Deflater_BEST_SPEED          = 1
	Deflater_BEST_COMPRESSION    = 9
	Deflater_DEFAULT_COMPRESSION = -1
	Deflater_FILTERED            = 1
	Deflater_HUFFMAN_ONLY        = 2
	Deflater_DEFAULT_STRATEGY    = 0
	Deflater_NO_FLUSH            = 0
	Deflater_SYNC_FLUSH          = 2
	Deflater_FULL_FLUSH          = 3
)

// deflaterChunk bounds the input handed to the compressor at a time, so
// that output is produced as the output buffers allow.
const deflaterChunk = 32 * 1024

// Deflater is java.util.zip.Deflater.
type Deflater struct {
	nowrap       bool
	input        []byte
	level        int
	strategy     int
	setParams    bool
	finish       bool
	finished     bool
	closed       bool
	bytesRead    int64
	bytesWritten int64

	w         *flate.Writer // the current compressor, nil before data or after a params change
	pending   bytes.Buffer  // compressed bytes not yet returned by deflate
	started   bool          // the zlib header is written
	streamEnd bool          // the final block (and trailer) are written
	unflushed bool          // input went to w since its last flush
	adler     hash.Hash32   // of the input taken so far (zlib wrapper)
}

// NewDeflater is new Deflater(): the default compression level, with the
// zlib header and checksum.
func NewDeflater() *Deflater {
	return NewDeflaterWithLevelAndNowrap(Deflater_DEFAULT_COMPRESSION, false)
}

// NewDeflaterWithLevel is new Deflater(int level).
func NewDeflaterWithLevel(level int) *Deflater {
	return NewDeflaterWithLevelAndNowrap(level, false)
}

// NewDeflaterWithLevelAndNowrap is new Deflater(int level, boolean nowrap):
// "If 'nowrap' is true then the ZLIB header and checksum fields will not
// be used in order to support the compression format used in both GZIP and
// PKZIP."
func NewDeflaterWithLevelAndNowrap(level int, nowrap bool) *Deflater {
	if (level < 0 || level > 9) && level != Deflater_DEFAULT_COMPRESSION {
		panic(NewIllegalArgumentException())
	}
	return &Deflater{
		nowrap:   nowrap,
		level:    level,
		strategy: Deflater_DEFAULT_STRATEGY,
		adler:    adler32.New(),
	}
}

func (d *Deflater) ensureOpen() {
	if d.closed {
		panic(NewNullPointerExceptionMsg("Deflater has been closed"))
	}
}

// SetInput is setInput(byte[] input, int off, int len) for the slice
// input: it replaces any input not yet processed. The Deflater keeps a
// reference to the slice until it has processed it.
func (d *Deflater) SetInput(input []byte) {
	d.input = input
}

// SetLevel is setLevel(int level).
func (d *Deflater) SetLevel(level int) {
	if (level < 0 || level > 9) && level != Deflater_DEFAULT_COMPRESSION {
		panic(NewIllegalArgumentExceptionMsg("invalid compression level"))
	}
	if d.level != level {
		d.level = level
		d.setParams = true
	}
}

// SetStrategy is setStrategy(int strategy).
func (d *Deflater) SetStrategy(strategy int) {
	switch strategy {
	case Deflater_DEFAULT_STRATEGY, Deflater_FILTERED, Deflater_HUFFMAN_ONLY:
	default:
		panic(NewIllegalArgumentException())
	}
	if d.strategy != strategy {
		d.strategy = strategy
		d.setParams = true
	}
}

// NeedsInput is needsInput(): true if no data remains in the input buffer.
func (d *Deflater) NeedsInput() bool {
	return len(d.input) == 0
}

// Finish is finish(): compression should end with the current contents of
// the input buffer.
func (d *Deflater) Finish() {
	d.finish = true
}

// Finished is finished(): true if the end of the compressed data output
// stream has been reached.
func (d *Deflater) Finished() bool {
	return d.finished
}

// Deflate is deflate(byte[] b) (and deflate(b, off, len) for the slice b):
// it compresses with NO_FLUSH and returns the number of bytes written to b.
func (d *Deflater) Deflate(b []byte) int {
	return d.DeflateWithFlush(b, Deflater_NO_FLUSH)
}

// flateLevel is the compress/flate level for the current parameters.
func (d *Deflater) flateLevel() int {
	if d.strategy == Deflater_HUFFMAN_ONLY {
		return flate.HuffmanOnly
	}
	return d.level // -1 is flate.DefaultCompression, 0..9 are the same
}

// writer returns the current compressor, starting one if needed.
func (d *Deflater) writer() *flate.Writer {
	if d.w == nil {
		w, err := flate.NewWriter(&d.pending, d.flateLevel())
		if err != nil {
			panic(NewIllegalStateExceptionMsg(err.Error()))
		}
		d.w = w
	}
	return d.w
}

// take compresses n bytes of the input.
func (d *Deflater) take(n int) {
	chunk := d.input[:n]
	d.writer().Write(chunk)
	if !d.nowrap {
		d.adler.Write(chunk)
	}
	d.input = d.input[n:]
	d.bytesRead += int64(n)
	d.unflushed = true
}

// zlibHeader is the RFC 1950 header: CM 8 (deflate) with a 32K window,
// and FLEVEL from the compression level (RFC 1950: 0 fastest, 1 fast,
// 2 default, 3 maximum compression).
func (d *Deflater) zlibHeader() []byte {
	const cmf = 0x78
	var flevel byte
	switch level := d.level; {
	case d.strategy == Deflater_HUFFMAN_ONLY || level == 0 || level == 1:
		flevel = 0
	case level >= 2 && level <= 5:
		flevel = 1
	case level == 6 || level == Deflater_DEFAULT_COMPRESSION:
		flevel = 2
	default:
		flevel = 3
	}
	flg := flevel << 6
	flg += byte(31 - (uint16(cmf)<<8|uint16(flg))%31)
	return []byte{cmf, flg}
}

// DeflateWithFlush is deflate(byte[] b, int off, int len, int flush) for
// the slice b.
func (d *Deflater) DeflateWithFlush(b []byte, flush int) int {
	d.ensureOpen()
	if flush != Deflater_NO_FLUSH && flush != Deflater_SYNC_FLUSH && flush != Deflater_FULL_FLUSH {
		panic(NewIllegalArgumentException())
	}
	if !d.started {
		d.started = true
		if !d.nowrap {
			d.pending.Write(d.zlibHeader())
		}
	}
	if d.setParams {
		d.setParams = false
		if d.w != nil {
			// The input available so far is compressed with the old
			// parameters and flushed; later data goes to a new compressor.
			if len(d.input) > 0 {
				d.take(len(d.input))
			}
			d.w.Flush()
			d.w = nil
			d.unflushed = false
		}
	}
	n := 0
	for {
		n += copy(b[n:], d.pending.Next(len(b)-n))
		if n == len(b) || d.streamEnd {
			break
		}
		if len(d.input) > 0 {
			d.take(min(len(d.input), deflaterChunk))
			continue
		}
		if d.finish {
			d.writer().Close()
			if !d.nowrap {
				sum := d.adler.Sum32()
				d.pending.Write([]byte{byte(sum >> 24), byte(sum >> 16), byte(sum >> 8), byte(sum)})
			}
			d.streamEnd = true
			continue
		}
		if flush != Deflater_NO_FLUSH && d.unflushed {
			d.writer().Flush()
			d.unflushed = false
			if flush == Deflater_FULL_FLUSH {
				d.w = nil
			}
			continue
		}
		break
	}
	if d.streamEnd && d.pending.Len() == 0 {
		d.finished = true
	}
	d.bytesWritten += int64(n)
	return n
}

// GetAdler is getAdler(): the ADLER-32 value of the uncompressed data (1
// for a raw deflate stream).
func (d *Deflater) GetAdler() int32 {
	d.ensureOpen()
	if d.nowrap {
		return 1
	}
	return int32(d.adler.Sum32())
}

// GetBytesRead is getBytesRead(): the total number of uncompressed bytes
// input so far.
func (d *Deflater) GetBytesRead() int64 {
	d.ensureOpen()
	return d.bytesRead
}

// GetBytesWritten is getBytesWritten(): the total number of compressed
// bytes output so far.
func (d *Deflater) GetBytesWritten() int64 {
	d.ensureOpen()
	return d.bytesWritten
}

// GetTotalIn is getTotalIn(): getBytesRead() as an int.
func (d *Deflater) GetTotalIn() int32 { return int32(d.GetBytesRead()) }

// GetTotalOut is getTotalOut(): getBytesWritten() as an int.
func (d *Deflater) GetTotalOut() int32 { return int32(d.GetBytesWritten()) }

// Reset is reset(): resets the deflater so that a new set of input data can
// be processed. Keeps current compression level and strategy settings.
func (d *Deflater) Reset() {
	d.ensureOpen()
	d.w = nil
	d.pending.Reset()
	d.started = false
	d.streamEnd = false
	d.unflushed = false
	d.adler.Reset()
	d.setParams = false
	d.finish = false
	d.finished = false
	d.input = nil
	d.bytesRead = 0
	d.bytesWritten = 0
}

// End is end(): closes the compressor and discards any unprocessed input.
// Once it has been called, the Deflater can no longer be used.
func (d *Deflater) End() {
	d.closed = true
	d.input = nil
	d.w = nil
	d.pending.Reset()
}

// DeflaterOutputStream is java.util.zip.DeflaterOutputStream.
type DeflaterOutputStream struct {
	out io.Writer
	// Def is the protected field def: the compressor for this stream.
	Def *Deflater
	// Buf is the protected field buf: the output buffer for writing
	// compressed data.
	Buf []byte

	syncFlush  bool
	closed     bool
	finishSelf func()
}

// NewDeflaterOutputStream is new DeflaterOutputStream(OutputStream out,
// Deflater def, int size, boolean syncFlush).
func NewDeflaterOutputStream(out io.Writer, def *Deflater, size int, syncFlush bool) *DeflaterOutputStream {
	d := &DeflaterOutputStream{}
	d.InitDeflaterOutputStream(out, def, size, syncFlush)
	return d
}

// InitDeflaterOutputStream is the constructor body, for subclasses.
func (d *DeflaterOutputStream) InitDeflaterOutputStream(out io.Writer, def *Deflater, size int, syncFlush bool) {
	if out == nil || def == nil {
		panic(NewNullPointerException())
	} else if size <= 0 {
		panic(NewIllegalArgumentExceptionMsg("buffer size <= 0"))
	}
	d.out = out
	d.Def = def
	d.Buf = make([]byte, size)
	d.syncFlush = syncFlush
	d.finishSelf = d.finishDeflater
}

func (d *DeflaterOutputStream) writeOut(b []byte) {
	if len(b) == 0 {
		return
	}
	if _, err := d.out.Write(b); err != nil {
		panic(NewIOException(err.Error()))
	}
}

// WriteByte is write(int b): writes a byte to the compressed output stream.
func (d *DeflaterOutputStream) WriteByte(b byte) error {
	d.writeBytes([]byte{b})
	return nil
}

// Write is write(byte[] b, int off, int len) for the slice b: writes an
// array of bytes to the compressed output stream.
func (d *DeflaterOutputStream) Write(b []byte) (int, error) {
	d.writeBytes(b)
	return len(b), nil
}

func (d *DeflaterOutputStream) writeBytes(b []byte) {
	if d.Def.Finished() {
		panic(NewIOException("write beyond end of stream"))
	}
	if len(b) == 0 {
		return
	}
	d.Def.SetInput(b)
	for !d.Def.NeedsInput() {
		d.DeflateNext()
	}
}

// DeflateNext is the protected deflate(): writes the next block of
// compressed data to the output stream.
func (d *DeflaterOutputStream) DeflateNext() {
	n := d.Def.Deflate(d.Buf)
	if n > 0 {
		d.writeOut(d.Buf[:n])
	}
}

// Finish is finish(): finishes writing compressed data to the output
// stream without closing the underlying stream.
func (d *DeflaterOutputStream) Finish() {
	d.finishSelf()
}

func (d *DeflaterOutputStream) finishDeflater() {
	if !d.Def.Finished() {
		d.Def.Finish()
		for !d.Def.Finished() {
			d.DeflateNext()
		}
	}
}

// Close is close(): writes remaining compressed data to the output stream
// and closes the underlying stream (when it is an io.Closer).
func (d *DeflaterOutputStream) Close() error {
	if d.closed {
		return nil
	}
	d.closed = true
	finished := false
	defer func() {
		if c, ok := d.out.(io.Closer); ok {
			if err := c.Close(); err != nil && finished {
				panic(NewIOException(err.Error()))
			}
		}
	}()
	d.Finish()
	finished = true
	return nil
}

// Flush is flush(): with syncFlush, it first flushes the compressor with
// SYNC_FLUSH; then it flushes the underlying stream (when it has a Flush
// method).
func (d *DeflaterOutputStream) Flush() {
	if d.syncFlush && !d.Def.Finished() {
		for {
			n := d.Def.DeflateWithFlush(d.Buf, Deflater_SYNC_FLUSH)
			if n <= 0 {
				break
			}
			d.writeOut(d.Buf[:n])
			if n < len(d.Buf) {
				break
			}
		}
	}
	if f, ok := d.out.(interface{ Flush() error }); ok {
		if err := f.Flush(); err != nil {
			panic(NewIOException(err.Error()))
		}
	}
}
