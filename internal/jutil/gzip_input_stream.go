package jutil

import (
	"bufio"
	"compress/gzip"
	"io"
)

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
