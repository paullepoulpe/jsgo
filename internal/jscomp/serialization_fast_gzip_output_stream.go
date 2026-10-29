// Port of com/google/javascript/jscomp/serialization/FastGzipOutputStream.java @ closure-compiler v20261021

/*
 * Copyright 2026 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"bufio"
	"bytes"
	"compress/flate"
	"hash/crc32"
	"io"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// FastGzipOutputStream is a GZIPOutputStream that uses BEST_SPEED
// compression level and a reasonable buffer size.
//
// Go-only: language-adapter — uses Go's compress/flate at BestSpeed: the same
// content as Java's, possibly different compressed bytes (PORTING.md §7).
//
// Java: public final class FastGzipOutputStream extends GZIPOutputStream
type FastGzipOutputStream struct {
	out      io.Writer
	bw       *bufio.Writer
	fw       *flate.Writer
	crc      uint32
	written  uint32
	finished bool
	closed   bool
}

// Java: FastGzipOutputStream.DEFAULT_BUFFER_SIZE
const fastGzipOutputStream_DEFAULT_BUFFER_SIZE = 64 * 1024

// Java: FastGzipOutputStream(OutputStream)
func NewFastGzipOutputStream(out io.Writer) *FastGzipOutputStream {
	return NewFastGzipOutputStreamWithSize(out, fastGzipOutputStream_DEFAULT_BUFFER_SIZE)
}

// Java: FastGzipOutputStream(OutputStream, int)
func NewFastGzipOutputStreamWithSize(out io.Writer, size int) *FastGzipOutputStream {
	if out == nil {
		panic(jutil.NewNullPointerException())
	}
	if size <= 0 {
		panic(jutil.NewIllegalArgumentExceptionMsg("buffer size <= 0"))
	}
	if _, err := out.Write(parallelGzipMemberHeader[:]); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	bw := bufio.NewWriterSize(out, size)
	fw, err := flate.NewWriter(bw, flate.BestSpeed)
	if err != nil {
		panic(jutil.NewIllegalStateExceptionMsg(err.Error()))
	}
	return &FastGzipOutputStream{
		out: out,
		bw:  bw,
		fw:  fw,
	}
}

func (f *FastGzipOutputStream) WriteByte(b byte) error {
	_, err := f.Write([]byte{b})
	return err
}

func (f *FastGzipOutputStream) Write(b []byte) (int, error) {
	if f.finished || f.closed {
		panic(jutil.NewIOException("write beyond end of stream"))
	}
	if len(b) == 0 {
		return 0, nil
	}
	if _, err := f.fw.Write(b); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	f.crc = crc32.Update(f.crc, crc32.IEEETable, b)
	f.written += uint32(len(b))
	return len(b), nil
}

func (f *FastGzipOutputStream) Finish() {
	if f.finished {
		return
	}
	f.finished = true
	if err := f.fw.Close(); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	var trailer [8]byte
	trailer[0] = byte(f.crc)
	trailer[1] = byte(f.crc >> 8)
	trailer[2] = byte(f.crc >> 16)
	trailer[3] = byte(f.crc >> 24)
	trailer[4] = byte(f.written)
	trailer[5] = byte(f.written >> 8)
	trailer[6] = byte(f.written >> 16)
	trailer[7] = byte(f.written >> 24)
	if _, err := f.bw.Write(trailer[:]); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	if err := f.bw.Flush(); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
}

func (f *FastGzipOutputStream) Flush() {
	if err := f.bw.Flush(); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	if flusher, ok := f.out.(interface{ Flush() error }); ok {
		if err := flusher.Flush(); err != nil {
			panic(jutil.NewIOException(err.Error()))
		}
	}
}

func (f *FastGzipOutputStream) Close() error {
	if f.closed {
		return nil
	}
	f.closed = true
	finished := false
	defer func() {
		if c, ok := f.out.(io.Closer); ok {
			if err := c.Close(); err != nil && finished {
				panic(jutil.NewIOException(err.Error()))
			}
		}
	}()
	f.Finish()
	finished = true
	return nil
}

const parallelGzipChunkSize = 512 * 1024

var parallelGzipMemberHeader = [10]byte{0x1f, 0x8b, 8, 0, 0, 0, 0, 0, 0, 0xff}

type parallelGzipWorker struct {
	fw     *flate.Writer
	outBuf bytes.Buffer
}

func (w *parallelGzipWorker) compressMember(chunk []byte) []byte {
	w.outBuf.Reset()
	if w.outBuf.Cap() < len(chunk)/2+64 {
		w.outBuf.Grow(len(chunk)/2 + 64)
	}
	w.outBuf.Write(parallelGzipMemberHeader[:])
	if w.fw == nil {
		fw, err := flate.NewWriter(&w.outBuf, flate.BestSpeed)
		if err != nil {
			panic(jutil.NewIllegalStateExceptionMsg(err.Error()))
		}
		w.fw = fw
	} else {
		w.fw.Reset(&w.outBuf)
	}
	if len(chunk) > 0 {
		if _, err := w.fw.Write(chunk); err != nil {
			panic(jutil.NewIOException(err.Error()))
		}
	}
	if err := w.fw.Close(); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	crc := crc32.ChecksumIEEE(chunk)
	isize := uint32(len(chunk))
	var trailer [8]byte
	trailer[0] = byte(crc)
	trailer[1] = byte(crc >> 8)
	trailer[2] = byte(crc >> 16)
	trailer[3] = byte(crc >> 24)
	trailer[4] = byte(isize)
	trailer[5] = byte(isize >> 8)
	trailer[6] = byte(isize >> 16)
	trailer[7] = byte(isize >> 24)
	w.outBuf.Write(trailer[:])
	return w.outBuf.Bytes()
}

// parallelFastGzipOutputStream compresses uncompressed data in fixed 512 KiB
// chunks across up to numThreads workers into valid concatenated RFC 1952 GZIP
// members, preserving deterministic output and bounded memory overhead.
type parallelFastGzipOutputStream struct {
	out         io.Writer
	numThreads  int
	curBuf      []byte
	batch       [][]byte
	freeChunks  [][]byte
	workers     []parallelGzipWorker
	wroteMember bool
	closed      bool
}

func newSaveStateGzipOutputStream(out io.Writer, numThreads int) io.WriteCloser {
	if numThreads <= 1 {
		return NewFastGzipOutputStream(out)
	}
	if out == nil {
		panic(jutil.NewNullPointerException())
	}
	return &parallelFastGzipOutputStream{
		out:        out,
		numThreads: numThreads,
		batch:      make([][]byte, 0, numThreads),
		freeChunks: make([][]byte, 0, numThreads+1),
		workers:    make([]parallelGzipWorker, numThreads),
	}
}

func (p *parallelFastGzipOutputStream) acquireChunk() []byte {
	if n := len(p.freeChunks); n > 0 {
		b := p.freeChunks[n-1]
		p.freeChunks = p.freeChunks[:n-1]
		return b[:0]
	}
	return make([]byte, 0, parallelGzipChunkSize)
}

func (p *parallelFastGzipOutputStream) Write(b []byte) (int, error) {
	if p.closed {
		panic(jutil.NewIOException("write beyond end of stream"))
	}
	total := len(b)
	for len(b) > 0 {
		if p.curBuf == nil {
			p.curBuf = p.acquireChunk()
		}
		space := parallelGzipChunkSize - len(p.curBuf)
		take := min(space, len(b))
		p.curBuf = append(p.curBuf, b[:take]...)
		b = b[take:]
		if len(p.curBuf) == parallelGzipChunkSize {
			p.batch = append(p.batch, p.curBuf)
			p.curBuf = nil
			if len(p.batch) == p.numThreads {
				p.flushBatch()
			}
		}
	}
	return total, nil
}

func (p *parallelFastGzipOutputStream) flushBatch() {
	n := len(p.batch)
	if n == 0 {
		return
	}
	if n == 1 {
		compressed := p.workers[0].compressMember(p.batch[0])
		if _, err := p.out.Write(compressed); err != nil {
			panic(jutil.NewIOException(err.Error()))
		}
	} else {
		var wg sync.WaitGroup
		wg.Add(n)
		for i := 0; i < n; i++ {
			idx := i
			go func() {
				defer wg.Done()
				p.workers[idx].compressMember(p.batch[idx])
			}()
		}
		wg.Wait()
		for i := 0; i < n; i++ {
			if _, err := p.out.Write(p.workers[i].outBuf.Bytes()); err != nil {
				panic(jutil.NewIOException(err.Error()))
			}
		}
	}
	p.wroteMember = true
	for i := 0; i < n; i++ {
		p.freeChunks = append(p.freeChunks, p.batch[i][:0])
		p.batch[i] = nil
	}
	p.batch = p.batch[:0]
}

func (p *parallelFastGzipOutputStream) Close() error {
	if p.closed {
		return nil
	}
	p.closed = true
	finished := false
	defer func() {
		if c, ok := p.out.(io.Closer); ok {
			if err := c.Close(); err != nil && finished {
				panic(jutil.NewIOException(err.Error()))
			}
		}
	}()
	if len(p.curBuf) > 0 || !p.wroteMember {
		if p.curBuf == nil {
			p.curBuf = make([]byte, 0)
		}
		p.batch = append(p.batch, p.curBuf)
		p.curBuf = nil
	}
	p.flushBatch()
	finished = true
	return nil
}
