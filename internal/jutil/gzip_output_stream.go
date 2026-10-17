package jutil

import (
	"hash/crc32"
	"io"
)

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
