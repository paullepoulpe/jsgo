package jutil

import "io"

// Go-only: language-adapter — java.util.zip.DeflaterOutputStream, written
// from its Javadoc, over a Go io.Writer: "This class implements an output
// stream filter for compressing data in the 'deflate' compression format."
// Write errors of the underlying stream are thrown as IOException. Write
// and Close also satisfy io.WriteCloser; they throw instead of returning
// an error.
//
// Observed behaviour (black-box, JDK 21, JutilDeflaterVectors.java):
// writing after finish throws IOException("write beyond end of stream");
// close finishes the stream, then closes the underlying stream; with
// syncFlush, flush makes everything written so far decodable from the
// output.

// DeflaterOutputStream is java.util.zip.DeflaterOutputStream.
type DeflaterOutputStream struct {
	out io.Writer
	// Def is the protected field def: the compressor for this stream.
	Def *Deflater
	// Buf is the protected field buf: the output buffer for writing
	// compressed data.
	Buf []byte

	syncFlush bool
	closed    bool
	// finishSelf is the virtual finish() (a subclass's override).
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
